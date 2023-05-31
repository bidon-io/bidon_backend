// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  ssr: false,
  css: [
    "primevue/resources/themes/lara-light-blue/theme.css",
    "primevue/resources/primevue.css",
    "primeicons/primeicons.css",
  ],
  modules: ["@nuxtjs/tailwindcss"],
  build: {
    transpile: ["primevue"],
  },
  routeRules: {
    "/": { redirect: "/auction_configurations" },
  },
});
