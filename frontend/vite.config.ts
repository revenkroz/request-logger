import path from 'node:path'
import { defineConfig } from 'vite'
import Vue from '@vitejs/plugin-vue'
import UnoCSS from 'unocss/vite'
import { webfontDownload } from 'vite-plugin-webfont-dl'

export default defineConfig({
  resolve: {
    alias: {
      '~/': `${path.resolve(__dirname, 'src')}/`,
    },
  },

  build: {
    rollupOptions: {
      input: {
        main: path.resolve(__dirname, 'index.html'),
      },
    },
  },

  plugins: [
    Vue({
      include: [/\.vue$/],
    }),
    UnoCSS(),
    webfontDownload([
      'https://fonts.googleapis.com/css2?family=JetBrains+Mono&display=swap',
    ]),
  ],
})
