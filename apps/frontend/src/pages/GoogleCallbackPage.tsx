import { useEffect, useRef, useState } from 'react';
import { AuthLayout } from '../components/layout/AuthLayout';
import { LinkButton } from '../components/ui/Button';
import { InlineError } from '../components/ui/Feedback';
import { consumeGoogleCallback } from '../lib/googleOAuth';
import { exchangeGoogleCode } from '../api/google';
import { errorText, refresh } from '../api/client';
import { go } from '../lib/navigation';
import { useNavigation } from '../stores/navigation';

export function GoogleCallbackPage() {
  const [query] = useState(() => new URLSearchParams(location.search));
  const [error, setError] = useState('');
  const completion = useRef<Promise<string> | null>(null);
  useEffect(() => {
    let active = true;
    if (!completion.current) {
      completion.current = (async () => {
        const transaction = consumeGoogleCallback(query);
        // Let the initial session probe finish before installing the new session.
        await refresh();
        await exchangeGoogleCode(transaction.code);
        return transaction.next;
      })();
      // Never leave authorization codes in history or subsequent request referrers.
      history.replaceState(history.state, '', location.pathname);
      useNavigation.getState().sync();
    }
    void completion.current.then(
      (next) => {
        if (active) go(next, true);
      },
      (reason: unknown) => {
        if (active)
          setError(
            reason instanceof Error
              ? errorText(reason, reason.message)
              : 'Google 로그인에 실패했습니다.',
          );
      },
    );
    return () => {
      active = false;
    };
  }, [query]);
  return (
    <AuthLayout
      title="Google 로그인"
      intro="Google 계정 인증 결과를 확인합니다."
    >
      {error ? (
        <>
          <InlineError>{error}</InlineError>
          <LinkButton href="/login" className="full">
            로그인으로 돌아가기
          </LinkButton>
        </>
      ) : (
        <p role="status">Google 로그인 중… 잠시만 기다려 주세요.</p>
      )}
    </AuthLayout>
  );
}
