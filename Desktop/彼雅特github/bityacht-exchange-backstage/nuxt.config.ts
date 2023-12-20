// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  modules: ["@nuxtjs/tailwindcss"],
  app: {
    head: {
      link: [
        { rel: "preconnect", href: "https://fonts.googleapis.com" },
        {
          rel: "preconnect",
          href: "https://fonts.gstatic.com",
          crossorigin: "",
        },
        {
          href: "https://fonts.googleapis.com/css2?family=Lato:wght@400;700&family=Noto+Sans+TC:wght@400;700&display=swap",
          rel: "stylesheet",
        },
        { rel: "icon", type: "image/svg+xml", href: "/favicon.svg" },
        { rel: "alternate icon", href: "/favicon.png" },
      ],
      title: "BitYacht",
    },
  },
  plugins: [
    { src: "@/plugins/antd", mode: "client" },
    { src: "@/plugins/vue-toastification", mode: "client" },
  ],
  runtimeConfig: {
    public: {
      apiURL: process.env.API_URL || "https://dev-exchange.skycloud.tw/api/v1",
    },
  },
  devtools: { enabled: true },
  experimental: {
    // WORKAROUND: https://github.com/nuxt-modules/i18n/issues/2177
    inlineSSRStyles: false,
    payloadExtraction: false,
  },
});
