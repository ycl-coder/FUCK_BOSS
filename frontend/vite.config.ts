import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
      '@/api': path.resolve(__dirname, './src/api'),
      '@/features': path.resolve(__dirname, './src/features'),
      '@/shared': path.resolve(__dirname, './src/shared'),
    },
  },
  server: {
    port: 8000,
    host: true, // 允许外部访问
    allowedHosts: [
      'localhost',
      '.ngrok-free.dev', // 允许所有 ngrok 域名
      '.ngrok.io', // 允许旧版 ngrok 域名
      '.ngrok.app', // 允许新版 ngrok 域名
    ],
    proxy: {
      // Proxy gRPC Web requests to backend
      '/api': {
        target: 'http://localhost:50051',
        changeOrigin: true,
        secure: false,
      },
    },
  },
  build: {
    outDir: 'dist',
    sourcemap: true,
  },
})
