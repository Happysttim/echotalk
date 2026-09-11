import { useEffect, useState } from 'react';
import { questionApi } from '../api/questions';
import { ApiError } from '../api/client';
import type { Survey } from '../types';

export function useSurvey(id: string) {
  const [survey, setSurvey] = useState<Survey | null>(null);
  const [error, setError] = useState('');
  const [version, setVersion] = useState(0);
  useEffect(() => {
    const controller = new AbortController();
    setSurvey(null);
    setError('');
    void questionApi
      .get(id, controller.signal)
      .then((data) => {
        if (!controller.signal.aborted) setSurvey(data);
      })
      .catch((reason: unknown) => {
        if (controller.signal.aborted) return;
        setError(
          reason instanceof ApiError && reason.status === 404
            ? '질문을 찾을 수 없습니다.'
            : reason instanceof ApiError && reason.status === 403
              ? '이 질문에 접근할 수 없습니다.'
              : '질문을 불러오지 못했습니다.',
        );
      });
    return () => controller.abort();
  }, [id, version]);
  return { survey, error, reload: () => setVersion((value) => value + 1) };
}
