// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  ssr:false,
  modules: ['@nuxtjs/tailwindcss','@pinia/nuxt','@nuxtjs/i18n','@ant-design-vue/nuxt','dayjs-nuxt'],
  // nitro:{
  //   devProxy:{
  //     "/dev-api": {
  //       target: 'http://10.66.13.110:9000',
  //       changeOrigin:true,
  //       prependPath:true,
  //     }
  //   }
  // },
  // routeRules:{
  //   "/dev-api/**":{
  //     proxy:'http://10.66.13.110:9000/**'
  //   }
  // },  
  build: {
    transpile: [/echarts/],
  },
  app:{
    baseURL:process.env.APP_BASE_URL ? process.env.APP_BASE_URL : '/',
    head: {
      link: [
        { rel: "icon", type: "image/svg+xml", href: "/favicon.svg" },
        { rel: "alternate icon", href: "/favicon.png" },
      ],
      title: "BitYacht",
    },
  },
  runtimeConfig: {
    public: {
      apiBase: 'http://10.66.13.110:9100/api/v1',
    }
  },
  devtools: { enabled: true },
  i18n: {
    defaultLocale: "zh-TW",
    locales: [
      { code: "zh-TW", iso: "zh-TW", file: "zh-TW.js" },
      { code: "en-US", iso: "en-US", file: "en-US.js" },
    ],
  },
  pinia: {
    autoImports: [
      // automatically imports `defineStore`
      'defineStore', // import { defineStore } from 'pinia'
      ['defineStore', 'definePiniaStore'], // import { defineStore as definePiniaStore } from 'pinia'
    ],
  },
})
