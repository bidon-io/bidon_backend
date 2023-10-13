import { API_URL } from "~/constants";
import { camelizeKeys } from "humps";
import { $fetch } from "ofetch";

export const $apiFetch = $fetch.create({
  baseURL: `${API_URL}api`,
  headers: {
    "X-Bidon-App": "web",
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
  async onResponseError({ response }) {
    if (response.status === 401) {
      useAuthStore().setUnauthorized();
      await navigateTo("/login");
    }
  },
});
