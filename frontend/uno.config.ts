import { defineConfig, presetTypography, presetUno, transformerDirectives, transformerVariantGroup } from 'unocss'

export default defineConfig({
  shortcuts: [
    ['btn-1', 'border-none outline-none p-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 cursor-pointer'],
    ['btn-2', 'border-none outline-none p-2 bg-slate-500 text-white rounded-md hover:bg-slate-600 cursor-pointer'],

    ['code-block', 'p-4 rounded-md bg-gray-100 color-gray-900'],
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
