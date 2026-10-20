import type { Author } from '../../types';
import { useSession } from '../../stores/session';
import { AuthorLine } from './QuestionRow';

export function WritingAuthor({ author }: { author?: Author | null }) {
  const userId = useSession((state) => state.userId);
  const knownAuthor = useSession((state) => state.author);
  const current =
    userId && author?.id === userId && author.nickname?.trim()
      ? author
      : knownAuthor?.id === userId
        ? knownAuthor
        : null;
  return (
    <div className="writing-author">
      <span>내 계정으로 작성</span>
      {current?.nickname?.trim() ? (
        <AuthorLine author={current} />
      ) : (
        <span className="help">등록 후 내 닉네임으로 표시됩니다.</span>
      )}
    </div>
  );
}
