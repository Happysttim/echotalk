import type { Answer, Survey } from '../types';

export const survey: Survey = {
  id: 'question-1',
  title: '새로운 기술을 배울 때, 어떤 순서로 공부하시나요?',
  content:
    '새로운 기술을 공부할 때마다 어디서 시작할지 고민하게 됩니다.\n\n공식 문서를 먼저 읽는 편인지, 작은 프로젝트를 만들면서 배우는 편인지 궁금합니다.',
  is_public: true,
  closed: false,
  author_id: 'user-1',
  expires_at: '2099-09-30T14:59:00Z',
  created_at: '2026-09-10T09:00:00Z',
  updated_at: '2026-09-10T09:00:00.123456789Z',
};
export const answer: Answer = {
  id: 'answer-1',
  survey_id: survey.id,
  author_id: 'user-1',
  is_anonymous: false,
  anonymous: '',
  content:
    '저는 먼저 아주 작은 결과물을 정합니다. 필요한 부분의 공식 문서를 읽으면서 배웁니다.',
  rate_up: 3,
  created_at: survey.created_at,
  updated_at: survey.updated_at,
};
export const token = (user = 'user-1') =>
  `header.${btoa(JSON.stringify({ user_id: user, token_type: 'access', exp: Math.floor(Date.now() / 1000) + 300 }))}.signature`;
