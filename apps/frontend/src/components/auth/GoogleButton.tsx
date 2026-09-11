import { useState } from 'react';
import { Button } from '../ui/Button';
import { createGoogleAuthorizationUrl } from '../../lib/googleOAuth';
import { safeNext } from '../../lib/navigation';

export function GoogleButton({ register = false }: { register?: boolean }) {
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);
  return (
    <>
      <Button
        className="google full"
        variant="secondary"
        disabled={busy}
        onClick={() => {
          if (busy) return;
          setError('');
          try {
            const next = safeNext(
              new URLSearchParams(location.search).get('next'),
            );
            const url = createGoogleAuthorizationUrl(next);
            setBusy(true);
            location.assign(url);
          } catch (reason) {
            setBusy(false);
            setError(
              reason instanceof Error
                ? reason.message
                : 'Google 로그인을 시작하지 못했습니다.',
            );
          }
        }}
      >
        <span className="google-g" aria-hidden>
          G
        </span>
        {busy
          ? 'Google로 이동 중…'
          : register
            ? 'Google로 계속하기'
            : 'Google로 로그인'}
      </Button>
      {error && (
        <p className="help" role="alert">
          {error}
        </p>
      )}
      <div className="divider">
        또는 이메일로 {register ? '가입' : '로그인'}
      </div>
    </>
  );
}
