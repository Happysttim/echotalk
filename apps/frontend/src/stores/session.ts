import { create } from 'zustand';
import type { Author } from '../types';
import {
  clearRatingReceipts,
  readRatingReceipts,
  writeRatingReceipts,
} from '../lib/ratingReceipts';

type Session = {
  token: string;
  userId: string | null;
  ready: boolean;
  author: Author | null;
  ratedAnswers: string[];
  rememberAuthor: (author: Author) => void;
  recordRating: (userId: string, answerId: string) => void;
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

export const useSession = create<Session>((set, get) => ({
  token: '',
  userId: null,
  ready: false,
  author: null,
  ratedAnswers: [],
  setToken: (token) => {
    const userId = tokenUserId(token);
    const previous = get();
    if (!userId || (previous.userId && previous.userId !== userId))
      clearRatingReceipts();
    set({
      token,
      userId,
      ready: true,
      author: userId && previous.userId === userId ? previous.author : null,
      ratedAnswers: userId
        ? previous.userId === userId
          ? previous.ratedAnswers
          : readRatingReceipts(userId)
        : [],
    });
  },
  rememberAuthor: (author) => {
    const state = get();
    if (!state.userId || author.id !== state.userId || !author.nickname?.trim())
      return;
    if (
      state.author?.nickname === author.nickname &&
      state.author.hash_number === author.hash_number
    )
      return;
    set({ author: { ...author } });
  },
  recordRating: (userId, answerId) => {
    const state = get();
    if (state.userId !== userId) return;
    const ratedAnswers = [...new Set([...state.ratedAnswers, answerId])].slice(
      -1000,
    );
    set({ ratedAnswers });
    writeRatingReceipts(userId, ratedAnswers);
  },
  clear: () => {
    clearRatingReceipts();
    set({
      token: '',
      userId: null,
      ready: true,
      author: null,
      ratedAnswers: [],
    });
  },
}));
