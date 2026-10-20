// UI-only receipts for successful votes in this tab, never authorization data.
const KEY = 'echotalk.rating-receipts.v1';
const LIMIT = 1000;

export function readRatingReceipts(userId: string): string[] {
  try {
    const data: unknown = JSON.parse(sessionStorage.getItem(KEY) || 'null');
    if (!data || typeof data !== 'object') return [];
    const receipt = data as { userId?: unknown; answers?: unknown };
    if (receipt.userId !== userId || !Array.isArray(receipt.answers)) return [];
    return [
      ...new Set(
        receipt.answers.filter(
          (id): id is string => typeof id === 'string' && id.length > 0,
        ),
      ),
    ].slice(-LIMIT);
  } catch {
    return [];
  }
}

export function writeRatingReceipts(userId: string, answers: string[]) {
  try {
    sessionStorage.setItem(
      KEY,
      JSON.stringify({ userId, answers: answers.slice(-LIMIT) }),
    );
  } catch {
    // Storage restrictions must not make a successful vote appear to fail.
  }
}

export function clearRatingReceipts() {
  try {
    sessionStorage.removeItem(KEY);
  } catch {
    // The in-memory state is cleared regardless of browser storage access.
  }
}
