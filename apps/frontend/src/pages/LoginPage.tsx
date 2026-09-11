import { useState } from 'react';
import { authApi } from '../api/auth';
import { ApiError, errorText } from '../api/client';
import { go } from '../lib/navigation';
import { AuthLayout } from '../components/layout/AuthLayout';
import { GoogleButton } from '../components/auth/GoogleButton';
import { AppLink, Button } from '../components/ui/Button';
import { PasswordField, TextField } from '../components/ui/Fields';
import { useSession } from '../stores/session';

export function LoginPage({ next }: { next: string }) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);
  const userId = useSession((state) => state.userId);
  return (
    <AuthLayout title="로그인" intro="가입한 방법으로 로그인하세요." login>
      <GoogleButton />
      {userId ? (
        <>
          <p className="verified">이미 로그인되어 있습니다.</p>
          <Button className="full" onClick={() => go(next)}>
            계속하기
          </Button>
        </>
      ) : (
        <form
          aria-busy={busy}
          onSubmit={(event) => {
            event.preventDefault();
            if (busy) return;
            setBusy(true);
            setError('');
            void authApi
              .login(email.trim(), password)
              .then(() => go(next, true))
              .catch((reason: unknown) =>
                setError(
                  reason instanceof ApiError && reason.status === 401
                    ? '이메일 또는 비밀번호를 확인해 주세요.'
                    : errorText(
                        reason,
                        '로그인하지 못했습니다. 잠시 후 다시 시도해 주세요.',
                      ),
                ),
              )
              .finally(() => setBusy(false));
          }}
        >
          <TextField
            label="이메일"
            type="email"
            autoComplete="email"
            placeholder="name@example.com"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            required
          />
          <PasswordField
            autoComplete="current-password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            required
            error={error || undefined}
          />
          <Button className="full" type="submit" disabled={busy}>
            {busy ? '로그인 중…' : '로그인'}
          </Button>
        </form>
      )}
      <div className="auth-links">
        <AppLink href="/change">비밀번호 변경</AppLink>
        <span />
        <AppLink href="/register">회원가입</AppLink>
      </div>
    </AuthLayout>
  );
}
