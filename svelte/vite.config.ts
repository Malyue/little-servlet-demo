import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [svelte()],
  resolve:{
    //设置路径别名
    alias: {
      '@': path.resolve(__dirname, './src'),
      },
  },
  server: {
		host: 'localhost',
		port: 5173,
		open: false,
		https: false,
		proxy: {
		  '/api': 'http://localhost:9090/',
		},
	  },
})
