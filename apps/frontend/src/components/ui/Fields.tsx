import { useId, useState } from 'react';
import type { InputHTMLAttributes, TextareaHTMLAttributes } from 'react';
import { Eye, EyeOff } from 'lucide-react';

type FieldProps = { label: string; help?: string; error?: string };
export function TextField({
  label,
  help,
  error,
  id: givenId,
  ...props
}: InputHTMLAttributes<HTMLInputElement> & FieldProps) {
  const generated = useId();
  const id = givenId || generated;
  return (
    <div className="field">
      <label htmlFor={id}>{label}</label>
      <input
        id={id}
        aria-invalid={!!error}
        aria-describedby={error || help ? `${id}-help` : undefined}
        {...props}
      />
      {(error || help) && (
        <p id={`${id}-help`} className={error ? 'field-error' : 'help'}>
          {error || help}
        </p>
      )}
    </div>
  );
}
export function PasswordField({
  label = '비밀번호',
  help,
  error,
  id: givenId,
  ...props
}: InputHTMLAttributes<HTMLInputElement> & Partial<FieldProps>) {
  const generated = useId();
  const id = givenId || generated;
  const [shown, setShown] = useState(false);
  return (
    <div className="field">
      <label htmlFor={id}>{label}</label>
      <div className="input-wrap">
        <input
          id={id}
          type={shown ? 'text' : 'password'}
          aria-invalid={!!error}
          aria-describedby={error || help ? `${id}-help` : undefined}
          {...props}
        />
        <button
          type="button"
          className="eye"
          aria-label={shown ? `${label} 숨기기` : `${label} 표시`}
          aria-pressed={shown}
          onClick={() => setShown(!shown)}
        >
          {shown ? (
            <EyeOff size={18} aria-hidden />
          ) : (
            <Eye size={18} aria-hidden />
          )}
        </button>
      </div>
      {(error || help) && (
        <p id={`${id}-help`} className={error ? 'field-error' : 'help'}>
          {error || help}
        </p>
      )}
    </div>
  );
}
export function TextAreaField({
  label,
  help,
  error,
  id: givenId,
  ...props
}: TextareaHTMLAttributes<HTMLTextAreaElement> & FieldProps) {
  const generated = useId();
  const id = givenId || generated;
  return (
    <div className="field">
      <label htmlFor={id}>{label}</label>
      <textarea
        id={id}
        aria-invalid={!!error}
        aria-describedby={error || help ? `${id}-help` : undefined}
        {...props}
      />
      {(error || help) && (
        <p id={`${id}-help`} className={error ? 'field-error' : 'help'}>
          {error || help}
        </p>
      )}
      <span className="char-count">{String(props.value ?? '').length}자</span>
    </div>
  );
}
