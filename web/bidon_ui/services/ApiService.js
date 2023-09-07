import axios from "axios";
import { camelizeKeys, decamelizeKeys } from "humps";
import { API_URL } from "@/constants/index.js";
// import { useAuthStore } from "@/stores/AuthStore.js";

const api = axios.create({
  baseURL: `${API_URL}api`,
  data: {},
});

// Axios middleware to convert all api responses to camelCase and logout user if token is invalid
api.interceptors.response.use(
  (response) => {
    if (
      response.data &&
      response.headers["content-type"].includes("application/json")
    ) {
      response.data = camelizeKeys(response.data);
    }

    return response;
  },
  (error) => {
    if (error.response.status === 401) {
      // useAuthStore().logout();
      // window.location.reload();
    }
  }
);

// Axios middleware to convert all api requests to snake_case and add authorization header
api.interceptors.request.use((config) => {
  const newConfig = { ...config };

  // const { accessToken } = useAuthStore();
  // if (accessToken) {
  //   newConfig.headers["Authorization"] = `Bearer ${accessToken}`;
  // }

  if (newConfig.headers["Content-Type"] === "multipart/form-data")
    return newConfig;

  if (config.params) {
    newConfig.params = decamelizeKeys(config.params);
  }

  if (config.data) {
    newConfig.data = decamelizeKeys(config.data);
  }

  return newConfig;
});

export default api;
