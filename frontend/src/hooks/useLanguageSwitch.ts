/* eslint-disable @typescript-eslint/no-misused-promises */
import { useCallback } from 'react';

import { useTranslation } from 'react-i18next';

type LocaleKeys = Record<string, string>;

type Locales = keyof typeof locales;

// writing this way helps in autocompletion
function localeKeys<T extends LocaleKeys>(arg: T): T {
  return arg;
}

const locales = localeKeys({
  EN: 'English',
  ES: 'Español',
});

const useLanguageSwitch = (): {
  activeLanguage: Locales;
  availableLocales: typeof locales;
  handleLanguageSwitch: (activeLocale: Locales) => void;
} => {
  const { i18n } = useTranslation();
  const locale = i18n?.language?.toUpperCase() as Locales;

  const handleLanguageSwitch = useCallback(
    async (code: Locales) => {
      await i18n.changeLanguage(code.toLowerCase());
    },
    [i18n],
  );

  return {
    activeLanguage: locale,
    availableLocales: locales,
    handleLanguageSwitch,
  };
};

export { useLanguageSwitch };
