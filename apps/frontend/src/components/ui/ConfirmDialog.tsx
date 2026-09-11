import { useState } from 'react';
import type { ReactNode } from 'react';
import { AlertDialog } from 'radix-ui';
import { Button } from './Button';
import { PasswordField } from './Fields';
import { InlineError } from './Feedback';

export function ConfirmDialog({
  trigger,
  title,
  description,
  action,
  onConfirm,
  passwordRequired = false,
  failureMessage = '요청을 처리하지 못했습니다. 다시 시도해 주세요.',
}: {
  trigger: ReactNode;
  title: string;
  description: string;
  action: string;
  onConfirm: (password: string) => Promise<unknown>;
  passwordRequired?: boolean;
  failureMessage?: string;
}) {
  const [open, setOpen] = useState(false);
  const [password, setPassword] = useState('');
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  return (
    <AlertDialog.Root
      open={open}
      onOpenChange={(value) => {
        if (busy) return;
        setOpen(value);
        if (!value) {
          setPassword('');
          setError('');
        }
      }}
    >
      <AlertDialog.Trigger asChild>{trigger}</AlertDialog.Trigger>
      <AlertDialog.Portal>
        <AlertDialog.Overlay className="dialog-overlay" />
        <AlertDialog.Content className="dialog-content">
          <AlertDialog.Title>{title}</AlertDialog.Title>
          <AlertDialog.Description>{description}</AlertDialog.Description>
          <form
            onSubmit={(event) => {
              event.preventDefault();
              if (busy) return;
              setBusy(true);
              setError('');
              void onConfirm(password)
                .then(() => {
                  setOpen(false);
                  setPassword('');
                })
                .catch(() => setError(failureMessage))
                .finally(() => setBusy(false));
            }}
          >
            {passwordRequired && (
              <PasswordField
                label="답변 비밀번호"
                autoComplete="off"
                value={password}
                onChange={(event) => setPassword(event.target.value)}
                required
                error={error || undefined}
              />
            )}
            {error && !passwordRequired && <InlineError>{error}</InlineError>}
            <div className="dialog-actions">
              <AlertDialog.Cancel asChild>
                <Button variant="secondary" disabled={busy}>
                  취소
                </Button>
              </AlertDialog.Cancel>
              <Button
                variant="danger"
                type="submit"
                disabled={busy || (passwordRequired && !password.trim())}
              >
                {busy ? '처리 중…' : action}
              </Button>
            </div>
          </form>
        </AlertDialog.Content>
      </AlertDialog.Portal>
    </AlertDialog.Root>
  );
}
