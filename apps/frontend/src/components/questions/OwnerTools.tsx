import { LockKeyhole, Pencil, Trash2 } from 'lucide-react';
import type { Survey } from '../../types';
import { questionApi } from '../../api/questions';
import { useSession } from '../../stores/session';
import { go } from '../../lib/navigation';
import { Button, LinkButton } from '../ui/Button';
import { ConfirmDialog } from '../ui/ConfirmDialog';

export function OwnerTools({
  survey,
  closed,
  onClosed,
}: {
  survey: Survey;
  closed: boolean;
  onClosed: () => void;
}) {
  const userId = useSession((state) => state.userId);
  if (!userId || userId !== survey.author_id) return null;
  return (
    <div className="owner-tools">
      <span className="eyebrow">내 질문 관리</span>
      <LinkButton href={`/questions/${survey.id}/edit`} variant="quiet">
        <Pencil size={15} aria-hidden />
        질문 수정
      </LinkButton>
      {!closed && (
        <ConfirmDialog
          trigger={
            <Button variant="quiet">
              <LockKeyhole size={15} aria-hidden />
              답변 마감
            </Button>
          }
          title="질문을 마감할까요?"
          description="마감하면 새 답변을 받을 수 없습니다. 기존 답변은 계속 표시됩니다."
          action="답변 마감"
          onConfirm={async () => {
            await questionApi.update(survey.id, {
              title: survey.title,
              content: survey.content,
              is_public: survey.is_public,
              expires_at: survey.expires_at,
              closed: true,
            });
            onClosed();
          }}
        />
      )}
      <ConfirmDialog
        trigger={
          <Button variant="quiet" className="danger">
            <Trash2 size={15} aria-hidden />
            질문 삭제
          </Button>
        }
        title="질문을 삭제할까요?"
        description="질문과 연결된 답변이 함께 삭제됩니다. 삭제한 내용은 되돌릴 수 없습니다."
        action="질문 삭제"
        onConfirm={async () => {
          await questionApi.remove(survey.id);
          go('/questions');
        }}
      />
    </div>
  );
}
