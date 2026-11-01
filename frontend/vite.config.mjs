import { defineConfig } from 'vite'
import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { VitePWA } from 'vite-plugin-pwa'

const __dirname = path.dirname(fileURLToPath(import.meta.url))

export default defineConfig(async () => {
  const { visualizer } = await import('rollup-plugin-visualizer')

  return {
    plugins: [
      vue(),
      tailwindcss(),
      process.env.ANALYZE === 'true' && visualizer({
        open: true, // Автоматически открыть карту размеров в браузере после сборки
        filename: 'stats.html', // Имя файла с отчетом
        gzipSize: true, // Показать размер с учетом сжатия gzip
        brotliSize: true, // Показать размер с учетом сжатия brotli
      }),

      VitePWA({
        workbox: {
          globPatterns: ['**/*.{js,css,html,ico,png,svg,webmanifest}'],
        },
        registerType: 'autoUpdate',
        manifest: {
          name: 'Polyxena',
          short_name: 'PLXN',
          display: 'standalone',
          description: 'Polyxena - Opens Source Monitoring System',
          theme_color: '#000000',
          icons: [
            {
              src: '/pwa-192x192.png',
              sizes: '192x192',
              type: 'image/png',
              purpose: 'any maskable',
            },
            {
              src: '/pwa-512x512.png',
              sizes: '512x512',
              type: 'image/png',
              purpose: 'any maskable',
            },
          ],
        },
      }),
    ].filter(Boolean),
    build: {
      chunkSizeWarningLimit: 700,
    },
    resolve: {
      extensions: ['.mjs', '.js', '.ts', '.jsx', '.tsx', '.json', '.vue'],
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
  }
})