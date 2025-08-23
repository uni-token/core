import antfu from '@antfu/eslint-config'

export default antfu({
  vue: true,
  typescript: true,
  markdown: true,
  formatters: {
    markdown: true,
    css: true,
    slidev: {
      files: [
        '**/slides.md',
      ],
    },
  },
  ignores: [
    'frontend/public/**',
    '**/.venv/**',
    '**/node_modules/**',
    '**/dist/**',
  ],
}, {
  rules: {
    'vue/custom-event-name-casing': 'off',
  },
}, {
  files: ['sdk/node/src/example.ts'],
  rules: {
    'no-console': 'off', // Allow console logs in example files
  },
})
