import type { AnchorHTMLAttributes, ButtonHTMLAttributes } from 'react';
import { go } from '../../lib/navigation';

type Variant = 'primary' | 'secondary' | 'quiet' | 'danger';
export function Button({
  variant = 'primary',
  className = '',
  type = 'button',
  ...props
}: ButtonHTMLAttributes<HTMLButtonElement> & { variant?: Variant }) {
  return (
    <button
      type={type}
      className={`${variant === 'quiet' ? 'text-link' : 'btn'} ${variant} ${className}`}
      {...props}
    />
  );
}
export function AppLink({
  href = '/',
  onClick,
  ...props
}: AnchorHTMLAttributes<HTMLAnchorElement>) {
  return (
    <a
      href={href}
      onClick={(event) => {
        onClick?.(event);
        if (
          !event.defaultPrevented &&
          event.button === 0 &&
          !event.metaKey &&
          !event.ctrlKey &&
          !event.shiftKey &&
          !event.altKey &&
          !props.target
        ) {
          event.preventDefault();
          go(href);
        }
      }}
      {...props}
    />
  );
}
export function LinkButton({
  variant = 'primary',
  className = '',
  ...props
}: AnchorHTMLAttributes<HTMLAnchorElement> & { variant?: Variant }) {
  return (
    <AppLink
      className={`${variant === 'quiet' ? 'text-link' : 'btn'} ${variant} ${className}`}
      {...props}
    />
  );
}
