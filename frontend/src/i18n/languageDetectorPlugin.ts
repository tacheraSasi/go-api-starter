/* eslint-disable import/no-import-module-exports */
import { STORAGE_KEYS } from 'src/constants';
import LocalStorageUtil from 'src/utils/local-storage';
import logger from 'src/utils/logger';

/**
 * Custom plugin to detect the preferred language of the user
 * @see https://www.i18next.com/misc/creating-own-plugins#languagedetector
 */
const languageDetectorPlugin = {
  async: true,

  cacheUserLanguage(language: string) {
    try {
      // save a user's language choice in Async storage
      LocalStorageUtil.setItem(STORAGE_KEYS.LANGUAGE, language.toLowerCase());
    } catch (error) {
      logger.error('i18n:', 'Error setting language');
    }
  },

  // eslint-disable-next-line consistent-return
  detect(callback: (lang: string) => void) {
    try {
      const language: string = LocalStorageUtil.getItem(STORAGE_KEYS.LANGUAGE) ?? 'en';
      if (language) {
        // if language was stored before, use this language in the app
        return callback(language.toLowerCase());
      }
      // if language was not stored yet, use "en"
      return callback('en');
    } catch (error) {
      logger.error('i18n:', 'Error reading language');
    }
  },

  init: () => {},
  type: 'languageDetector',
};

export { languageDetectorPlugin };
