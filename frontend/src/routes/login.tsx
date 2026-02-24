import { createFileRoute } from '@tanstack/react-router';

import { Seo } from 'src/components';

const Login = () => {
  return (
    <>
      <Seo title='Login' />
      <p>Login</p>
    </>
  );
};

const Route = createFileRoute('/login')({
  component: Login,
});

export { Route };
export default Login;
