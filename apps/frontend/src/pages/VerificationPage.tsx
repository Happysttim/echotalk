import { useEffect, useState } from 'react';
import { Check, CheckCircle2, Mail } from 'lucide-react';
import { authApi } from '../api/auth';
import { errorText } from '../api/client';
import { AuthLayout } from '../components/layout/AuthLayout';
import { GoogleButton } from '../components/auth/GoogleButton';
import { AppLink, Button, LinkButton } from '../components/ui/Button';
import { PasswordField, TextField } from '../components/ui/Fields';
import { InlineError } from '../components/ui/Feedback';
import {
  useEmailVerification,
  type VerificationKind,
} from '../hooks/useEmailVerification';

export function VerificationPage({
  kind,
  verifyToken = null,
}: {
  kind: VerificationKind;
  verifyToken?: string | null;
}) {
  const verification = useEmailVerification(verifyToken, kind);
  const [waiting, setWaiting] = useState(false);
  const [email, setEmail] = useState('');
  const [nickname, setNickname] = useState('');
  const [password, setPassword] = useState('');
  const [confirmation, setConfirmation] = useState('');
  const [error, setError] = useState('');
  const [notice, setNotice] = useState('');
  const [busy, setBusy] = useState(false);
  const [done, setDone] = useState(false);
  const register = kind === 'register';
  const title = register ? '회원가입' : '비밀번호 변경';
  const step =
    verification.status === 'verified'
      ? 3
      : waiting || verification.status !== 'idle'
        ? 2
        : 1;
  const checking = verification.status === 'checking';
  const invalid = verification.status === 'invalid';
  useEffect(() => {
    if (verifyToken === null) return;
    setPassword('');
    setConfirmation('');
    setNickname('');
    setError('');
    setNotice('');
    const verifyMe = async () => {
      const myEmail = await authApi.verifyMe();
      setEmail(myEmail)
    };

    void verifyMe();
    setDone(false);
  }, [verifyToken]);
  const restart = () => {
    verification.reset();
    setWaiting(false);
    setError('');
    setNotice('');
    setPassword('');
    setConfirmation('');
    setNickname('');
  };
  async function run(
    action: () => Promise<unknown>,
    next: () => void,
    fallback: string,
  ) {
    if (busy) return;
    setBusy(true);
    setError('');
    setNotice('');
    try {
      await action();
      next();
    } catch (reason) {
      setError(errorText(reason, fallback));
    } finally {
      setBusy(false);
    }
  }
  const send = () => {
    void run(
      () => authApi.sendVerification(email.trim(), kind),
      () => {
        setWaiting(true);
        setNotice('인증 메일을 보냈습니다. 메일에 있는 링크를 열어 주세요.');
      },
      '인증 메일을 보내지 못했습니다. 이메일을 확인해 주세요.',
    );
  };
  const submit = () => {
    if (busy || checking) return;
    if (step === 1) send();
    else if (step !== 3 || verification.status !== 'verified') return;
    else if (password !== confirmation)
      setError('비밀번호 확인이 일치하지 않습니다.');
    else
      void run(
        () =>
          register
            ? authApi.register(email.trim(), password, nickname.trim())
            : authApi.resetPassword(email.trim(), password),
        () => {
          setDone(true);
          setPassword('');
          setConfirmation('');
        },
        '요청을 완료하지 못했습니다. 입력 내용을 확인해 주세요.',
      );
  };
  return (
    <AuthLayout
      title={title}
      intro={
        register
          ? 'Google 계정 또는 이메일로 가입하세요.'
          : '이메일로 가입한 계정의 비밀번호를 변경합니다.'
      }
    >
      {done ? (
        <section className="success-state">
          <Check size={32} aria-hidden />
          <h2>
            {register
              ? '회원가입이 완료되었습니다.'
              : '비밀번호가 변경되었습니다.'}
          </h2>
          <p>
            {register
              ? '이제 질문과 답변을 나눌 수 있습니다.'
              : '새 비밀번호로 로그인해 주세요.'}
          </p>
          <LinkButton href={register ? '/questions' : '/login'}>
            {register ? '질문 둘러보기' : '로그인으로 이동'}
          </LinkButton>
        </section>
      ) : (
        <>
          {step === 1 && register && <GoogleButton register />}
          <ol className="steps" aria-label="진행 단계">
            {['이메일', '인증', register ? '가입 정보' : '비밀번호'].map(
              (label, index) => (
                <li
                  key={label}
                  className={step === index + 1 ? 'current' : ''}
                  aria-current={step === index + 1 ? 'step' : undefined}
                >
                  <span>
                    {step > index + 1 ? (
                      <Check size={13} aria-hidden />
                    ) : (
                      index + 1
                    )}
                  </span>
                  {label}
                </li>
              ),
            )}
          </ol>
          <form
            aria-busy={busy || checking}
            onSubmit={(event) => {
              event.preventDefault();
              submit();
            }}
          >
            <fieldset disabled={busy || checking}>
              {step === 1 ? (
                <TextField
                  label="이메일"
                  type="email"
                  autoComplete="email"
                  placeholder="name@example.com"
                  value={email}
                  onChange={(event) => setEmail(event.target.value)}
                  required
                  error={error || undefined}
                />
              ) : step === 2 ? (
                <>
                  {checking ? (
                    <p role="status">이메일 인증 링크를 확인하고 있습니다…</p>
                  ) : invalid ? (
                    <>
                      <InlineError>
                        인증 링크를 확인하지 못했습니다. 만료되거나 이미 사용된
                        링크일 수 있습니다. 새 인증 메일을 받아 주세요.
                      </InlineError>
                      <Button variant="secondary" onClick={restart}>
                        인증 메일 다시 요청
                      </Button>
                    </>
                  ) : (
                    <>
                      <div className="email-summary">
                        <Mail size={17} aria-hidden />
                        <span>{email}</span>
                        <Button variant="quiet" onClick={restart}>
                          변경
                        </Button>
                      </div>
                      <p className="help">
                        메일의 인증 링크를 열면 해당 화면에서 다음 단계를 진행할
                        수 있습니다. 스팸 메일함도 확인해 주세요.
                      </p>
                      <Button variant="quiet" className="resend" onClick={send}>
                        인증 메일 다시 받기
                      </Button>
                      {error && <InlineError>{error}</InlineError>}
                    </>
                  )}
                </>
              ) : (
                <>
                  <p className="verified">
                    <CheckCircle2 size={17} aria-hidden />
                    이메일 인증 완료
                  </p>
                  <TextField
                    label="인증받은 이메일"
                    type="email"
                    autoComplete="email"
                    value={email}
                    disabled={true}
                    required
                  />
                  {register && (
                    <TextField
                      label="닉네임"
                      autoComplete="nickname"
                      value={nickname}
                      onChange={(event) => setNickname(event.target.value)}
                      required
                    />
                  )}
                  <PasswordField
                    label={register ? '비밀번호' : '새 비밀번호'}
                    autoComplete="new-password"
                    value={password}
                    onChange={(event) => setPassword(event.target.value)}
                    required
                  />
                  <PasswordField
                    label="비밀번호 확인"
                    autoComplete="new-password"
                    value={confirmation}
                    onChange={(event) => setConfirmation(event.target.value)}
                    required
                    error={
                      password !== confirmation ? error || undefined : undefined
                    }
                  />
                  {error && password === confirmation && (
                    <InlineError>{error}</InlineError>
                  )}
                </>
              )}
            </fieldset>
            {notice && (
              <p className="help" role="status">
                {notice}
              </p>
            )}
            {step !== 2 && (
              <Button
                className="full"
                type="submit"
                disabled={
                  busy ||
                  (step === 3 &&
                    (!email.trim() ||
                      !password ||
                      !confirmation ||
                      (register && !nickname.trim())))
                }
              >
                {busy ? '처리 중…' : step === 1 ? '인증 메일 받기' : title}
              </Button>
            )}
          </form>
        </>
      )}
      {register ? (
        <p className="auth-bottom">
          이미 가입하셨나요?<AppLink href="/login">로그인</AppLink>
        </p>
      ) : (
        <p className="auth-bottom">
          Google 계정은 Google에서 비밀번호를 관리합니다.
        </p>
      )}
    </AuthLayout>
  );
}
