import { defineStore } from "pinia";
import { useAsyncState } from "@vueuse/core";

interface Resources {
  [key: string]: {
    key: string;
    permissions: {
      read: boolean;
      create: boolean;
    };
  };
}

export const useResources = defineStore("resources", () => {
  const { state, isReady, isLoading } = useAsyncState(
    $apiFetch<Resources>("/rest/resources"),
    {},
  );

  return { state, isReady, isLoading };
});
