import { useToast } from "primevue/usetoast";
import axios from "@/services/ApiService.js";

export default function ({
  path,
  method,
  message,
  onSuccess = async () => {
    /* no operation function */
  },
  onError = async () => {
    /* no operation function */
  },
}) {
  const toast = useToast();
  const handleSubmit = async (event) => {
    try {
      const response = await axios[method](path, event);

      const id = response.data.id;
      await onSuccess(id);
      toast.add({
        severity: "success",
        summary: "Success",
        detail: message,
        life: 3000,
      });
    } catch (error) {
      console.error(error);
      await onError(error);
    }
  };
  return handleSubmit;
}
