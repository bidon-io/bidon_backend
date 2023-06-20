import { useToast } from "primevue/usetoast";
import axios from "@/services/ApiService.js";

export default function (
  path,
  method,
  message,
  callback = () => {
    /* no operation function */
  }
) {
  const toast = useToast();
  const handleSubmit = (event) => {
    axios[method](path, event)
      .then(async (response) => {
        const id = response.data.id;
        await callback(id);

        toast.add({
          severity: "success",
          summary: "Success",
          detail: message,
        });
      })
      .catch((error) => {
        console.error(error);
        toast.add({
          severity: "error",
          summary: "Error",
          detail: error.error.message,
        });
      });
  };
  return handleSubmit;
}
