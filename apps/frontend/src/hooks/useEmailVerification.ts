import { useEffect, useRef, useState } from 'react';
import { authApi } from '../api/auth';
import { useNavigation } from '../stores/navigation';

export type VerificationKind = 'register' | 'password';

// A mail link is single-use. Share its request across StrictMode effect replays,
// but never persist the link token or the HttpOnly verification cookie.
export function useEmailVerification(
  token: string | null,
  kind: VerificationKind,
) {
  const [status, setStatus] = useState<
    'idle' | 'checking' | 'verified' | 'invalid'
  >(token === null ? 'idle' : 'checking');
  const request = useRef<{ token: string; promise: Promise<void> } | null>(
    null,
  );
  useEffect(() => {
    if (token === null) {
      // A navigation that removes a still-pending link cancels UI progression.
      // Removing it after a completed exchange must preserve the result form.
      if (request.current) {
        request.current = null;
        setStatus('idle');
      }
      return;
    }
    let active = true;
    setStatus('checking');
    if (request.current?.token !== token) {
      request.current = {
        token,
        promise: token.trim()
          ? authApi.verifyCode(token, kind)
          : Promise.reject(new Error('Empty verification link')),
      };
    }
    const clearLink = () => {
      const url = new URL(location.href);
      if (url.searchParams.get('verify') !== token) return;
      url.searchParams.delete('verify');
      history.replaceState(
        history.state,
        '',
        `${url.pathname}${url.search}${url.hash}`,
      );
      useNavigation.getState().sync();
    };
    void request.current.promise.then(
      () => {
        if (!active) return;
        setStatus('verified');
        request.current = null;
        clearLink();
      },
      () => {
        if (!active) return;
        setStatus('invalid');
        request.current = null;
        clearLink();
      },
    );
    return () => {
      active = false;
    };
  }, [token, kind]);
  return {
    status,
    reset: () => {
      request.current = null;
      setStatus('idle');
    },
  };
}
