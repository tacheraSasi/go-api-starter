import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { RouterProvider, ErrorComponent, createRouter } from '@tanstack/react-router';

import { AuthProvider } from './context/authContext';
import { useAuth } from './hooks';
import { AxiosProvider } from './providers/AxiosProvider';
import { routeTree } from './routeTree.gen';

const queryClient = new QueryClient();

const router = createRouter({
  context: {
    auth: undefined!, // This will be set after we wrap the app in an AuthProvider
    queryClient,
  },
  defaultErrorComponent: ({ error }) => <ErrorComponent error={error} />,
  defaultPendingComponent: () => (
    <div>
      <p>Loading...</p>
    </div>
  ),
  defaultPreload: 'intent',
  // Since we're using React Query, we don't want loader calls to ever be stale
  // This will ensure that the loader is always called when the route is preloaded or visited
  defaultPreloadStaleTime: 0,
  routeTree,
});

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}

const InnerApp = () => {
  const auth = useAuth();
  return <RouterProvider router={router} context={{ auth }} />;
};

const App = () => {
  return (
    <AuthProvider>
      <AxiosProvider env='production'>
        <QueryClientProvider client={queryClient}>
          <InnerApp />
        </QueryClientProvider>
      </AxiosProvider>
    </AuthProvider>
  );
};

// eslint-disable-next-line react-refresh/only-export-components
export { queryClient };
export default App;
