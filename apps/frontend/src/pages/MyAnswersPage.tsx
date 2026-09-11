import { useCallback } from 'react';
import { MessageSquare } from 'lucide-react';
import { answerApi } from '../api/questions';
import type { Answer } from '../types';
import { useCursorFeed } from '../hooks/useCursorFeed';
import { PersonalNavigation } from '../components/layout/PersonalNavigation';
import { MyAnswerGroup } from '../components/questions/MyAnswerGroup';
import { LinkButton } from '../components/ui/Button';
import {
  EmptyState,
  LoadingRows,
  PaginationMore,
  RetryState,
} from '../components/ui/Feedback';

export function groupMyAnswers(answers: Answer[]) {
  const groups = new Map<string, Answer[]>();
  for (const answer of answers) {
    const group = groups.get(answer.survey_id) ?? [];
    group.push(answer);
    groups.set(answer.survey_id, group);
  }
  return [...groups].map(([surveyId, items]) => ({ surveyId, answers: items }));
}
export function MyAnswersPage() {
  const feed = useCursorFeed(
    useCallback(
      (cursor: string, signal: AbortSignal) => answerApi.mine(cursor, signal),
      [],
    ),
  );
  return (
    <main id="main-content" className="page-width feed-layout">
      <aside className="feed-side">
        <span className="eyebrow">내 활동</span>
        <h2>내가 작성한 답변</h2>
        <p>
          내가 남긴 답변을
          <br />
          질문별로 모아보세요.
        </p>
        <LinkButton href="/questions">질문 둘러보기</LinkButton>
        <div className="side-note">
          <MessageSquare size={24} aria-hidden />
          <p>
            질문 제목을 선택하면
            <br />
            해당 질문으로 이동합니다.
          </p>
        </div>
      </aside>
      <section className="feed">
        <div className="feed-title">
          <h1 tabIndex={-1}>내가 작성한 답변</h1>
        </div>
        <PersonalNavigation current="answers" />
        <div className="sort-line">
          <span>질문별로 모아보기</span>
        </div>
        {feed.status === 'loading' ? (
          <LoadingRows />
        ) : feed.status === 'error' ? (
          <RetryState
            message="내 답변을 불러오지 못했습니다."
            onRetry={feed.reload}
          />
        ) : feed.items.length ? (
          groupMyAnswers(feed.items).map((group) => (
            <MyAnswerGroup key={group.surveyId} {...group} />
          ))
        ) : (
          <EmptyState title="아직 작성한 답변이 없습니다.">
            <LinkButton href="/questions">질문 둘러보기</LinkButton>
          </EmptyState>
        )}
        {feed.status === 'done' && (
          <PaginationMore
            hasNext={feed.hasNext}
            busy={feed.moreBusy}
            error={feed.moreError}
            onMore={feed.more}
            label="내 답변 더 보기"
            count={feed.items.length}
          />
        )}
      </section>
    </main>
  );
}
