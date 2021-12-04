import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import viteSvgIcons from 'vite-plugin-svg-icons'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import { resolve } from 'path'

function resolvePath(dir) {
  return resolve(__dirname, dir)
}

// https://vitejs.dev/config/
export default defineConfig({
  base: '/web',
  plugins: [
    vue(),
    viteSvgIcons({
      // Specify the icon folder to be cached
      iconDirs: [resolvePath('src/assets/icon')],
      // Specify symbolId format
      symbolId: 'icon-[name]'
    }),
    Components({
      useSource: true,
      resolvers: [ElementPlusResolver({
        importStyle: 'sass'
      })]
    })
  ],
  resolve: {
    alias: {
      '@': resolvePath('src') // 设置 `@` 指向 `src` 目录
    }
  },
  css: {
    preprocessorOptions: {
      scss: {
        charset: false,
        additionalData: `@use "@/assets/css/var.scss" as *;`
      }
    }
  },
  optimizeDeps: {
    include: [
      'axios', 'vue', 'vue-router', 'element-plus'
    ]
  },
  server: {
    port: 3001, // 设置服务启动端口号
    open: true, // 设置服务启动时是否自动打开浏览器
    cors: true // 允许跨域
    // 设置代理，根据我们项目实际情况配置
    // proxy: {
    //   '/api': {
    //     target: 'http://xxx.xxx.xxx.xxx:x000',
    //     changeOrigin: true,
    //     secure: false,
    //     rewrite: (path) => path.replace('/api/', '/')
    //   }
    // },
  },
  build: {
    target: 'modules',
    outDir: '../web',
    assetsDir: 'assets',
    sourcemap: false,
    minify: 'terser', // 压缩混淆
    manifest: false,
    chunkSizeWarningLimit: 2000,
    terserOptions: {
      compress: {
        drop_console: true,
        drop_debugger: true
      }
    },
    rollupOptions: {
      output: {
        // manualChunks: {
        //   vue: ['vue', 'vue-router', 'vuex'],
        //   'element-plus': ['element-plus'],
        //   echarts: ['echarts']
        // }
        // entryFileNames: '[name].js',
        // chunkFileNames: '[name].js',
        // assetFileNames: '[name].[ext]',
        manualChunks(id) {
          if (id.includes('/node_modules/')) {
            const expansions = []
            if (expansions.some(exp => id.includes(`/node_modules/${exp}`))) {
              return 'expansion'
            } else {
              return 'vendor'
            }
          }
        }
      }
    }
  }
})

