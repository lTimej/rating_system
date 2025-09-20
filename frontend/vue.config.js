const { defineConfig } = require('@vue/cli-service')

module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    port: 3000,
    host: "198.19.249.53",
    proxy: {
      '/api': {
        target: 'http://198.19.249.53:8080',
        changeOrigin: false,
        pathRewrite: {
          '^': '' // 重写路径：去掉路径中开头的''
        }
      }
    }
  }
})
