import { useState } from 'react';
import { Pencil, ThumbsUp, Trash2 } from 'lucide-react';
import type { Answer } from '../../types';
import { answerApi } from '../../api/questions';
import { useSession } from '../../stores/session';
import { go, loginPath } from '../../lib/navigation';
import { AuthorLine } from './QuestionRow';
import { Button, LinkButton } from '../ui/Button';
import { ConfirmDialog } from '../ui/ConfirmDialog';
import { InlineError } from '../ui/Feedback';

export function AnswerRow({
  answer,
  onDeleted,
}: {
  answer: Answer;
  onDeleted: () => void;
}) {
  const userId = useSession((state) => state.userId);
  const [count, setCount] = useState(answer.rate_up);
  const [busy, setBusy] = useState(false);
  const [rated, setRated] = useState(false);
  const [error, setError] = useState('');
  const owner = !!userId && userId === answer.author_id && !answer.is_anonymous;
  async function rate() {
    if (!userId) {
      go(loginPath(`/questions/${answer.survey_id}`));
      return;
    }
    if (busy || rated) return;
    setBusy(true);
    setError('');
    try {
      await answerApi.rateUp(answer.id);
      setRated(true);
      const updated = await answerApi.get(answer.id);
      setCount(updated.rate_up);
    } catch {
      setError('추천 또는 추천 수 갱신을 처리하지 못했습니다.');
    } finally {
      setBusy(false);
    }
  }
  return (
    <article className="answer-row" id={`answer-${answer.id}`} tabIndex={-1}>
      <AuthorLine
        anonymous={answer.is_anonymous ? answer.anonymous || '익명' : undefined}
        date={answer.created_at}
      />
      <p className="prose">{answer.content}</p>
      <div className="answer-actions">
        <button
          className={`vote ${rated ? 'voted' : ''}`}
          disabled={busy || rated}
          onClick={() => {
            void rate();
          }}
          aria-label={`추천 ${count}${rated ? ', 추천 완료' : ''}`}
        >
          <ThumbsUp size={15} aria-hidden />
          {busy ? '처리 중…' : rated ? '추천 완료' : '추천'}
          <span>{count}</span>
        </button>
        {owner && (
          <LinkButton
            variant="quiet"
            href={`/questions/${answer.survey_id}/answer?edit=${answer.id}`}
          >
            <Pencil size={14} aria-hidden />
            수정
          </LinkButton>
        )}
        {(owner || answer.is_anonymous) && (
          <ConfirmDialog
            trigger={
              <Button variant="quiet" className="danger">
                <Trash2 size={14} aria-hidden />
                삭제
              </Button>
            }
            title="답변을 삭제할까요?"
            description={
              answer.is_anonymous
                ? '이 답변을 작성할 때 설정한 비밀번호를 입력하세요.'
                : '삭제한 답변은 되돌릴 수 없습니다.'
            }
            action="답변 삭제"
            passwordRequired={answer.is_anonymous}
            failureMessage={
              answer.is_anonymous
                ? '비밀번호를 확인해 주세요. 답변을 삭제하지 못했습니다.'
                : undefined
            }
            onConfirm={async (password) => {
              await answerApi.remove(answer.id, answer.is_anonymous, password);
              onDeleted();
            }}
          />
        )}
      </div>
      {error && <InlineError>{error}</InlineError>}
    </article>
  );
}
