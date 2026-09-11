import { useState } from 'react';
import {
  ArrowUpRight,
  FileText,
  MessageSquare,
  ChevronDown,
  LogOut,
  MessageCircle,
  Plus,
  UserRound,
} from 'lucide-react';
import { DropdownMenu } from 'radix-ui';
import { authApi } from '../../api/auth';
import { useSession } from '../../stores/session';
import { useNavigation } from '../../stores/navigation';
import { go, loginPath } from '../../lib/navigation';
import { AppLink, LinkButton } from '../ui/Button';
import { InlineError } from '../ui/Feedback';

export function ProductHeader() {
  const userId = useSession((state) => state.userId);
  const path = useNavigation((state) => state.path);
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);
  const authPage = [
    '/login',
    '/register',
    '/change',
    '/password/reset',
    '/auth/google/callback',
  ].includes(path);
  return (
    <>
      <header className="product-header">
        <div className="header-inner">
          <AppLink href="/" className="brand" aria-label="echotalk 홈">
            <MessageCircle strokeWidth={2.5} aria-hidden />
            echotalk<span className="brand-dot">.</span>
          </AppLink>
          {!authPage && (
            <nav className="top-nav" aria-label="주 메뉴">
              <AppLink
                href="/questions"
                aria-current={path === '/questions' ? 'page' : undefined}
              >
                질문 목록
              </AppLink>
            </nav>
          )}
          <div className="header-actions">
            {authPage ? (
              <LinkButton href="/questions" variant="quiet">
                질문 둘러보기
                <ArrowUpRight size={16} aria-hidden />
              </LinkButton>
            ) : userId ? (
              <>
                <LinkButton
                  href="/questions/new"
                  className="small write-desktop"
                >
                  <Plus size={16} aria-hidden />
                  질문 작성
                </LinkButton>
                <DropdownMenu.Root>
                  <DropdownMenu.Trigger
                    className="account-button"
                    aria-label="계정 메뉴"
                  >
                    <span className="avatar">
                      <UserRound size={16} aria-hidden />
                    </span>
                    <ChevronDown size={14} aria-hidden />
                  </DropdownMenu.Trigger>
                  <DropdownMenu.Portal>
                    <DropdownMenu.Content
                      className="menu-content"
                      align="end"
                      sideOffset={8}
                    >
                      <DropdownMenu.Label className="menu-label">
                        내 계정
                      </DropdownMenu.Label>
                      <DropdownMenu.Item asChild className="menu-item">
                        <AppLink href="/my/questions">
                          <FileText size={16} aria-hidden />
                          내가 작성한 질문
                        </AppLink>
                      </DropdownMenu.Item>
                      <DropdownMenu.Item asChild className="menu-item">
                        <AppLink href="/my/answers">
                          <MessageSquare size={16} aria-hidden />
                          내가 작성한 답변
                        </AppLink>
                      </DropdownMenu.Item>
                      <DropdownMenu.Item
                        className="menu-item"
                        disabled={busy}
                        onSelect={() => {
                          if (busy) return;
                          setBusy(true);
                          void authApi
                            .logout()
                            .then(() => {
                              setError('');
                              go('/');
                            })
                            .catch(() =>
                              setError(
                                '로그아웃하지 못했습니다. 다시 시도해 주세요.',
                              ),
                            )
                            .finally(() => setBusy(false));
                        }}
                      >
                        <LogOut size={16} aria-hidden />
                        {busy ? '로그아웃 중…' : '로그아웃'}
                      </DropdownMenu.Item>
                    </DropdownMenu.Content>
                  </DropdownMenu.Portal>
                </DropdownMenu.Root>
              </>
            ) : (
              <>
                <LinkButton
                  href="/login"
                  variant="quiet"
                  className="login-link"
                >
                  로그인
                </LinkButton>
                <LinkButton href="/register" className="small">
                  회원가입
                </LinkButton>
              </>
            )}
          </div>
        </div>
      </header>
      {error && (
        <div className="page-width">
          <InlineError>{error}</InlineError>
        </div>
      )}
    </>
  );
}
export function useWriteQuestionPath() {
  return useSession((state) => state.userId)
    ? '/questions/new'
    : loginPath('/questions/new');
}
export function ProductFooter() {
  return (
    <footer className="product-footer">
      <span>echotalk.</span>
      <span>질문과 답변을 나누는 공간</span>
    </footer>
  );
}
