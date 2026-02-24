import { QueryClient } from '@tanstack/react-query';
import { ReactQueryDevtools } from '@tanstack/react-query-devtools';
import { Outlet, createRootRouteWithContext } from '@tanstack/react-router';
import { TanStackRouterDevtools } from '@tanstack/router-devtools';

import { SHOW_DEV_TOOLS } from 'src/constants/app-constants';
import { AuthContextProps } from 'src/context/authContext';

interface MyRouterContext {
  auth: AuthContextProps;
  queryClient: QueryClient;
}

const RootComponent = () => {
  return (
    <>
      <Outlet />
      {SHOW_DEV_TOOLS && (
        <>
          <ReactQueryDevtools buttonPosition='top-right' />
          <TanStackRouterDevtools position='bottom-right' />
        </>
      )}
    </>
  );
};

const Route = createRootRouteWithContext<MyRouterContext>()({
  component: RootComponent,
  notFoundComponent: () => {
    return <p>Page you are looking for does not exist</p>;
  },
});

export { Route };
