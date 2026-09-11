import { FileText, MessageSquare } from 'lucide-react';
import { AppLink } from '../ui/Button';

export function PersonalNavigation({
  current,
}: {
  current: 'questions' | 'answers';
}) {
  return (
    <nav className="personal-navigation" aria-label="내 활동">
      <AppLink
        href="/my/questions"
        aria-current={current === 'questions' ? 'page' : undefined}
      >
        <FileText size={16} aria-hidden />
        내가 작성한 질문
      </AppLink>
      <AppLink
        href="/my/answers"
        aria-current={current === 'answers' ? 'page' : undefined}
      >
        <MessageSquare size={16} aria-hidden />
        내가 작성한 답변
      </AppLink>
    </nav>
  );
}
