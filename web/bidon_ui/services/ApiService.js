import axios from "axios";
import { API_URL } from "@/constants/index.js";

export default axios.create({
  baseURL: API_URL,
  data: {},
});
