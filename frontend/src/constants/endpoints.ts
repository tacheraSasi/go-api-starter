const API_URLS = {
  production: 'https://api.example.com',
  staging: 'https://api.example.com',
} as const;

const baseUrl = (env: keyof typeof API_URLS) => API_URLS[env];

/**
 * Below endpoints do not require authentication.
 * Auth token will not be sent in the request header.
 *  @readonly
 */
const OPEN_ENDPOINTS = ['/products'] as const;

export { API_URLS, baseUrl, OPEN_ENDPOINTS };
