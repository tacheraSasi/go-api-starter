/**
 * @readonly
 * Used to identify cache for queries
 */
const CACHE_KEYS = {
  LOGIN: 'login',
} as const;

/**
 * @readonly
 * Used for debugging purposes during mutations
 */
const MUTATION_KEYS = {
  LOGIN: 'user-login',
  LOGOUT: 'user-logout',
} as const;

/**
 * readonly
 * These caches will clear once user is logged out
 *
 * @see useLogout
 */
const LOGOUT_CACHE_CLEAR = [] as const;

export { CACHE_KEYS, MUTATION_KEYS, LOGOUT_CACHE_CLEAR };
