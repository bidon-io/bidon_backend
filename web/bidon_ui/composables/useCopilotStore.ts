import { defineStore } from "pinia";
import { useStorage } from "@vueuse/core";

// Persist threadId in localStorage so conversations survive reloads
export const useCopilotStore = defineStore("copilotStore", () => {
  const threadId = useStorage<string | null>("copilot.threadId", null);

  function setThreadId(id: string | null) {
    threadId.value = id;
  }

  function reset() {
    threadId.value = null;
  }

  return {
    threadId,
    setThreadId,
    reset,
  };
});
