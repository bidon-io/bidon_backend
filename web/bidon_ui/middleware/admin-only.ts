export default defineNuxtRouteMiddleware(() => {
  const authStore = useAuthStore();
  const user = authStore.currentUser as { isAdmin?: boolean } | null;

  // Redirect non-admins (or unauthenticated users) to home
  if (!user || user.isAdmin !== true) {
    return navigateTo("/");
  }
});
