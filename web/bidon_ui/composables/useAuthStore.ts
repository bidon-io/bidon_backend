import { defineStore } from "pinia";
import { useToast } from "primevue/usetoast";

export const useAuthStore = defineStore("authStore", () => {
  const currentUser = ref(null);
  const isAuthorized = ref(true);

  function setCurrentUser(user) {
    currentUser.value = user;
    isAuthorized.value = true;
  }

  function setUnauthorized() {
    currentUser.value = null;
    isAuthorized.value = false;
  }

  function setAuthorized() {
    isAuthorized.value = true;
  }

  const toast = useToast();

  async function login(email: string, password: string) {
    try {
      await $fetch("/auth/login", {
        method: "POST",
        body: { email, password },
      });

      setAuthorized();

      return navigateTo("/");
    } catch (error) {
      toast.add({
        severity: "error",
        summary: "Error",
        detail: error?.data?.error?.message || "Something went wrong",
      });
    }
  }

  async function logout() {
    await $fetch("/auth/logout", {
      method: "POST",
    });

    setUnauthorized();

    return navigateTo("/login");
  }

  return {
    currentUser,
    isAuthorized,
    setCurrentUser,
    setUnauthorized,
    setAuthorized,
    login,
    logout,
  };
});
