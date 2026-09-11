import { request } from './client';
import { useSession, tokenUserId } from '../stores/session';

// Kept separate from email authApi so the two login flows can evolve independently.
export async function exchangeGoogleCode(code: string) {
  const data = await request<{ accessToken: string }>(
    `/auth/google?${new URLSearchParams({ code })}`,
    { cache: 'no-store', referrerPolicy: 'no-referrer' },
  );
  if (!data?.accessToken || !tokenUserId(data.accessToken))
    throw new Error(
      'Google 로그인 응답을 확인하지 못했습니다. 다시 시도해 주세요.',
    );
  useSession.getState().setToken(data.accessToken);
}
