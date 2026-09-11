export function normalizeApiBase(value: string | undefined): string {
  const base = (value || '/api').trim().replace(/\/+$/, '');
  if (base.startsWith('/') && !base.startsWith('//') && !/[?#\\]/.test(base))
    return base;
  const url = new URL(base);
  if (
    !['http:', 'https:'].includes(url.protocol) ||
    url.username ||
    url.password ||
    url.search ||
    url.hash
  ) {
    throw new Error(
      'VITE_API_BASE_URL must be an HTTP(S) URL or an absolute path.',
    );
  }
  return base;
}

export const env = {
  apiBaseUrl: normalizeApiBase(import.meta.env.VITE_API_BASE_URL),
  googleClientId: import.meta.env.VITE_GOOGLE_CLIENT_ID?.trim() || '',
  googleRedirectUri:
    import.meta.env.VITE_GOOGLE_REDIRECT_URI?.trim() || '/auth/google/callback',
};
