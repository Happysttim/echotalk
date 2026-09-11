import { useEffect, useState } from 'react';
import { ArrowLeft, ArrowUpRight } from 'lucide-react';
import { Switch } from 'radix-ui';
import { questionApi } from '../api/questions';
import { errorText } from '../api/client';
import {
  fromKstInput,
  toKstInput,
  minimumDeadline,
  isFutureDeadline,
} from '../lib/date';
import { useClock } from '../hooks/useClock';
import { go } from '../lib/navigation';
import { useSession } from '../stores/session';
import { useNavigation } from '../stores/navigation';
import { useDirtyForm } from '../hooks/useDirtyForm';
import { AppLink, Button, LinkButton } from '../components/ui/Button';
import { TextAreaField, TextField } from '../components/ui/Fields';
import {
  InlineError,
  LoadingRows,
  RetryState,
} from '../components/ui/Feedback';

export function SurveyEditorPage({ id }: { id?: string }) {
  const userId = useSession((state) => state.userId);
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [isPublic, setPublic] = useState(true);
  const [expiresAt, setExpiresAt] = useState('');
  const [closed, setClosed] = useState(false);
  const [initial, setInitial] = useState('');
  const [loaded, setLoaded] = useState(!id);
  const [loadError, setLoadError] = useState('');
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);
  const [version, setVersion] = useState(0);
  const time = useClock();
  const validDeadline = isFutureDeadline(expiresAt, time);
  useDirtyForm(
    loaded &&
      (id
        ? JSON.stringify({ title, content, isPublic, expiresAt }) !== initial
        : !!(title || content || expiresAt || !isPublic)),
  );
  useEffect(() => {
    if (!id) return;
    const controller = new AbortController();
    setLoadError('');
    void questionApi
      .get(id, controller.signal)
      .then((survey) => {
        if (controller.signal.aborted) return;
        if (survey.author_id !== userId) {
          setLoadError('본인이 작성한 질문만 수정할 수 있습니다.');
          return;
        }
        const date = toKstInput(survey.expires_at);
        setTitle(survey.title);
        setContent(survey.content);
        setPublic(survey.is_public);
        setExpiresAt(date);
        setClosed(survey.closed);
        setInitial(
          JSON.stringify({
            title: survey.title,
            content: survey.content,
            isPublic: survey.is_public,
            expiresAt: date,
          }),
        );
        setLoaded(true);
      })
      .catch(() => {
        if (!controller.signal.aborted)
          setLoadError('질문을 불러오지 못했습니다.');
      });
    return () => controller.abort();
  }, [id, userId, version]);
  async function submit() {
    if (busy || !title.trim() || !content.trim() || !expiresAt) return;
    if (!isFutureDeadline(expiresAt)) {
      setError('답변 마감일은 현재 시간 이후로 설정해 주세요.');
      return;
    }
    setError('');
    setBusy(true);
    try {
      const payload = {
        title: title.trim(),
        content: content.trim(),
        is_public: isPublic,
        expires_at: fromKstInput(expiresAt),
      };
      const questionId = id || (await questionApi.create(payload)).id;
      if (id) await questionApi.update(id, { ...payload, closed });
      useNavigation.getState().setDirty(false);
      go(`/questions/${questionId}`);
    } catch (reason) {
      setError(
        errorText(
          reason,
          '질문을 저장하지 못했습니다. 입력 내용을 확인해 주세요.',
        ),
      );
    } finally {
      setBusy(false);
    }
  }
  const back = id ? `/questions/${id}` : '/questions';
  if (loadError)
    return (
      <main id="main-content" className="page-width">
        <RetryState
          message={loadError}
          onRetry={() => setVersion(version + 1)}
        />
        <LinkButton href={back}>질문으로 돌아가기</LinkButton>
      </main>
    );
  if (!loaded)
    return (
      <main id="main-content" className="page-width">
        <LoadingRows />
      </main>
    );
  return (
    <main id="main-content" className="page-width writing-layout">
      <AppLink className="back-link" href={back}>
        <ArrowLeft size={16} aria-hidden />
        {id ? '질문으로 돌아가기' : '질문 목록'}
      </AppLink>
      <div className="writing-grid">
        <section>
          <span className="eyebrow">TEXT QUESTION</span>
          <h1 tabIndex={-1}>{id ? '질문 수정' : '질문 작성'}</h1>
          <p className="form-intro">
            제목과 내용을 작성하고 공개 여부와 마감일을 설정하세요.
          </p>
          <form
            aria-busy={busy}
            onSubmit={(event) => {
              event.preventDefault();
              void submit();
            }}
          >
            <fieldset disabled={busy}>
              <TextField
                label="제목"
                value={title}
                onChange={(event) => setTitle(event.target.value)}
                placeholder="궁금한 내용을 제목으로 작성하세요."
                required
              />
              <TextAreaField
                label="질문 내용"
                value={content}
                onChange={(event) => setContent(event.target.value)}
                placeholder="질문의 배경과 답변받고 싶은 내용을 작성하세요."
                rows={8}
                required
              />
              <div className="settings-row">
                <div>
                  <label htmlFor="public-switch">공개 여부</label>
                  <p>{isPublic ? '공개' : '비공개'}</p>
                </div>
                <Switch.Root
                  id="public-switch"
                  className="public-switch"
                  checked={isPublic}
                  onCheckedChange={setPublic}
                >
                  <Switch.Thumb className="switch-thumb" />
                </Switch.Root>
              </div>
              <TextField
                label="답변 마감일"
                type="datetime-local"
                min={minimumDeadline(time)}
                step={60}
                value={expiresAt}
                onChange={(event) => setExpiresAt(event.target.value)}
                help="한국 표준시(KST) 기준이며, 현재 시간 이후만 선택할 수 있습니다."
                error={
                  expiresAt && !validDeadline
                    ? '현재 시간 이후의 마감일을 선택해 주세요.'
                    : undefined
                }
                required
              />
            </fieldset>
            {error && <InlineError>{error}</InlineError>}
            <div className="form-actions">
              <LinkButton href={back} variant="secondary">
                취소
              </LinkButton>
              <Button
                type="submit"
                disabled={
                  busy || !title.trim() || !content.trim() || !validDeadline
                }
              >
                {busy ? '저장 중…' : id ? '질문 저장' : '질문 등록'}
                <ArrowUpRight size={17} aria-hidden />
              </Button>
            </div>
          </form>
        </section>
        <aside className="writing-side">
          <span className="eyebrow">작성 안내</span>
          <h2>무엇이 궁금한가요?</h2>
          <ol>
            {[
              ['질문을 요약한 제목', '무엇을 묻는지 제목에 담아주세요.'],
              ['답변에 필요한 배경', '상황이나 조건을 본문에 적어주세요.'],
              ['답변 마감일', '답변을 받을 기한을 설정하세요.'],
            ].map(([heading, text], index) => (
              <li key={heading}>
                <span>0{index + 1}</span>
                <div>
                  <strong>{heading}</strong>
                  <p>{text}</p>
                </div>
              </li>
            ))}
          </ol>
        </aside>
      </div>
    </main>
  );
}
