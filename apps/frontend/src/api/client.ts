import { env } from '../config/env';
import { tokenUserId, useSession } from '../stores/session';

type Envelope<T> =
  | { status: 'ok'; data: T }
  | { status: 'failed'; error?: string; reason?: string };
type AuthMode = 'none' | 'required' | 'optional';
export class ApiError extends Error {
  constructor(
    message: string,
    public readonly status: number,
  ) {
    super(message);
  }
}

let refreshing: Promise<boolean> | null = null;
export function refresh(): Promise<boolean> {
  if (!refreshing) {
    refreshing = request<{ accessToken: string }>('/auth/refresh')
      .then((data) => {
        useSession.getState().setToken(data.accessToken);
        return true;
      })
      .catch(() => {
        useSession.getState().clear();
        return false;
      })
      .finally(() => {
        refreshing = null;
      });
  }
  return refreshing;
}

export async function request<T = void>(
  path: string,
  init: RequestInit = {},
  auth: AuthMode = 'none',
  retry = true,
): Promise<T> {
  let token = useSession.getState().token;
  if (auth !== 'none' && token && !tokenUserId(token)) {
    await refresh();
    token = useSession.getState().token;
  }
  if (auth === 'required' && !token)
    throw new ApiError('로그인이 필요합니다.', 401);
  const headers = new Headers(init.headers);
  headers.delete('Authorization');
  if (init.body) headers.set('Content-Type', 'application/json');
  if (auth !== 'none' && token) headers.set('Authorization', `Bearer ${token}`);
  let response: Response;
  try {
    response = await fetch(`${env.apiBaseUrl}${path}`, {
      ...init,
      headers,
      credentials: 'include',
    });
  } catch (error) {
    if (init.signal?.aborted) throw error;
    throw new ApiError(
      '서버에 연결하지 못했습니다. 잠시 후 다시 시도해 주세요.',
      0,
    );
  }
  // Retry read-only requests once. Never replay a submitted mutation automatically.
  if (response.status === 401 && auth !== 'none' && retry) {
    const refreshed =
      useSession.getState().token !== token || (await refresh());
    if (refreshed && (!init.method || init.method === 'GET'))
      return request<T>(path, init, auth, false);
  }
  const body = (await response.json().catch(() => null)) as Envelope<T> | null;
  if (!response.ok || body?.status !== 'ok')
    throw new ApiError(
      '요청을 처리하지 못했습니다. 입력 내용과 권한을 확인해 주세요.',
      response.status,
    );
  return body.data;
}

export const jsonBody = (payload: unknown) => JSON.stringify(payload);
export function errorText(
  error: unknown,
  fallback = '요청을 처리하지 못했습니다.',
) {
  if (error instanceof ApiError && error.status === 0) return error.message;
  if (error instanceof ApiError && error.status === 401)
    return '로그인이 필요하거나 세션이 만료되었습니다. 다시 로그인해 주세요.';
  return fallback;
}
