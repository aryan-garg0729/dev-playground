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
export const HexIcon = (p: IconProps) => (
  <Icon {...p}>
    <path d="M12 3 20 7.5v9L12 21l-8-4.5v-9Z" />
  </Icon>
);
export const SnowflakeIcon = (p: IconProps) => (
  <Icon {...p}>
    <path d="M12 3v18M4.2 7.5l15.6 9M4.2 16.5l15.6-9" />
  </Icon>
);
export const BoltIcon = (p: IconProps) => (
  <Icon {...p}>
    <path d="M13 3 5 14h6l-1 7 8-11h-6Z" />
  </Icon>
);
export const RotateIcon = (p: IconProps) => (
  <Icon {...p}>
    <path d="M4 12a8 8 0 1 0 2.3-5.7M4 4v4h4" />
  </Icon>
);
export const BulbIcon = (p: IconProps) => (
  <Icon {...p}>
    <path d="M9 18h6M10 21h4M12 3a6 6 0 0 0-3.5 10.9c.6.5 1 1.2 1 2.1h5c0-.9.4-1.6 1-2.1A6 6 0 0 0 12 3Z" />
  </Icon>
);
export const EyeIcon = (p: IconProps) => (
  <Icon {...p}>
    <path d="M2 12s3.6-7 10-7 10 7 10 7-3.6 7-10 7S2 12 2 12Z" />
    <circle cx="12" cy="12" r="3" />
  </Icon>
);
export const ShieldIcon = (p: IconProps) => (
  <Icon {...p}>
    <path d="M12 3 4 6v6c0 5 3.5 8 8 9 4.5-1 8-4 8-9V6Z" />
  </Icon>
);
export const FlameOffIcon = (p: IconProps) => (
  <Icon {...p}>
    <path d="M12 3c.5 3 3.5 4.5 3.5 8.5a3.5 3.5 0 0 1-7 0c0-1.6.7-2.7 1.5-3.5M3 3l18 18" />
  </Icon>
);
export const CalendarIcon = (p: IconProps) => (
  <Icon {...p}>
    <rect x="4" y="5" width="16" height="16" rx="2" />
    <path d="M4 10h16M9 3v4M15 3v4" />
  </Icon>
);
export const DiceIcon = (p: IconProps) => (
  <Icon {...p}>
    <rect x="3" y="3" width="11" height="11" rx="2" />
    <rect x="10" y="10" width="11" height="11" rx="2" />
    <circle cx="7" cy="7" r="1" fill="currentColor" />
    <circle cx="15.5" cy="15.5" r="1" fill="currentColor" />
  </Icon>
);
export const MedalIcon = (p: IconProps) => (
  <Icon {...p}>
    <path d="M7 3l3 6M17 3l-3 6" />
    <circle cx="12" cy="15" r="6" />
    <circle cx="12" cy="15" r="3" />
  </Icon>
);
