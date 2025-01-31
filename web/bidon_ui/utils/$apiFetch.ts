import { API_URL } from "~/constants";
import { camelizeKeys } from "humps";
import { $fetch } from "ofetch";

const baseURL = `${API_URL}api`;

export const $apiFetch = $fetch.create({
  baseURL: baseURL,
  headers: {
    "X-Bidon-App": "web",
  },
  onRequest(context) {
    // baseURL does not prepend if request matches it
    // https://github.com/unjs/ufo/blob/496140d0abcc3c3636409eb403c4916fade58203/src/utils.ts#L231

    if (
      typeof context.request === "string" &&
      context.request.startsWith("/api_keys")
    ) {
      context.request = `${baseURL}${context.request}`;
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
  async onResponseError({ response }) {
    if (response.status === 401) {
      useAuthStore().setUnauthorized();
      await navigateTo("/login");
    }
  },
});
