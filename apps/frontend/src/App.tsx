import { useEffect } from 'react';
import { refresh } from './api/client';
import { useNavigation } from './stores/navigation';
import {
  ProductFooter,
  ProductHeader,
} from './components/layout/ProductHeader';
import { AppLink } from './components/ui/Button';
import { NavigationGuard } from './routing/NavigationGuard';
import { Routes } from './routing/Routes';

export function App() {
  const path = useNavigation((state) => state.path);
  const search = useNavigation((state) => state.search);
  useEffect(() => {
    void refresh();
  }, []);
  useEffect(() => {
    window.scrollTo(0, 0);
    document.title = '질문과 답변 · echotalk';
    const updateHeading = () => {
      const heading = document.querySelector<HTMLElement>('main h1');
      if (!heading) return false;
      if (!location.hash) heading.focus({ preventScroll: true });
      document.title = `${heading.textContent} · echotalk`;
      return true;
    };
    if (updateHeading()) return;
    // Async routes render their heading after loading the resource.
    const observer = new MutationObserver(() => {
      if (updateHeading()) observer.disconnect();
    });
    observer.observe(document.getElementById('root') ?? document.body, {
      childList: true,
      subtree: true,
    });
    return () => observer.disconnect();
  }, [path, search]);
  return (
    <div className="product">
      <AppLink
        className="skip-link"
        href="#main-content"
        onClick={(event) => {
          event.preventDefault();
          document.getElementById('main-content')?.scrollIntoView();
          document.querySelector<HTMLElement>('main h1, main button')?.focus();
        }}
      >
        본문 바로가기
      </AppLink>
      <NavigationGuard />
      <ProductHeader />
      <Routes />
      <ProductFooter />
    </div>
  );
}
