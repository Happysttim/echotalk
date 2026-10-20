const ZERO_TIME = Date.parse('0001-01-01T00:00:00Z');

function dateTimestamp(value: unknown) {
  if (typeof value !== 'string' || !value.trim()) return NaN;
  const time = Date.parse(value);
  return time === ZERO_TIME ? NaN : time;
}

export function dateText(value: string | null | undefined) {
  const time = dateTimestamp(value);
  if (!Number.isFinite(time)) return '날짜 정보 없음';
  const date = new Date(time);
  return new Intl.DateTimeFormat('ko-KR', {
    timeZone: 'Asia/Seoul',
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  }).format(date);
}

export function toKstInput(value: string | null | undefined) {
  const time = dateTimestamp(value);
  return !Number.isFinite(time)
    ? ''
    : new Date(time + 9 * 60 * 60 * 1000).toISOString().slice(0, 16);
}

export function fromKstInput(value: string) {
  if (!/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}$/.test(value))
    throw new Error('마감일을 확인해 주세요.');
  const iso = new Date(`${value}:00+09:00`).toISOString();
  if (toKstInput(iso) !== value) throw new Error('마감일을 확인해 주세요.');
  return iso;
}

// datetime-local uses minute precision: the current partial minute is already past.
export const minimumDeadline = (now = Date.now()) =>
  toKstInput(new Date((Math.floor(now / 60000) + 1) * 60000).toISOString());
export function isFutureDeadline(value: string, now = Date.now()) {
  try {
    return Date.parse(fromKstInput(value)) > now;
  } catch {
    return false;
  }
}

export function deadlineText(value: string | null | undefined) {
  const input = toKstInput(value);
  return input ? `${input.replace('T', ' ')} (KST)` : '마감일 정보 없음';
}
export const isClosed = (
  survey: { closed: boolean; expires_at: string },
  time = Date.now(),
) => survey.closed || Date.parse(survey.expires_at) <= time;
