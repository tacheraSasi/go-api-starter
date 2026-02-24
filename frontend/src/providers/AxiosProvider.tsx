/* eslint-disable no-param-reassign */
import { ReactNode, createContext, useMemo } from 'react';

import Axios, { AxiosInstance, AxiosError, InternalAxiosRequestConfig, AxiosResponse } from 'axios';
import { isEmpty } from 'radash';

import { API_URLS, baseUrl, OPEN_ENDPOINTS } from 'src/constants/endpoints';

const AxiosContext = createContext<AxiosInstance>(Axios.create({}));

// Abstracted interceptor logic for requests
const configureRequestInterceptor = (axiosInstance: AxiosInstance, token?: string) => {
  axiosInstance.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
      const isOpenEndpoint = config.url ? OPEN_ENDPOINTS.some((endpoint) => config.url!.includes(endpoint)) : false;
      if (!isEmpty(token) && !isOpenEndpoint) {
        config.headers.Authorization = token;
      } else {
        delete config.headers.Authorization;
      }
      return config;
    },
    (error: AxiosError) => Promise.reject(error),
  );
};

// Abstracted interceptor logic for responses
const configureResponseInterceptor = (axiosInstance: AxiosInstance) => {
  axiosInstance.interceptors.response.use(
    (response) => response,
    (error: AxiosError) => {
      // Change according to your apis
      const errorObject = {
        error: error.response?.data,
        status: error.response?.status,
      };
      return Promise.reject(errorObject);
    },
  );
};

function axios(env: keyof typeof API_URLS) {
  return Axios.create({
    baseURL: baseUrl(env),
    headers: {
      'Content-Type': 'application/json',
    },
    timeout: 5000,
  });
}

const AxiosProvider = ({
  token,
  children,
  env,
}: {
  children: ReactNode;
  env: keyof typeof API_URLS;
  token?: string;
}) => {
  const axiosInstance: AxiosInstance = useMemo(() => {
    const getAxiosInstance = axios(env);

    // Apply interceptors
    configureRequestInterceptor(getAxiosInstance, token);
    configureResponseInterceptor(getAxiosInstance);

    return getAxiosInstance;
  }, [env, token]);

  return <AxiosContext.Provider value={axiosInstance}>{children}</AxiosContext.Provider>;
};

export { AxiosContext, AxiosProvider };
export type { AxiosError, AxiosResponse };
