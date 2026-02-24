import { createContext, type ReactNode, useMemo, useState } from 'react';

import { isEmpty } from 'radash';

import logger from 'src/utils/logger';

interface AuthContextProps {
  authToken: string | null;
  isAuthenticated: boolean;
  logout: () => void;
  removeAuthToken: () => void;
  setAuthToken: (token: string | null) => void;
}

const AuthContext = createContext<AuthContextProps | null>(null);

const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [authToken, setToken] = useState<string | null>(''); // Replace with a storage mechanism of your choice
  logger.debug('authToken', authToken);

  const value = useMemo(() => {
    const isAuthenticated = !isEmpty(authToken);

    const setAuthToken = (token: string | null) => {
      setToken(token ?? '');
    };

    const removeAuthToken = () => {
      setAuthToken(null);
    };

    const logout = () => {
      removeAuthToken();
    };

    return {
      authToken,
      isAuthenticated,
      logout,
      removeAuthToken,
      setAuthToken,
    };
  }, [authToken, setToken]);

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export type { AuthContextProps };
export { AuthProvider, AuthContext };
