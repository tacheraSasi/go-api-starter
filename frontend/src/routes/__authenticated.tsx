/* eslint-disable @typescript-eslint/no-throw-literal */
import { createFileRoute, redirect } from '@tanstack/react-router';

export const Route = createFileRoute('/__authenticated')({
  beforeLoad: ({ location }) => {
    const isAuthenticated = false;
    // If the user is logged out, redirect them to the login page
    if (!isAuthenticated) {
      throw redirect({
        search: {
          redirect: location.href,
        },
        to: '/login',
      });
    }

    // Otherwise, return the user in context
    return {};
  },
});
