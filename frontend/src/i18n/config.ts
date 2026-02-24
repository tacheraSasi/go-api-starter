/* eslint-disable @typescript-eslint/ban-ts-comment */
/* eslint-disable @typescript-eslint/no-floating-promises */
/* eslint-disable @typescript-eslint/no-unsafe-assignment */
/* eslint-disable @typescript-eslint/no-var-requires */
/* eslint-disable global-require */
import i18n, { CustomTypeOptions } from 'i18next';
import { initReactI18next } from 'react-i18next';

import { languageDetectorPlugin } from './languageDetectorPlugin';
import enCommon from './locales/en/common.json';
import esCommon from './locales/es/common.json';

const INTL_NAMESPACE = {
  COMMON: 'common',
} as const;

//* Add types for below in types/i18next.d.ts file
const resources = Object.freeze({
  en: {
    [INTL_NAMESPACE.COMMON]: enCommon,
  } as CustomTypeOptions['resources'],
  es: {
    [INTL_NAMESPACE.COMMON]: esCommon,
  } as CustomTypeOptions['resources'],
});

i18n
  // .use(FsBackend)
  .use(initReactI18next)
  //   @ts-ignore
  .use(languageDetectorPlugin)
  .init({
    // backend: {
    //   loadPath: `${DocumentDirectoryPath}/locales/{{lng}}/{{ns}}.json`,
    // },
    compatibilityJSON: 'v3',
    fallbackLng: 'en',
    interpolation: {
      escapeValue: false,
    },
    lng: 'en',
    react: {
      useSuspense: false,
    },
    resources,
  });

export { resources, INTL_NAMESPACE };
export default i18n;
