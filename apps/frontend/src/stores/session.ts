import { create } from 'zustand';

type Session = {
  token: string;
  userId: string | null;
  ready: boolean;
  setToken: (token: string) => void;
  clear: () => void;
};

// Claims are only used to present ownership controls. The server authorizes every write.
export function tokenUserId(token: string): string | null {
  try {
    const segment = token.split('.')[1];
    if (!segment) return null;
    const claims = JSON.parse(
      atob(segment.replace(/-/g, '+').replace(/_/g, '/')),
    ) as Record<string, unknown>;
    return claims.token_type === 'access' &&
      typeof claims.user_id === 'string' &&
      typeof claims.exp === 'number' &&
      claims.exp * 1000 > Date.now()
      ? claims.user_id
      : null;
  } catch {
    return null;
  }
}

export const useSession = create<Session>((set) => ({
  token: '',
  userId: null,
  ready: false,
  setToken: (token) => set({ token, userId: tokenUserId(token), ready: true }),
  clear: () => set({ token: '', userId: null, ready: true }),
}));
