import { useEffect, useRef } from 'react';
import { AlertDialog } from 'radix-ui';
import { useNavigation } from '../stores/navigation';
import { Button } from '../components/ui/Button';
function navigate(path: string, replace: boolean) {
  if (replace) history.replaceState({}, '', path);
  else history.pushState({}, '', path);
  useNavigation.getState().sync();
}
export function NavigationGuard() {
  const pending = useNavigation((state) => state.pending);
  const returnFocus = useRef<HTMLElement | null>(null);
  useEffect(() => {
    const onNavigate = (event: Event) => {
      const detail = (event as CustomEvent<{ path: string; replace: boolean }>)
        .detail;
      if (useNavigation.getState().dirty) {
        returnFocus.current =
          document.activeElement instanceof HTMLElement
            ? document.activeElement
            : null;
        useNavigation.setState({ pending: detail });
      } else navigate(detail.path, detail.replace);
    };
    const onPop = () => {
      const state = useNavigation.getState();
      const target = `${location.pathname}${location.search}${location.hash}`;
      if (state.dirty) {
        returnFocus.current =
          document.activeElement instanceof HTMLElement
            ? document.activeElement
            : null;
        history.pushState({}, '', `${state.path}${state.search}${state.hash}`);
        useNavigation.setState({ pending: { path: target, replace: true } });
      } else state.sync();
    };
    const onUnload = (event: BeforeUnloadEvent) => {
      if (useNavigation.getState().dirty) event.preventDefault();
    };
    window.addEventListener('echotalk:navigate', onNavigate);
    window.addEventListener('popstate', onPop);
    window.addEventListener('beforeunload', onUnload);
    return () => {
      window.removeEventListener('echotalk:navigate', onNavigate);
      window.removeEventListener('popstate', onPop);
      window.removeEventListener('beforeunload', onUnload);
    };
  }, []);
  return (
    <AlertDialog.Root
      open={!!pending}
      onOpenChange={(open) => {
        if (!open) useNavigation.setState({ pending: null });
      }}
    >
      <AlertDialog.Portal>
        <AlertDialog.Overlay className="dialog-overlay" />
        <AlertDialog.Content
          className="dialog-content"
          onCloseAutoFocus={(event) => {
            event.preventDefault();
            if (
              useNavigation.getState().dirty &&
              returnFocus.current?.isConnected
            )
              returnFocus.current.focus();
          }}
        >
          <AlertDialog.Title>작성 중인 내용을 나갈까요?</AlertDialog.Title>
          <AlertDialog.Description>
            저장하지 않은 내용은 사라집니다.
          </AlertDialog.Description>
          <div className="dialog-actions">
            <AlertDialog.Cancel asChild>
              <Button variant="secondary">계속 작성</Button>
            </AlertDialog.Cancel>
            <AlertDialog.Action asChild>
              <Button
                variant="danger"
                onClick={() => {
                  if (pending) {
                    useNavigation.getState().setDirty(false);
                    navigate(pending.path, pending.replace);
                  }
                }}
              >
                나가기
              </Button>
            </AlertDialog.Action>
          </div>
        </AlertDialog.Content>
      </AlertDialog.Portal>
    </AlertDialog.Root>
  );
}
