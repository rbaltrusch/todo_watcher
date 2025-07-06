import { fileURLToPath, URL } from 'node:url'

import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig(({mode})=> {
  
  const env = loadEnv(mode, process.cwd());

  return {
  plugins: [vue(), vueJsx(), vueDevTools()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },

  // proxies API requests to the backend server, configured via environment variable VITE_API_URL
  server: {
    proxy: {
      '/api': {
        target: env.VITE_API_URL || 'http://localhost:8080', // default to localhost:8080 if VITE_API_URL is not set
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '/api'), // optional if your path stays the same
      },
    },
  },
}})
