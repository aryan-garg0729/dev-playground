import type { ButtonHTMLAttributes, ReactNode } from 'react';

type Variant = 'primary' | 'secondary' | 'tutor';

const variants: Record<Variant, string> = {
  // One primary per screen: accent fill with accent-ink text.
  primary: 'h-[34px] px-4 bg-accent text-accent-ink enabled:hover:brightness-105',
  // Light theme adds a 1px border (see Button README).
  secondary: 'h-[34px] px-4 bg-surface-raised text-text [[data-theme=light]_&]:border [[data-theme=light]_&]:border-border',
  tutor: 'h-10 w-full bg-transparent text-accent-text border border-dashed border-accent-line',
};

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant;
}

export function Button({ variant = 'secondary', className = '', type = 'button', ...rest }: ButtonProps) {
  return (
    <button
      type={type}
      className={`text-control inline-flex items-center justify-center gap-2 rounded-md disabled:cursor-not-allowed disabled:text-text-faint ${variants[variant]} ${className}`}
      {...rest}
    />
  );
}

interface IconButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  'aria-label': string;
  children: ReactNode;
}

/** 32px square icon-only button; always labelled. */
export function IconButton({ className = '', type = 'button', ...rest }: IconButtonProps) {
  return (
    <button
      type={type}
      className={`inline-flex h-8 w-8 items-center justify-center rounded-md border border-border text-text-secondary hover:bg-surface-raised disabled:text-text-faint ${className}`}
      {...rest}
    />
  );
}
