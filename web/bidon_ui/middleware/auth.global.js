import { useAuthStore } from "~/stores/AuthStore";

export default defineNuxtRouteMiddleware((to) => {
  const router = useRouter();
  const { user } = useAuthStore();

  const normilizedPath = to.path.replace(/\/$/, "");
  if (["/login", "/signup"].includes(normilizedPath) && user) {
    return router.push("/");
  }
  if (!["/login", "/signup"].includes(normilizedPath) && !user) {
    return router.push("/login");
  }
});
