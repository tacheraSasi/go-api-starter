import { createFileRoute } from '@tanstack/react-router';
import { useTranslation } from 'react-i18next';

import { Seo } from 'src/components';
import { INTL_NAMESPACE } from 'src/i18n/config';
import homeStyles from 'src/styles/routes/home.module.scss';
import logger from 'src/utils/logger';

const Home = () => {
  const { t } = useTranslation([INTL_NAMESPACE.COMMON]);
  logger.info('ENV', process.env.NODE_ENV);
  return (
    <>
      <Seo title='Home' />
      <div className={homeStyles.container}>{t('common:greet')}</div>
    </>
  );
};

const Route = createFileRoute('/')({
  component: Home,
});

export { Route };
export default Home;
