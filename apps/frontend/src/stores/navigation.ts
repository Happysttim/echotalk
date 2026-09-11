import { create } from 'zustand';

type LocationState = {
  path: string;
  search: string;
  hash: string;
  dirty: boolean;
  pending: { path: string; replace: boolean } | null;
  sync: () => void;
  setDirty: (dirty: boolean) => void;
};
export const useNavigation = create<LocationState>((set) => ({
  path: window.location.pathname,
  search: window.location.search,
  hash: window.location.hash,
  dirty: false,
  pending: null,
  sync: () =>
    set({
      path: location.pathname,
      search: location.search,
      hash: location.hash,
      pending: null,
    }),
  setDirty: (dirty) => set({ dirty }),
}));
