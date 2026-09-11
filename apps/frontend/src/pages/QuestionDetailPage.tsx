import { useCallback, useEffect, useState } from 'react';
import type { Answer } from '../types';
import {
  ArrowLeft,
  CalendarDays,
  Globe2,
  LockKeyhole,
  MessageCircle,
} from 'lucide-react';
import { answerApi } from '../api/questions';
import { useSurvey } from '../hooks/useSurvey';
import { useCursorFeed } from '../hooks/useCursorFeed';
import { useClock } from '../hooks/useClock';
import { useSession } from '../stores/session';
import { useNavigation } from '../stores/navigation';
import { deadlineText, isClosed } from '../lib/date';
import { loginPath } from '../lib/navigation';
import { AppLink, LinkButton } from '../components/ui/Button';
import {
  EmptyState,
  LoadingRows,
  PaginationMore,
  RetryState,
} from '../components/ui/Feedback';
import { AuthorLine, StatusBadge } from '../components/questions/QuestionRow';
import { AnswerRow } from '../components/questions/AnswerRow';
import { OwnerTools } from '../components/questions/OwnerTools';

export function QuestionDetailPage({ id }: { id: string }) {
  const { survey, error, reload } = useSurvey(id);
  const userId = useSession((state) => state.userId);
  const time = useClock();
  const hash = useNavigation((state) => state.hash);
  const [linkedAnswer, setLinkedAnswer] = useState<Answer | null>(null);
  const feed = useCursorFeed(
    useCallback(
      (cursor: string, signal: AbortSignal) =>
        answerApi.list(id, cursor, signal),
      [id],
    ),
  );
  const linkedId = hash.startsWith('#answer-') ? hash.slice(8) : '';
  const hasLinkedAnswer = feed.items.some((answer) => answer.id === linkedId);
  useEffect(() => {
    setLinkedAnswer(null);
    if (!linkedId || hasLinkedAnswer || feed.status !== 'done') return;
    const controller = new AbortController();
    void answerApi
      .get(linkedId, controller.signal)
      .then((answer) => {
        if (!controller.signal.aborted && answer?.survey_id === id)
          setLinkedAnswer(answer);
      })
      .catch(() => {
        /* The regular feed remains available if the link is stale. */
      });
    return () => controller.abort();
  }, [linkedId, hasLinkedAnswer, feed.status, id]);
  useEffect(() => {
    if (hash && feed.status === 'done')
      document.getElementById(hash.slice(1))?.focus();
  }, [hash, feed.status, linkedAnswer, survey]);
  if (error)
    return (
      <main id="main-content" className="page-width">
        <RetryState message={error} onRetry={reload} />
        <LinkButton href="/questions" variant="quiet">
          질문 목록
        </LinkButton>
      </main>
    );
  if (!survey)
    return (
      <main id="main-content" className="page-width">
        <LoadingRows />
      </main>
    );
  const closed = isClosed(survey, time);
  const answerPath = userId
    ? `/questions/${id}/answer`
    : loginPath(`/questions/${id}/answer`);
  const primary = userId ? '답변 작성' : '로그인하고 답변';
  return (
    <main
      id="main-content"
      className={`page-width reading-layout ${!closed ? 'has-answer-bar' : ''}`}
    >
      <AppLink href="/questions" className="back-link">
        <ArrowLeft size={16} aria-hidden />
        질문 목록
      </AppLink>
      <div className="reading-grid">
        <article className="question-content">
          <div className="question-meta">
            <StatusBadge closed={closed} />
            <span>
              {survey.is_public ? (
                <Globe2 size={13} aria-hidden />
              ) : (
                <LockKeyhole size={13} aria-hidden />
              )}
              {survey.is_public ? '공개' : '비공개'}
            </span>
          </div>
          <h1 tabIndex={-1}>{survey.title}</h1>
          <AuthorLine date={survey.created_at} />
          <div className="question-prose prose">{survey.content}</div>
          <p className="deadline">
            <CalendarDays size={16} aria-hidden />
            답변 마감 · {deadlineText(survey.expires_at)}
          </p>
          <div className="mobile-owner">
            <OwnerTools survey={survey} closed={closed} onClosed={reload} />
          </div>
          <section className="answer-section">
            {linkedAnswer && !hasLinkedAnswer && (
              <section aria-label="링크로 이동한 답변">
                <h2>선택한 답변</h2>
                <AnswerRow
                  answer={linkedAnswer}
                  onDeleted={() => {
                    setLinkedAnswer(null);
                    useNavigation.setState({ hash: '' });
                    history.replaceState(
                      {},
                      '',
                      `${location.pathname}${location.search}`,
                    );
                    feed.reload();
                  }}
                />
              </section>
            )}
            <div className="section-head">
              <h2>답변</h2>
              <span className="help">등록된 순서로 표시</span>
            </div>
            {feed.status === 'loading' ? (
              <LoadingRows />
            ) : feed.status === 'error' ? (
              <RetryState
                message="답변을 불러오지 못했습니다."
                onRetry={feed.reload}
              />
            ) : feed.items.length ? (
              feed.items.map((answer) => (
                <AnswerRow
                  key={answer.id}
                  answer={answer}
                  onDeleted={feed.reload}
                />
              ))
            ) : (
              <EmptyState title="아직 답변이 없습니다." />
            )}
            {feed.status === 'done' && (
              <PaginationMore
                hasNext={feed.hasNext}
                busy={feed.moreBusy}
                error={feed.moreError}
                onMore={feed.more}
                label="답변 더 보기"
                count={feed.items.length}
              />
            )}
          </section>
        </article>
        <aside className="answer-aside">
          <div className="answer-prompt">
            <MessageCircle size={28} aria-hidden />
            <h2>{closed ? '답변이 마감되었습니다.' : '이 질문에 답변하기'}</h2>
            <p>
              {closed
                ? '기존 답변은 계속 읽을 수 있습니다.'
                : '본인의 경험이나 생각을 남겨주세요.'}
            </p>
            {!closed && (
              <>
                <LinkButton href={answerPath} className="full">
                  {primary}
                </LinkButton>
                <LinkButton
                  href={`/questions/${id}/answer?mode=anonymous`}
                  variant="secondary"
                  className="full"
                >
                  익명으로 답변
                </LinkButton>
                <span className="help">
                  익명 답변에는 삭제용 비밀번호가 필요합니다.
                </span>
              </>
            )}
          </div>
          <OwnerTools survey={survey} closed={closed} onClosed={reload} />
        </aside>
      </div>
      {!closed && (
        <div className="mobile-answer-bar">
          <LinkButton
            href={`/questions/${id}/answer?mode=anonymous`}
            variant="secondary"
          >
            익명 답변
          </LinkButton>
          <LinkButton href={answerPath}>{primary}</LinkButton>
        </div>
      )}
    </main>
  );
}
