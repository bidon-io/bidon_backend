import { defineStore } from "pinia";
import { camelizeKeys } from "humps";
import axios from "axios";
import authorizedApi from "@/services/ApiService.js";
import { useToast } from "primevue/usetoast";
import { API_URL } from "@/constants/index.js";

export const useAuthStore = defineStore("authStore", () => {
  const api = axios.create({ baseURL: API_URL });
  const toast = useToast();
  const router = useRouter();

  const localStorageUser = localStorage.getItem("user");
  const user = ref(
    localStorageUser ? camelizeKeys(JSON.parse(localStorageUser)) : null,
  );
  const accessToken = ref(localStorage.getItem("accessToken") || null);
  const permissions = ref([]);

  async function login(email, password) {
    api
      .post("/auth/login", { email, password })
      .then((response) => {
        user.value = camelizeKeys(response.data.user);
        accessToken.value = response.data["access_token"];
        localStorage.setItem("user", JSON.stringify(user.value));
        localStorage.setItem("accessToken", accessToken.value);
        router.push("/");
      })
      .catch((error) => {
        toast.add({
          severity: "error",
          summary: "Error",
          detail:
            error?.response?.data?.error?.message || "Something went wrong",
        });
      });
  }

  async function logout() {
    user.value = null;
    accessToken.value = null;
    localStorage.removeItem("user");
    localStorage.removeItem("accessToken");
    router.push("/login");
  }

  async function getResourcesPermissions() {
    if (permissions.value.length > 0) return permissions.value;

    const response = await authorizedApi.get("/permissions");
    permissions.value = response.data;

    return permissions.value;
  }

  async function getResourcePermissionsByPath(path) {
    const permissions = await getResourcesPermissions();
    return (
      permissions.find((permission) => permission.path === path)?.actions || {}
    );
  }

  return {
    user,
    accessToken,
    login,
    logout,
    getResourcesPermissions,
    getResourcePermissionsByPath,
  };
});
