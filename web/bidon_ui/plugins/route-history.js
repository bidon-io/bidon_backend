export default defineNuxtPlugin((nuxtApp) => {
  const routeHistory = useState("route-history", () => []);
  const MAX_HISTORY_LENGTH = 5;

  nuxtApp.vueApp.use(({}) => {
    const router = useRouter();

    router.beforeEach((_to, from) => {
      if (from.fullPath) {
        routeHistory.value.unshift({
          path: from.path,
          fullPath: from.fullPath,
          query: { ...from.query },
          timestamp: Date.now(),
        });

        if (routeHistory.value.length > MAX_HISTORY_LENGTH) {
          routeHistory.value.pop();
        }
      }
    });
  });

  return {
    provide: {
      routeHistory: () => ({
        getPrevious: () => routeHistory.value[0] || null,
        getPreviousFullPath: () => routeHistory.value[0]?.fullPath || null,
        getPreviousQuery: () => routeHistory.value[0]?.query || {},
        getAll: () => [...routeHistory.value],
      }),
    },
  };
});
