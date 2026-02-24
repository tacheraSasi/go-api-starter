import LocalStorageUtil from './local-storage';

interface Storage {
  getItem: (key: string) => string | null;
  removeItem: (key: string) => void;
  setItem: (key: string, value: string) => void;
}

/**
 * Storage adapter used with @tanstack/react-query persistor
 */
function createStoragePersistor(): Storage {
  return {
    getItem: (storageKey: string): string | null => {
      return LocalStorageUtil.getItem(storageKey);
    },
    removeItem: (storageKey: string) => {
      LocalStorageUtil.clearItem(storageKey);
    },
    setItem: (storageKey: string, value: string) => {
      LocalStorageUtil.setItem(storageKey, value);
    },
  };
}

export type { Storage };
export default createStoragePersistor;
