import { useToast } from "primevue/usetoast";
import axios from "@/services/ApiService.js";

export default function (
  resourcesPath,
  message,
  callback = () => {
    /* no operation function */
  }
) {
  const toast = useToast();
  const handleSubmit = (event) => {
    axios
      .post(resourcesPath, event)
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
          detail: error.message,
        });
      });
  };
  return handleSubmit;
}
