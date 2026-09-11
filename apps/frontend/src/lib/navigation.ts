export function safeNext(value: string | null): string {
  if (
    !value ||
    !value.startsWith('/') ||
    value.startsWith('//') ||
    /[\\\u0000-\u0020]/.test(value)
  )
    return '/';
  const url = new URL(value, 'https://echotalk.invalid');
  return /^\/(?:questions(?:\/[a-zA-Z0-9_-]+)*(?:\/)?|my\/(?:questions|answers)|)$/.test(
    url.pathname,
  )
    ? `${url.pathname}${url.search}${url.hash}`
    : '/';
}

export const loginPath = (next: string) =>
  `/login?next=${encodeURIComponent(safeNext(next))}`;

export function go(path: string, replace = false) {
  window.dispatchEvent(
    new CustomEvent('echotalk:navigate', { detail: { path, replace } }),
  );
}
