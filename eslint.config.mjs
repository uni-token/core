import antfu from '@antfu/eslint-config'

export default antfu({
  vue: true,
  typescript: true,
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
})
