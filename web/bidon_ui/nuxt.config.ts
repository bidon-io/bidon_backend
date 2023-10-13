// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  alias: {
    assets: "/<rootDir>/assets",
  },
  ssr: false,
  css: [
    "primevue/resources/themes/lara-light-blue/theme.css",
    "primevue/resources/primevue.css",
    "primeicons/primeicons.css",
  ],
  components: [
    {
      path: "~/components",
      pathPrefix: false,
    },
  ],
  modules: ["@nuxtjs/tailwindcss", "@pinia/nuxt", "@vee-validate/nuxt"],
  build: {
    transpile: ["primevue"],
  },
  routeRules: {
    "/auth/**": { proxy: "http://localhost:1323/auth/**" },
    "/api/**": { proxy: "http://localhost:1323/api/**" },
  },
});
