import { API_URL } from "~/constants";
import { camelizeKeys } from "humps";
import { useAuthStore } from "~/stores/AuthStore";

export const $apiFetch = $fetch.create({
  baseURL: `${API_URL}api`,
  onRequest({ options }) {
    const auth = useAuthStore();
    if (auth.accessToken) {
      options.headers = {
        ...options.headers,
        Authorization: `Bearer ${auth.accessToken}`,
      };
    }
  },
  onResponse({ response }) {
    if (
      response._data &&
      response.headers.get("Content-Type")?.includes("application/json")
    ) {
      response._data = camelizeKeys(response._data, (key, convert, options) =>
        key[0] === "_" ? key : convert(key, options),
      );
    }
  },
  onResponseError({ response }) {
    if (response.status === 401) {
      useAuthStore().logout();
      window.location.reload();
    }
  },
});
