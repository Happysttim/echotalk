import { useEffect } from 'react';
import { useSession } from '../stores/session';
import { useNavigation } from '../stores/navigation';
import { go, loginPath, safeNext } from '../lib/navigation';
import { LinkButton } from '../components/ui/Button';
import { EmptyState, LoadingRows } from '../components/ui/Feedback';
import { HomePage } from '../pages/HomePage';
import { LoginPage } from '../pages/LoginPage';
import { VerificationPage } from '../pages/VerificationPage';
import { QuestionsPage } from '../pages/QuestionsPage';
import { QuestionDetailPage } from '../pages/QuestionDetailPage';
import { SurveyEditorPage } from '../pages/SurveyEditorPage';
import { AnswerComposerPage } from '../pages/AnswerComposerPage';
import { GoogleCallbackPage } from '../pages/GoogleCallbackPage';
import { MyAnswersPage } from '../pages/MyAnswersPage';
export function Routes() {
  const { path, search } = useNavigation();
  const userId = useSession((state) => state.userId);
  const ready = useSession((state) => state.ready);
  const query = new URLSearchParams(search);
  const match = /^\/questions\/([^/]+)(?:\/(edit|answer))?$/.exec(path);
  const anonymous = query.get('mode') === 'anonymous' && !query.get('edit');
  const protectedPage =
    path === '/my/questions' ||
    path === '/my/answers' ||
    path === '/questions/new' ||
    match?.[2] === 'edit' ||
    (match?.[2] === 'answer' && !anonymous);
  useEffect(() => {
    if (protectedPage && ready && !userId)
      go(loginPath(`${path}${search}`), true);
  }, [protectedPage, ready, userId, path, search]);
  if (protectedPage && (!ready || !userId))
    return (
      <main id="main-content" className="page-width">
        <LoadingRows />
      </main>
    );
  if (path === '/') return <HomePage />;
  if (path === '/auth/google/callback') return <GoogleCallbackPage />;
  if (path === '/login')
    return <LoginPage next={safeNext(query.get('next'))} />;
  if (path === '/register')
    return (
      <VerificationPage
        key="register"
        kind="register"
        verifyToken={query.get('verify')}
      />
    );
  if (path === '/change' || path === '/password/reset')
    return (
      <VerificationPage
        key="password"
        kind="password"
        verifyToken={query.get('verify')}
      />
    );
  if (path === '/questions') return <QuestionsPage />;
  if (path === '/my/questions')
    return <QuestionsPage key={`my-questions-${userId}`} mine />;
  if (path === '/my/answers')
    return <MyAnswersPage key={`my-answers-${userId}`} />;
  if (path === '/questions/new') return <SurveyEditorPage key="new" />;
  if (match?.[1] && match[2] === 'edit')
    return <SurveyEditorPage key={path} id={match[1]} />;
  if (match?.[1] && match[2] === 'answer')
    return (
      <AnswerComposerPage
        key={path + search}
        id={match[1]}
        anonymous={anonymous}
        editId={query.get('edit') || undefined}
      />
    );
  if (match?.[1]) return <QuestionDetailPage key={path} id={match[1]} />;
  return (
    <main id="main-content" className="page-width">
      <EmptyState title="페이지를 찾을 수 없습니다.">
        <LinkButton href="/">홈으로 이동</LinkButton>
      </EmptyState>
    </main>
  );
}
