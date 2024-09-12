import { defineConfig, presetTypography, presetUno, transformerDirectives, transformerVariantGroup } from 'unocss'

export default defineConfig({
  shortcuts: [
    ['btn-1', 'border-none outline-none p-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 cursor-pointer'],

    ['text-secondary', 'text-sm color-gray-500'],

    ['code-block', 'p-4 rounded-md bg-gray-100 color-gray-900 overflow-x-auto'],

    ['top-bar', 'h-80px mb-4'],
  ],
  presets: [
    presetUno(),
    presetTypography(),
  ],
  transformers: [
    transformerDirectives(),
    transformerVariantGroup(),
  ],
  safelist: [],
  theme: {
    fontFamily: {
      body: ['"JetBrains Mono"', 'monospace'].join(','),
    },
  },
})
