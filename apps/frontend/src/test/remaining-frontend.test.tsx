// @vitest-environment jsdom
import { act, cleanup, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { afterEach, expect, it, vi } from 'vitest';
import { dateText, deadlineText, toKstInput } from '../lib/date';
import { readRatingReceipts, writeRatingReceipts } from '../lib/ratingReceipts';
import { AuthorLine } from '../components/questions/QuestionRow';
import { WritingAuthor } from '../components/questions/WritingAuthor';
import { AnswerRow } from '../components/questions/AnswerRow';
import { useSession } from '../stores/session';
import { answerApi } from '../api/questions';
import { answer, token } from './fixtures';

const author = { id: 'user-1', nickname: '서윤', hash_number: 10000001 };
afterEach(() => {
  cleanup();
  vi.restoreAllMocks();
  useSession.getState().clear();
  sessionStorage.clear();
});

it.each([
  '0001-01-01T00:00:00Z',
  '0001-01-01T00:00:00.000000000+00:00',
  '',
  'invalid',
  null,
  undefined,
])('does not render missing or zero timestamps as real dates: %s', (value) => {
  expect(dateText(value)).toBe('날짜 정보 없음');
  expect(toKstInput(value)).toBe('');
  expect(deadlineText(value)).toBe('마감일 정보 없음');
});

it('keeps valid KST formatting and the Unix epoch', () => {
  expect(dateText('2026-10-20T15:30:00Z')).toContain('2026');
  expect(toKstInput('2026-10-20T15:30:00Z')).toBe('2026-10-21T00:30');
  expect(deadlineText('2026-10-20T15:30:00Z')).toBe('2026-10-21 00:30 (KST)');
  expect(toKstInput('1970-01-01T00:00:00Z')).toBe('1970-01-01T09:00');
});

it('reuses a confirmed own nickname for new writing and clears it on account changes', () => {
  useSession.getState().setToken(token());
  const view = render(<AuthorLine author={author} />);
  view.unmount();
  render(<WritingAuthor />);
  expect(screen.getByText('서윤')).toBeTruthy();
  act(() => useSession.getState().setToken(token()));
  expect(screen.getByText('서윤')).toBeTruthy();
  act(() => useSession.getState().setToken(token('other')));
  expect(screen.queryByText('서윤')).toBeNull();
  expect(useSession.getState().author).toBeNull();
});

it('does not remember anonymous or other users as the current author', () => {
  useSession.getState().setToken(token());
  const view = render(<AuthorLine author={author} anonymous="익명" />);
  expect(useSession.getState().author).toBeNull();
  view.rerender(<AuthorLine author={{ ...author, id: 'other' }} />);
  expect(useSession.getState().author).toBeNull();
});

it('restores successful votes after remount and same-tab session restoration', async () => {
  useSession.getState().setToken(token());
  const vote = vi.spyOn(answerApi, 'rateUp').mockResolvedValue(undefined);
  vi.spyOn(answerApi, 'get').mockResolvedValue({ ...answer, rate_up: 4 });
  const view = render(<AnswerRow answer={answer} onDeleted={() => {}} />);
  await userEvent.click(screen.getByRole('button', { name: '추천 3' }));
  await waitFor(() =>
    expect(
      screen.getByRole('button', { name: '추천 4, 추천 완료' }),
    ).toHaveProperty('disabled', true),
  );
  view.unmount();
  // Simulate losing only in-memory state on reload, without logout.
  useSession.setState({
    token: '',
    userId: null,
    ratedAnswers: [],
    ready: false,
  });
  useSession.getState().setToken(token());
  render(<AnswerRow answer={{ ...answer, rate_up: 4 }} onDeleted={() => {}} />);
  expect(
    screen.getByRole('button', { name: '추천 4, 추천 완료' }),
  ).toHaveProperty('disabled', true);
  expect(vote).toHaveBeenCalledTimes(1);
  act(() => useSession.getState().setToken(token('other')));
  expect(screen.getByRole('button', { name: '추천 4' })).toHaveProperty(
    'disabled',
    false,
  );
  expect(readRatingReceipts('user-1')).toEqual([]);
});

it('does not mark rejected votes as successful', async () => {
  useSession.getState().setToken(token());
  vi.spyOn(answerApi, 'rateUp').mockRejectedValue(new Error('failed'));
  const get = vi.spyOn(answerApi, 'get');
  render(<AnswerRow answer={answer} onDeleted={() => {}} />);
  await userEvent.click(screen.getByRole('button', { name: '추천 3' }));
  expect(
    await screen.findByText(
      '추천을 처리하지 못했습니다. 잠시 후 다시 시도해 주세요.',
    ),
  ).toBeTruthy();
  expect(screen.getByRole('button', { name: '추천 3' })).toHaveProperty(
    'disabled',
    false,
  );
  expect(readRatingReceipts('user-1')).toEqual([]);
  expect(get).not.toHaveBeenCalled();
});

it('preserves successful votes when refreshing the count fails', async () => {
  useSession.getState().setToken(token());
  vi.spyOn(answerApi, 'rateUp').mockResolvedValue(undefined);
  vi.spyOn(answerApi, 'get').mockRejectedValue(new Error('offline'));
  render(<AnswerRow answer={answer} onDeleted={() => {}} />);
  await userEvent.click(screen.getByRole('button', { name: '추천 3' }));
  expect(
    await screen.findByText(
      '추천은 완료되었지만 추천 수를 불러오지 못했습니다.',
    ),
  ).toBeTruthy();
  expect(
    screen.getByRole('button', { name: '추천 3, 추천 완료' }),
  ).toHaveProperty('disabled', true);
  expect(readRatingReceipts('user-1')).toEqual([answer.id]);
});

it('does not leak an in-flight vote into a different account', async () => {
  useSession.getState().setToken(token());
  let finish!: () => void;
  vi.spyOn(answerApi, 'rateUp').mockImplementation(
    () =>
      new Promise<void>((resolve) => {
        finish = resolve;
      }),
  );
  vi.spyOn(answerApi, 'get').mockResolvedValue({ ...answer, rate_up: 4 });
  render(<AnswerRow answer={answer} onDeleted={() => {}} />);
  await userEvent.click(screen.getByRole('button', { name: '추천 3' }));
  act(() => useSession.getState().setToken(token('other')));
  await act(async () => {
    finish();
    await Promise.resolve();
  });
  expect(useSession.getState().ratedAnswers).toEqual([]);
  expect(screen.getByRole('button', { name: '추천 3' })).toHaveProperty(
    'disabled',
    false,
  );
});

it('uses new server counts when the answer prop changes', () => {
  const view = render(<AnswerRow answer={answer} onDeleted={() => {}} />);
  view.rerender(
    <AnswerRow answer={{ ...answer, rate_up: 9 }} onDeleted={() => {}} />,
  );
  expect(screen.getByRole('button', { name: '추천 9' })).toBeTruthy();
});

it('handles unavailable storage without failing successful votes', () => {
  useSession.getState().setToken(token());
  vi.spyOn(Storage.prototype, 'setItem').mockImplementation(() => {
    throw new Error('blocked');
  });
  vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => {
    throw new Error('blocked');
  });
  expect(readRatingReceipts('user-1')).toEqual([]);
  expect(() =>
    useSession.getState().recordRating('user-1', answer.id),
  ).not.toThrow();
  expect(useSession.getState().ratedAnswers).toEqual([answer.id]);
});

it('validates stored receipts, bounds them, and removes them on logout', () => {
  sessionStorage.setItem('echotalk.rating-receipts.v1', '{broken');
  expect(readRatingReceipts('user-1')).toEqual([]);
  writeRatingReceipts(
    'user-1',
    Array.from({ length: 1100 }, (_, i) => `a-${i}`),
  );
  expect(readRatingReceipts('other')).toEqual([]);
  expect(readRatingReceipts('user-1')).toHaveLength(1000);
  useSession.getState().setToken(token());
  useSession.getState().clear();
  expect(readRatingReceipts('user-1')).toEqual([]);
  expect(useSession.getState().ratedAnswers).toEqual([]);
});
