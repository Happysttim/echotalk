import { useCallback } from 'react';
import { ArrowRight, ArrowUpRight, Plus } from 'lucide-react';
import { questionApi } from '../api/questions';
import { useCursorFeed } from '../hooks/useCursorFeed';
import { QuestionRow } from '../components/questions/QuestionRow';
import { LinkButton } from '../components/ui/Button';
import { EmptyState, LoadingRows, RetryState } from '../components/ui/Feedback';
import { useWriteQuestionPath } from '../components/layout/ProductHeader';

export function HomePage() {
  const feed = useCursorFeed(
    useCallback(
      (cursor: string, signal: AbortSignal) =>
        questionApi.list('_id', cursor, signal),
      [],
    ),
  );
  const writePath = useWriteQuestionPath();
  return (
    <main id="main-content" className="page-width">
      <section className="hero">
        <div>
          <span className="eyebrow">질문과 답변</span>
          <h1 tabIndex={-1}>
            질문을 올리고,
            <br />
            <span>답변을 나눕니다.</span>
          </h1>
          <p>
            궁금한 내용을 텍스트로 작성하고
            <br className="mobile-only" /> 다른 사용자의 답변을 확인하세요.
          </p>
          <div className="hero-actions">
            <LinkButton href={writePath}>
              질문 작성
              <Plus size={18} aria-hidden />
            </LinkButton>
            <LinkButton href="/questions" variant="quiet">
              질문 둘러보기
              <ArrowRight size={18} aria-hidden />
            </LinkButton>
          </div>
        </div>
        <aside className="hero-note">
          <div className="note-rule" />
          <span className="eyebrow">ECHOTALK</span>
          <p>
            하나의 질문,
            <br />
            각자의 답변.
          </p>
          <span className="note-caption">
            계정으로 답변하거나
            <br />
            익명으로 답변할 수 있습니다.
          </span>
        </aside>
      </section>
      <section className="home-feed">
        <div className="section-head">
          <h2>최근 질문</h2>
          <LinkButton href="/questions" variant="quiet">
            전체 보기
            <ArrowUpRight size={17} aria-hidden />
          </LinkButton>
        </div>
        {feed.status === 'loading' ? (
          <LoadingRows />
        ) : feed.status === 'error' ? (
          <RetryState
            message="목록을 불러오지 못했습니다."
            onRetry={feed.reload}
          />
        ) : feed.items.length ? (
          feed.items
            .slice(0, 3)
            .map((item, index) => (
              <QuestionRow key={item.id} survey={item} index={index} />
            ))
        ) : (
          <EmptyState title="등록된 질문이 없습니다.">
            <LinkButton href={writePath}>질문 작성</LinkButton>
          </EmptyState>
        )}
      </section>
    </main>
  );
}
