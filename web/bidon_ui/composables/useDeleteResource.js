import { useConfirm } from "primevue/useconfirm";
import { useToast } from "primevue/usetoast";
import axios from "@/services/ApiService.js";

export default function (path, onDelete) {
  const confirm = useConfirm();
  const toast = useToast();

  async function deleteResource(id, callback) {
    await axios.delete(`${path}/${id}`);
    await callback();
  }

  const deleteHandle = (id) => {
    confirm.require({
      message: "Do you want to delete this record?",
      header: "Delete Confirmation",
      icon: "pi pi-info-circle",
      acceptClass: "p-button-danger",
      accept: () => {
        deleteResource(id, () => {
          onDelete(id);
          toast.add({ severity: "info", summary: "Success", detail: "Record deleted", life: 3000 });
        });
      },
    });
  };
  return deleteHandle;
}
