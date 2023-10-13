export default defineNuxtRouteMiddleware(async (to) => {
  const authStore = useAuthStore();

  if (authStore.isAuthorized && !authStore.currentUser) {
    try {
      const currentUser = await $apiFetch("users/me");
      authStore.setCurrentUser(currentUser);
    } catch (e) {
      if (e.response.status === 401) {
        authStore.setUnauthorized();
        return navigateTo("/login");
      }
    }
  }

  const normalizedPath = to.path.replace(/\/$/, "");
  if (["/login", "/signup"].includes(normalizedPath) && authStore.currentUser) {
    return navigateTo("/");
  }
  if (
    !["/login", "/signup"].includes(normalizedPath) &&
    !authStore.currentUser
  ) {
    return navigateTo("/login");
  }
});
