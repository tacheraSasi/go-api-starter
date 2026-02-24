/// <reference types="./vite-env-override.d.ts" />
/// <reference types="vite/client" />

interface ImportMetaEnv {
  APP_ENVIRONMENT: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
