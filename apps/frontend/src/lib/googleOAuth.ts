import { env } from '../config/env';
import { safeNext } from './navigation';

export const GOOGLE_STATE_KEY = 'echotalk.google-oauth';
const MAX_STATE_AGE = 10 * 60 * 1000;
type Settings = { clientId: string; redirectUri: string };
type Transaction = {
  state: string;
  next: string;
  createdAt: number;
  redirectUri: string;
  clientId: string;
};

export function googleSettings(
  origin = location.origin,
  settings: Settings = {
    clientId: env.googleClientId,
    redirectUri: env.googleRedirectUri,
  },
): Settings {
  if (!settings.clientId)
    throw new Error(
      'Google 로그인 설정이 필요합니다. VITE_GOOGLE_CLIENT_ID를 확인해 주세요.',
    );
  const redirect = new URL(settings.redirectUri, origin);
  if (
    redirect.origin !== origin ||
    redirect.pathname !== '/auth/google/callback' ||
    redirect.search ||
    redirect.hash ||
    redirect.username ||
    redirect.password
  )
    throw new Error(
      'Google Redirect URI는 현재 사이트의 /auth/google/callback 경로여야 합니다.',
    );
  return { clientId: settings.clientId, redirectUri: redirect.href };
}

export function createGoogleAuthorizationUrl(
  next: string,
  settings = googleSettings(),
  storage = sessionStorage,
  now = Date.now(),
) {
  const state = Array.from(crypto.getRandomValues(new Uint8Array(32)), (byte) =>
    byte.toString(16).padStart(2, '0'),
  ).join('');
  const transaction: Transaction = {
    state,
    next: safeNext(next),
    createdAt: now,
    ...settings,
  };
  try {
    storage.setItem(GOOGLE_STATE_KEY, JSON.stringify(transaction));
  } catch {
    throw new Error(
      '로그인 확인 정보를 저장할 수 없습니다. 브라우저의 사이트 저장소 설정을 확인해 주세요.',
    );
  }
  const url = new URL('https://accounts.google.com/o/oauth2/v2/auth');
  url.search = new URLSearchParams({
    client_id: settings.clientId,
    redirect_uri: settings.redirectUri,
    response_type: 'code',
    scope: 'openid email profile',
    state,
  }).toString();
  return url.href;
}

export function consumeGoogleCallback(
  query: URLSearchParams,
  settings = googleSettings(),
  storage = sessionStorage,
  now = Date.now(),
) {
  let transaction: Partial<Transaction> | null = null;
  try {
    const raw = storage.getItem(GOOGLE_STATE_KEY);
    storage.removeItem(GOOGLE_STATE_KEY);
    transaction = raw ? (JSON.parse(raw) as Partial<Transaction>) : null;
  } catch {
    /* Missing or inaccessible transaction fails closed below. */
  }
  const state = query.get('state');
  if (
    !state ||
    query.getAll('state').length !== 1 ||
    !transaction ||
    transaction.state !== state ||
    typeof transaction.createdAt !== 'number' ||
    now < transaction.createdAt ||
    now - transaction.createdAt > MAX_STATE_AGE ||
    transaction.clientId !== settings.clientId ||
    transaction.redirectUri !== settings.redirectUri
  )
    throw new Error(
      '로그인 요청을 확인하지 못했거나 만료되었습니다. Google 로그인을 다시 시작해 주세요.',
    );
  if (query.has('error'))
    throw new Error(
      query.get('error') === 'access_denied'
        ? 'Google 로그인이 취소되었습니다. 다시 시도하거나 이메일로 로그인해 주세요.'
        : 'Google에서 로그인을 완료하지 못했습니다. 다시 시도해 주세요.',
    );
  const code = query.get('code');
  if (!code || query.getAll('code').length !== 1)
    throw new Error(
      'Google 인증 코드가 없습니다. 로그인을 다시 시작해 주세요.',
    );
  return {
    code,
    next: safeNext(
      typeof transaction.next === 'string' ? transaction.next : null,
    ),
  };
}
