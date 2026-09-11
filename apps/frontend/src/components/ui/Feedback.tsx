import type { ReactNode } from 'react';
import { ArrowDown, MessageSquare, RotateCcw } from 'lucide-react';
import { Button } from './Button';

export function InlineError({ children }: { children: ReactNode }) {
  return (
    <p className="field-error inline-error" role="alert">
      {children}
    </p>
  );
}
export function LoadingRows() {
  return (
    <div className="loading-rows" role="status">
      <span className="sr-only">불러오는 중</span>
      {[1, 2, 3].map((key) => (
        <div key={key} aria-hidden>
          <i />
          <i />
          <i />
        </div>
      ))}
    </div>
  );
}
export function EmptyState({
  title,
  children,
}: {
  title: string;
  children?: ReactNode;
}) {
  return (
    <section className="empty-state">
      <MessageSquare size={32} aria-hidden />
      <h2>{title}</h2>
      {children}
    </section>
  );
}
export function RetryState({
  message,
  onRetry,
}: {
  message: string;
  onRetry: () => void;
}) {
  return (
    <EmptyState title={message}>
      <Button variant="secondary" onClick={onRetry}>
        <RotateCcw size={16} aria-hidden />
        다시 시도
      </Button>
    </EmptyState>
  );
}
export function PaginationMore({
  hasNext,
  busy,
  error,
  onMore,
  label,
  count,
}: {
  hasNext: boolean;
  busy: boolean;
  error: string;
  onMore: () => void;
  label: string;
  count: number;
}) {
  return (
    <>
      {error && <InlineError>{error}</InlineError>}
      {hasNext ? (
        <Button
          variant="secondary"
          className="load-more"
          disabled={busy}
          onClick={onMore}
        >
          {busy ? '불러오는 중…' : label}
          <ArrowDown size={16} aria-hidden />
        </Button>
      ) : (
        count > 0 && <p className="end-list">마지막 항목입니다.</p>
      )}
    </>
  );
}
