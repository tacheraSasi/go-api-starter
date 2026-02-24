import { type ReactNode } from 'react';

// import { Helmet } from 'react-helmet-async';

interface SeoProps {
  children?: ReactNode;
  description?: string;
  title: string;
}

/**
 * Adds title to browser tab amongst other Seo tags
 *
 * Note: In React 19, the `Helmet` package is not needed.
 * @see https://react.dev/blog/2024/12/05/react-19#support-for-metadata-tags
 */
const Seo = ({ title, description, children }: SeoProps) => {
  return (
    <>
      <title>{title}</title>
      <meta name='title' content={title} />
      <meta property='og:title' name='og:title' content={title} />
      {description ? (
        <>
          <meta name='description' content={description} />
          <meta property='og:description' name='og:description' content={description} />
        </>
      ) : null}
      {children ?? null}
    </>
  );
};

export default Seo;
