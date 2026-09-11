import { useEffect, useState } from 'react';
import { ArrowLeft, ArrowUpRight } from 'lucide-react';
import { answerApi } from '../api/questions';
import { errorText } from '../api/client';
import { useSurvey } from '../hooks/useSurvey';
import { useDirtyForm } from '../hooks/useDirtyForm';
import { useClock } from '../hooks/useClock';
import { useSession } from '../stores/session';
import { useNavigation } from '../stores/navigation';
import { isClosed } from '../lib/date';
import { go, loginPath } from '../lib/navigation';
import { AuthorLine } from '../components/questions/QuestionRow';
import { AppLink, Button, LinkButton } from '../components/ui/Button';
import {
  PasswordField,
  TextAreaField,
  TextField,
} from '../components/ui/Fields';
import {
  InlineError,
  LoadingRows,
  RetryState,
} from '../components/ui/Feedback';

export function AnswerComposerPage({
  id,
  anonymous,
  editId,
}: {
  id: string;
  anonymous: boolean;
  editId?: string;
}) {
  const { survey, error: surveyError, reload } = useSurvey(id);
  const userId = useSession((state) => state.userId);
  const time = useClock();
  const [content, setContent] = useState('');
  const [initial, setInitial] = useState('');
  const [name, setName] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loadError, setLoadError] = useState('');
  const [ready, setReady] = useState(!editId);
  const [busy, setBusy] = useState(false);
  const [version, setVersion] = useState(0);
  useDirtyForm(ready && (content !== initial || !!name || !!password));
  useEffect(() => {
    if (!editId) return;
    const controller = new AbortController();
    setLoadError('');
    void answerApi
      .get(editId, controller.signal)
      .then((answer) => {
        if (controller.signal.aborted) return;
        if (
          answer.survey_id !== id ||
          answer.is_anonymous ||
          answer.author_id !== userId
        ) {
          setLoadError('본인이 작성한 일반 답변만 수정할 수 있습니다.');
          return;
        }
        setContent(answer.content);
        setInitial(answer.content);
        setReady(true);
      })
      .catch(() => {
        if (!controller.signal.aborted)
          setLoadError('답변을 불러오지 못했습니다.');
      });
    return () => controller.abort();
  }, [editId, id, userId, version]);
  if (surveyError || loadError)
    return (
      <main id="main-content" className="page-width">
        <RetryState
          message={surveyError || loadError}
          onRetry={() => {
            reload();
            setVersion(version + 1);
          }}
        />
        <LinkButton href={`/questions/${id}`}>질문으로 돌아가기</LinkButton>
      </main>
    );
  if (!survey || !ready)
    return (
      <main id="main-content" className="page-width">
        <LoadingRows />
      </main>
    );
  const closed = isClosed(survey, time);
  const back = `/questions/${id}`;
  async function submit() {
    if (
      busy ||
      !content.trim() ||
      (closed && !editId) ||
      (anonymous && !password.trim())
    )
      return;
    setBusy(true);
    setError('');
    try {
      let answerId = editId;
      if (editId) await answerApi.update(id, editId, content.trim());
      else
        answerId = (
          await answerApi.create({
            survey_id: id,
            content: content.trim(),
            is_anonymous: anonymous,
            anonymous: anonymous ? name.trim() : '',
            answer_password: anonymous ? password : '',
          })
        ).id;
      useNavigation.getState().setDirty(false);
      go(`${back}#answer-${answerId}`);
    } catch (reason) {
      setError(
        errorText(
          reason,
          '답변을 저장하지 못했습니다. 작성한 내용은 유지됩니다.',
        ),
      );
    } finally {
      setBusy(false);
    }
  }
  return (
    <main id="main-content" className="page-width writing-layout">
      <AppLink href={back} className="back-link">
        <ArrowLeft size={16} aria-hidden />
        질문으로 돌아가기
      </AppLink>
      <div className="writing-grid">
        <section>
          <span className="eyebrow">
            {anonymous ? 'ANONYMOUS ANSWER' : 'ANSWER'}
          </span>
          <h1 tabIndex={-1}>
            {editId ? '답변 수정' : anonymous ? '익명으로 답변' : '답변 작성'}
          </h1>
          <section className="question-context">
            <span>대상 질문</span>
            <h2>{survey.title}</h2>
            <AppLink href={back}>
              질문 전체 보기
              <ArrowUpRight size={13} aria-hidden />
            </AppLink>
          </section>
          {!anonymous && (
            <div className="writing-author">
              <span>내 계정으로 작성</span>
              <AuthorLine />
            </div>
          )}
          <form
            aria-busy={busy}
            onSubmit={(event) => {
              event.preventDefault();
              void submit();
            }}
          >
            <fieldset disabled={busy}>
              {anonymous && (
                <div className="anon-fields">
                  <TextField
                    label="표시 이름 (선택)"
                    value={name}
                    onChange={(event) => setName(event.target.value)}
                    placeholder="답변에 표시할 이름"
                  />
                  <PasswordField
                    label="답변 비밀번호"
                    value={password}
                    onChange={(event) => setPassword(event.target.value)}
                    autoComplete="new-password"
                    help="익명 답변을 삭제할 때 필요합니다."
                    required
                  />
                </div>
              )}
              <TextAreaField
                label="답변 내용"
                value={content}
                onChange={(event) => setContent(event.target.value)}
                placeholder="질문에 대한 경험이나 생각을 작성하세요."
                rows={9}
                required
              />
            </fieldset>
            {closed && !editId && (
              <InlineError>
                답변이 마감되었습니다. 새로운 답변을 등록할 수 없습니다.
              </InlineError>
            )}
            {error && <InlineError>{error}</InlineError>}
            <div className="form-actions">
              <LinkButton href={back} variant="secondary">
                취소
              </LinkButton>
              <Button
                type="submit"
                disabled={
                  busy ||
                  !content.trim() ||
                  (anonymous && !password.trim()) ||
                  (closed && !editId)
                }
              >
                {busy ? '저장 중…' : editId ? '답변 저장' : '답변 등록'}
                <ArrowUpRight size={17} aria-hidden />
              </Button>
            </div>
          </form>
        </section>
        <aside className="writing-side">
          <span className="eyebrow">작성 안내</span>
          <h2>{anonymous ? '익명으로 남기는 답변' : '나의 경험을 답변으로'}</h2>
          <p>
            {anonymous
              ? '설정한 표시 이름이 답변에 나타납니다. 삭제할 때 사용할 비밀번호를 확인해 주세요.'
              : '질문과 관련된 경험이나 생각을 자유롭게 작성해 주세요.'}
          </p>
          <hr />
          <p>입력한 내용은 답변 피드에 표시됩니다.</p>
          {anonymous && (
            <AppLink
              className="text-link"
              href={loginPath(`/questions/${id}/answer`)}
            >
              계정으로 답변하려면 로그인
              <ArrowUpRight size={16} aria-hidden />
            </AppLink>
          )}
        </aside>
      </div>
    </main>
  );
}
