import { ArrowUpRight, UserRound } from 'lucide-react';
import { useEffect } from 'react';
import { useSession } from '../../stores/session';
import type { Author, Survey } from '../../types';
import { dateText, isClosed } from '../../lib/date';
import { AppLink } from '../ui/Button';

export function StatusBadge({ closed }: { closed: boolean }) {
  return (
    <span className={`status ${closed ? 'is-closed' : ''}`}>
      {closed ? '마감' : '답변 받는 중'}
    </span>
  );
}
// Never expose the account summary for an anonymous answer.
export function AuthorLine({
  author,
  anonymous,
  date,
}: {
  author?: Author | null;
  anonymous?: string;
  date?: string;
}) {
  const userId = useSession((state) => state.userId);
  useEffect(() => {
    if (anonymous === undefined && author && userId === author.id)
      useSession.getState().rememberAuthor(author);
  }, [author, anonymous, userId]);
  return (
    <div className={date ? 'author-line' : 'byline'}>
      <span className="avatar" aria-hidden>
        {anonymous ? '익' : <UserRound size={15} />}
      </span>
      <span className="author-name">
        {anonymous || author?.nickname?.trim() || '작성자'}
        {!anonymous &&
          author?.nickname?.trim() &&
          Number.isSafeInteger(author.hash_number) &&
          author.hash_number >= 0 && (
            <span className="hash">#{author.hash_number}</span>
          )}
      </span>
      {anonymous && <span className="anon-tag">익명</span>}
      {date && <time dateTime={date}>{dateText(date)}</time>}
    </div>
  );
}
export function QuestionRow({
  survey,
  index,
}: {
  survey: Survey;
  index: number;
}) {
  const href = `/questions/${survey.id}`;
  return (
    <article className="question-row">
      <span className="row-index" aria-hidden>
        {String(index + 1).padStart(2, '0')}
      </span>
      <div className="question-row-body">
        <div className="row-meta">
          <StatusBadge closed={isClosed(survey)} />
          <time dateTime={survey.created_at}>
            {dateText(survey.created_at)}
          </time>
        </div>
        <AppLink className="question-title" href={href}>
          {survey.title}
        </AppLink>
        <p className="excerpt">{survey.content}</p>
        <AuthorLine author={survey.author} />
      </div>
      <AppLink
        className="row-arrow"
        href={href}
        aria-label={`${survey.title} 읽기`}
      >
        <ArrowUpRight size={22} aria-hidden />
      </AppLink>
    </article>
  );
}
