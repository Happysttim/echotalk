import { ApiError, jsonBody, request } from './client';
import { useSession } from '../stores/session';

export const authApi = {
  async login(email: string, password: string) {
    const data = await request<{ accessToken: string }>('/auth/local', {
      method: 'POST',
      body: jsonBody({ email, password }),
    });
    useSession.getState().setToken(data.accessToken);
  },
  async logout() {
    await request('/auth/logout', {}, 'required');
    useSession.getState().clear();
  },
  // The register response creates a server session. Keep the UI consistent with it.
  async register(email: string, password: string, nickname: string) {
    const data = await request<{ accessToken: string }>('/auth/register', {
      method: 'POST',
      body: jsonBody({ email, password, nickname }),
    });
    useSession.getState().setToken(data.accessToken);
  },
  sendVerification: (email: string, verifyType: 'register' | 'password') =>
    request(`/verify?${new URLSearchParams({ email, verifyType })}`),
  verifyMe: () => request<string>('/verify/me'),
  async verifyCode(code: string, verifyType: 'register' | 'password') {
    const verified = await request<boolean | void>(
      `/code?${new URLSearchParams({ code, verifyType })}`,
      { cache: 'no-store', referrerPolicy: 'no-referrer' },
    );
    // Current backend returns status: ok without data; an explicit false is never success.
    if (verified === false)
      throw new ApiError('인증 링크를 확인해 주세요.', 400);
  },
  resetPassword: (email: string, password: string) =>
    request('/auth/change', {
      method: 'POST',
      body: jsonBody({ email, password }),
    }),
};
