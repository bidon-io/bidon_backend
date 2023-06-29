import axios from "axios";
import { camelizeKeys, decamelizeKeys } from "humps";
import { API_URL } from "@/constants/index.js";

const api = axios.create({
  baseURL: API_URL,
  data: {},
});

// Axios middleware to convert all api responses to camelCase
api.interceptors.response.use((response) => {
  if (response.data && response.headers["content-type"].includes("application/json")) {
    response.data = camelizeKeys(response.data);
  }

  return response;
});

// Axios middleware to convert all api requests to snake_case
api.interceptors.request.use((config) => {
  const newConfig = { ...config };

  if (newConfig.headers["Content-Type"] === "multipart/form-data") return newConfig;

  if (config.params) {
    newConfig.params = decamelizeKeys(config.params);
  }

  if (config.data) {
    newConfig.data = decamelizeKeys(config.data);
  }

  return newConfig;
});

export default api;
