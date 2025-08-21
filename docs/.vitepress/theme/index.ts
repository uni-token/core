import DefaultTheme from 'vitepress/theme'
import Layout from './Layout.vue'
import 'virtual:uno.css'

/** @type {import('vitepress').Theme} */
export default {
  extends: DefaultTheme,
  Layout,
}
