import path from 'path';

import { TanStackRouterVite } from '@tanstack/router-vite-plugin';
import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';

// https://vitejs.dev/config/
export default defineConfig({
  // css: {
  //   modules: {
  //     generateScopedName: () => {
  //       if (process.env.NODE_ENV === 'production') {
  //         // Use a custom hashing function or any strategy you prefer
  //         // This is a very simplistic approach for demonstration
  //         return 'style--[hash:base64:5]';
  //       }
  //       // For development, use a more descriptive name
  //       return `[name]__[local]___[hash:base64:5]`;
  //     },
  //   },
  // },
  plugins: [react(), TanStackRouterVite()],
  resolve: {
    alias: {
      src: path.resolve(__dirname, './src'),
    },
  },
});
