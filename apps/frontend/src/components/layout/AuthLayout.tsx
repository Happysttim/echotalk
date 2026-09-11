import type { ReactNode } from 'react';
import { ArrowLeft, MessageCircle } from 'lucide-react';
import { AppLink } from '../ui/Button';

export function AuthLayout({
  title,
  intro,
  children,
  login = false,
}: {
  title: string;
  intro: string;
  children: ReactNode;
  login?: boolean;
}) {
  return (
    <main id="main-content" className="auth-layout">
      <aside className="auth-aside">
        <span className="eyebrow">ECHOTALK</span>
        <h2>
          질문을 읽고.
          <br />
          생각을 적고.
          <br />
          <span>답변을 나누고.</span>
        </h2>
        <div className="auth-aside-bottom">
          <MessageCircle size={36} aria-hidden />
          <p>
            질문 작성과 일반 답변 작성에는
            <br />
            로그인이 필요합니다.
          </p>
        </div>
      </aside>
      <section className="auth-form">
        <AppLink href={login ? '/' : '/login'} className="back-link">
          <ArrowLeft size={16} aria-hidden />
          {login ? '메인으로' : '로그인으로'}
        </AppLink>
        <h1 tabIndex={-1}>{title}</h1>
        <p className="form-intro">{intro}</p>
        {children}
      </section>
    </main>
  );
}
