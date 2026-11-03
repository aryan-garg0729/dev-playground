// Stroke icons on a 24px grid, 2px stroke, round caps/joins; color comes from currentColor.
import type { SVGProps } from 'react';

type IconProps = SVGProps<SVGSVGElement> & { size?: number };

function Icon({ size = 16, strokeWidth = 2, children, ...rest }: IconProps) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth={strokeWidth}
      strokeLinecap="round"
      strokeLinejoin="round"
      aria-hidden="true"
      {...rest}
    >
      {children}
    </svg>
  );
}

export const CodeIcon = (p: IconProps) => (
  <Icon strokeWidth={2.5} {...p}>
    <path d="m8 7-5 5 5 5M16 7l5 5-5 5" />
  </Icon>
);
export const SunIcon = (p: IconProps) => (
  <Icon {...p}>
    <circle cx="12" cy="12" r="4" />
    <path d="M12 2v2M12 20v2M4.9 4.9l1.4 1.4M17.7 17.7l1.4 1.4M2 12h2M20 12h2M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4" />
  </Icon>
);
export const MoonIcon = (p: IconProps) => (
  <Icon {...p}>
    <path d="M21 12.8A9 9 0 1 1 11.2 3a7 7 0 0 0 9.8 9.8Z" />
  </Icon>
);
export const ChevronLeftIcon = (p: IconProps) => (
  <Icon strokeWidth={2.5} {...p}>
    <path d="m15 18-6-6 6-6" />
  </Icon>
);
export const ChevronRightIcon = (p: IconProps) => (
  <Icon strokeWidth={2.5} {...p}>
    <path d="m9 18 6-6-6-6" />
  </Icon>
);
export const CheckIcon = (p: IconProps) => (
  <Icon strokeWidth={3} {...p}>
    <path d="M20 6 9 17l-5-5" />
  </Icon>
);
export const CrossIcon = (p: IconProps) => (
  <Icon strokeWidth={3} {...p}>
    <path d="M18 6 6 18M6 6l12 12" />
  </Icon>
);
export const LockIcon = (p: IconProps) => (
  <Icon {...p}>
    <rect x="5" y="11" width="14" height="10" rx="2" />
    <path d="M8 11V7a4 4 0 0 1 8 0v4" />
  </Icon>
);
export const PanelLeftIcon = (p: IconProps) => (
  <Icon {...p}>
    <rect x="3" y="4" width="18" height="16" rx="2" />
    <path d="M9 4v16M13 10l2 2-2 2" />
  </Icon>
);
export const FlameIcon = (p: IconProps) => (
  <Icon {...p}>
    <path d="M12 3c.5 3 3.5 4.5 3.5 8.5a3.5 3.5 0 0 1-7 0c0-1.6.7-2.7 1.5-3.5.2 1.3.8 2 1.5 2.5C11 8.5 11 5.5 12 3Z" />
    <path d="M6.5 13.5a5.5 5.5 0 0 0 11 0" />
  </Icon>
);
export const ChestIcon = (p: IconProps) => (
  <Icon {...p}>
    <rect x="3" y="8" width="18" height="12" rx="2" />
    <path d="M3 13h18M8 8V6a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2M11 13v2h2v-2" />
  </Icon>
);
export const StarIcon = (p: IconProps) => (
  <Icon {...p}>
    <path d="m12 3 2.6 5.6 6 .7-4.5 4.1 1.2 6L12 16.4 6.7 19.4l1.2-6-4.5-4.1 6-.7Z" />
  </Icon>
);
