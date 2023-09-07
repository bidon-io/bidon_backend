import { useAuthStore } from "~/stores/AuthStore";

export default defineNuxtRouteMiddleware((to) => {
  const router = useRouter();
  const { user } = useAuthStore();

  if (["/login", "/signup"].includes(to.path) && user) {
    return router.push("/");
  }
  if (!["/login", "/signup"].includes(to.path) && !user) {
    return router.push("/login");
  }
});
