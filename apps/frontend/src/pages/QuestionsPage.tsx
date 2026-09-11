import { useCallback, useState } from 'react';
import { Select } from 'radix-ui';
import { Check, ChevronDown, MessageCircle, Plus } from 'lucide-react';
import { questionApi } from '../api/questions';
import { useCursorFeed } from '../hooks/useCursorFeed';
import { useWriteQuestionPath } from '../components/layout/ProductHeader';
import { LinkButton } from '../components/ui/Button';
import {
  EmptyState,
  LoadingRows,
  PaginationMore,
  RetryState,
} from '../components/ui/Feedback';
import { QuestionRow } from '../components/questions/QuestionRow';
import { PersonalNavigation } from '../components/layout/PersonalNavigation';

export function QuestionsPage({ mine = false }: { mine?: boolean }) {
  const title = mine ? '내가 작성한 질문' : '질문 목록';
  const [sort, setSort] = useState<'_id' | 'updated_at'>('_id');
  const writePath = useWriteQuestionPath();
  const feed = useCursorFeed(
    useCallback(
      (cursor: string, signal: AbortSignal) =>
        mine
          ? questionApi.mine(sort, cursor, signal)
          : questionApi.list(sort, cursor, signal),
      [sort, mine],
    ),
  );
  return (
    <main id="main-content" className="page-width feed-layout">
      <aside className="feed-side">
        <span className="eyebrow">{mine ? '내 활동' : '둘러보기'}</span>
        <h2>{title}</h2>
        {mine ? (
          <p>
            내가 올린 질문을 확인하고
            <br />
            답변을 이어서 읽어보세요.
          </p>
        ) : (
          <p>
            다른 사용자의 질문을 읽고
            <br />
            답변을 남겨보세요.
          </p>
        )}
        <LinkButton href={writePath}>
          <Plus size={16} aria-hidden />
          질문 작성
        </LinkButton>
        <div className="side-note">
          <MessageCircle size={24} aria-hidden />
          {mine ? (
            <p>
              질문을 선택하면
              <br />
              상세와 답변을 확인할 수 있습니다.
            </p>
          ) : (
            <p>
              로그인 없이도 익명 답변을
              <br />
              작성할 수 있습니다.
            </p>
          )}
        </div>
      </aside>
      <section className="feed">
        <div className="feed-title">
          <h1 tabIndex={-1}>{title}</h1>
          <LinkButton href={writePath} className="mobile-only">
            질문 작성
          </LinkButton>
        </div>
        {mine && <PersonalNavigation current="questions" />}
        <div className="sort-line">
          <span>{mine ? '내가 작성한 공개 질문' : '공개 질문'}</span>
          <Select.Root
            value={sort}
            onValueChange={(value) => setSort(value as typeof sort)}
          >
            <Select.Trigger className="sort-select" aria-label="질문 정렬">
              <Select.Value />
              <Select.Icon>
                <ChevronDown size={16} aria-hidden />
              </Select.Icon>
            </Select.Trigger>
            <Select.Portal>
              <Select.Content
                className="menu-content"
                position="popper"
                align="end"
                sideOffset={6}
              >
                <Select.Viewport>
                  {[
                    ['_id', '등록순'],
                    ['updated_at', '수정순'],
                  ].map(([value, label]) => (
                    <Select.Item
                      value={value!}
                      key={value}
                      className="menu-item"
                    >
                      <Select.ItemText>{label}</Select.ItemText>
                      <Select.ItemIndicator>
                        <Check size={15} aria-hidden />
                      </Select.ItemIndicator>
                    </Select.Item>
                  ))}
                </Select.Viewport>
              </Select.Content>
            </Select.Portal>
          </Select.Root>
        </div>
        {feed.status === 'loading' ? (
          <LoadingRows />
        ) : feed.status === 'error' ? (
          <RetryState
            message="목록을 불러오지 못했습니다."
            onRetry={feed.reload}
          />
        ) : feed.items.length ? (
          feed.items.map((item, index) => (
            <QuestionRow key={item.id} survey={item} index={index} />
          ))
        ) : (
          <EmptyState
            title={
              mine
                ? '아직 작성한 공개 질문이 없습니다.'
                : '등록된 질문이 없습니다.'
            }
          >
            <LinkButton href={writePath}>질문 작성</LinkButton>
          </EmptyState>
        )}
        {feed.status === 'done' && (
          <PaginationMore
            hasNext={feed.hasNext}
            busy={feed.moreBusy}
            error={feed.moreError}
            onMore={feed.more}
            label="질문 더 보기"
            count={feed.items.length}
          />
        )}
      </section>
    </main>
  );
}
