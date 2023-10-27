import axios from "axios";
import { camelizeKeys, decamelizeKeys } from "humps";
import { API_URL } from "@/constants/index.js";

const api = axios.create({
  baseURL: `${API_URL}api`,
  data: {},
  headers: {
    "X-Bidon-App": "web",
  },
});

// Axios middleware to convert all api responses to camelCase and logout user if token is invalid
api.interceptors.response.use(
  (response) => {
    if (
      response.data &&
      response.headers["content-type"].includes("application/json")
    ) {
      response.data = camelizeKeys(response.data, (key, convert, options) =>
        key[0] === "_" ? key : convert(key, options),
      );
    }

    return response;
  },
  (error) => {
    if (error.response.status === 401) {
      useAuthStore().setUnauthorized();
      return navigateTo("/login");
    }
    throw error;
  },
);

// Axios middleware to convert all api requests to snake_case and add authorization header
api.interceptors.request.use((config) => {
  const newConfig = { ...config };

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
