import type { Answer, Page, Survey } from '../types';
import { ApiError, jsonBody, request } from './client';

export function encodeCursor(value: Record<string, string>) {
  return btoa(
    String.fromCharCode(...new TextEncoder().encode(JSON.stringify(value))),
  )
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=+$/, '');
}
type SurveyInput = Pick<
  Survey,
  'title' | 'content' | 'is_public' | 'expires_at'
>;
export const questionApi = {
  async mine(
    type: '_id' | 'updated_at',
    cursor = '',
    signal?: AbortSignal,
  ): Promise<Page<Survey>> {
    const data = await request<{
      surveys: Survey[] | null;
      skip: string;
      hasNext: boolean;
    }>(
      `/surveys/me?${new URLSearchParams({ type, limit: '10', cursor })}`,
      { signal },
      'required',
    );
    return {
      items: data.surveys ?? [],
      nextCursor: data.skip,
      hasNext: data.hasNext,
    };
  },
  async list(
    type: '_id' | 'updated_at',
    cursor = '',
    signal?: AbortSignal,
  ): Promise<Page<Survey>> {
    const data = await request<{
      surveys: Survey[] | null;
      skip: string;
      hasNext: boolean;
    }>(`/surveys?${new URLSearchParams({ type, limit: '10', cursor })}`, {
      signal,
    });
    return {
      items: data.surveys ?? [],
      nextCursor: data.skip,
      hasNext: data.hasNext,
    };
  },
  async get(id: string, signal?: AbortSignal) {
    const data = await request<Survey | null>(
      `/surveys/${encodeURIComponent(id)}`,
      { signal },
    );
    if (!data) throw new ApiError('질문을 찾을 수 없습니다.', 404);
    return data;
  },
  create: (payload: SurveyInput) =>
    request<Survey>(
      '/surveys',
      { method: 'POST', body: jsonBody(payload) },
      'required',
    ),
  update: (id: string, payload: SurveyInput & { closed: boolean }) =>
    request(
      '/surveys',
      { method: 'PATCH', body: jsonBody({ survey_id: id, ...payload }) },
      'required',
    ),
  remove: (id: string) =>
    request(
      '/surveys',
      { method: 'DELETE', body: jsonBody({ survey_id: id }) },
      'required',
    ),
};
export const answerApi = {
  async mine(cursor = '', signal?: AbortSignal): Promise<Page<Answer>> {
    // /answers/me takes a plain answer ID, not the public feed's base64 cursor.
    const data = await request<{
      answers: Answer[] | null;
      skip: string;
      hasNext: boolean;
    }>(
      `/answers/me?${new URLSearchParams({ cursor })}`,
      { signal },
      'required',
    );
    return {
      items: data.answers ?? [],
      nextCursor: data.skip,
      hasNext: data.hasNext,
    };
  },
  async list(
    surveyId: string,
    cursor = '',
    signal?: AbortSignal,
  ): Promise<Page<Answer>> {
    const initial = encodeCursor({ survey_id: surveyId, answer_id: '' });
    const data = await request<{
      answers: Answer[] | null;
      skip: string;
      hasNext: boolean;
    }>(`/answers?${new URLSearchParams({ cursor: cursor || initial })}`, {
      signal,
    });
    return {
      items: data.answers ?? [],
      nextCursor: data.skip,
      hasNext: data.hasNext,
    };
  },
  get: (id: string, signal?: AbortSignal) =>
    request<Answer>(`/answers/${encodeURIComponent(id)}`, { signal }),
  create: (payload: {
    survey_id: string;
    content: string;
    is_anonymous: boolean;
    anonymous: string;
    answer_password: string;
  }) =>
    request<Answer>(
      '/answers',
      { method: 'POST', body: jsonBody(payload) },
      payload.is_anonymous ? 'none' : 'required',
    ),
  update: (surveyId: string, id: string, content: string) =>
    request(
      '/answers',
      {
        method: 'PATCH',
        body: jsonBody({ survey_id: surveyId, answer_id: id, content }),
      },
      'required',
    ),
  remove: (id: string, anonymous: boolean, password: string) =>
    request(
      '/answers',
      {
        method: 'DELETE',
        body: jsonBody({ answer_id: id, answer_password: password }),
      },
      anonymous ? 'none' : 'required',
    ),
  rateUp: (id: string) =>
    request(
      '/answers/rateup',
      { method: 'POST', body: jsonBody({ answer_id: id }) },
      'required',
    ),
};
