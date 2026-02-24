import { useContext } from 'react';

import { AxiosContext } from 'src/providers/AxiosProvider';

function useAxios() {
  const axiosInstance = useContext(AxiosContext);

  if (!axiosInstance) {
    throw new Error('useAxios must be used within an AxiosProvider');
  }

  return axiosInstance;
}

export { useAxios };
