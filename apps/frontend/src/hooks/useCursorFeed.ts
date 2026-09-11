import { useCallback, useEffect, useRef, useState } from 'react';
import type { Page } from '../types';

export function appendUnique<T extends { id: string }>(
  old: T[],
  incoming: T[],
) {
  const ids = new Set(old.map((item) => item.id));
  return [
    ...old,
    ...incoming.filter((item) => {
      if (ids.has(item.id)) return false;
      ids.add(item.id);
      return true;
    }),
  ];
}
export function useCursorFeed<T extends { id: string }>(
  loadPage: (cursor: string, signal: AbortSignal) => Promise<Page<T>>,
) {
  const [items, setItems] = useState<T[]>([]);
  const [status, setStatus] = useState<'loading' | 'error' | 'done'>('loading');
  const [moreBusy, setMoreBusy] = useState(false);
  const [moreError, setMoreError] = useState('');
  const [hasNext, setHasNext] = useState(false);
  const cursor = useRef('');
  const generation = useRef(0);
  const controller = useRef<AbortController | null>(null);
  const locked = useRef(false);
  const load = useCallback(
    async (more: boolean) => {
      if (more && locked.current) return;
      if (!more) {
        controller.current?.abort();
        generation.current++;
        cursor.current = '';
        setItems([]);
        setStatus('loading');
        setHasNext(false);
      }
      const current = generation.current;
      controller.current = new AbortController();
      locked.current = true;
      setMoreBusy(more);
      setMoreError('');
      try {
        const page = await loadPage(
          more ? cursor.current : '',
          controller.current.signal,
        );
        if (current !== generation.current) return;
        setItems((old) => appendUnique(more ? old : [], page.items));
        cursor.current = page.nextCursor;
        setHasNext(page.hasNext && !!page.nextCursor);
        setStatus('done');
      } catch {
        if (
          current === generation.current &&
          !controller.current?.signal.aborted
        ) {
          if (more)
            setMoreError(
              '추가 항목을 불러오지 못했습니다. 다시 시도해 주세요.',
            );
          else setStatus('error');
        }
      } finally {
        if (current === generation.current) {
          locked.current = false;
          setMoreBusy(false);
        }
      }
    },
    [loadPage],
  );
  useEffect(() => {
    void load(false);
    return () => {
      generation.current++;
      controller.current?.abort();
      locked.current = false;
    };
  }, [load]);
  return {
    items,
    status,
    moreBusy,
    moreError,
    hasNext,
    more: () => {
      void load(true);
    },
    reload: () => {
      void load(false);
    },
  };
}
