export let API_URL;

if (import.meta.env.VITE_APP_ENV === "production") {
  API_URL = "/";
} else {
  API_URL = "http://localhost:1323/";
}
