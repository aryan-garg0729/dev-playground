interface ImportMetaEnv {
  /** Origin of the sandbox runner page (separate origin; dev default http://localhost:5174). */
  readonly VITE_SANDBOX_ORIGIN?: string;
  /** Build id sent with error reports (e.g. a git sha); "dev" when unset. */
  readonly VITE_RELEASE?: string;
}
