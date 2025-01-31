<template>
  <NavigationContainer v-if="resource?.permissions?.create">
    <Button
      :disabled="buttonDisabled"
      class="p-button-success"
      label="Generate Token"
      icon="pi pi-cog"
      @click="createApiToken"
    ></Button>
  </NavigationContainer>
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup lang="ts">
import { $apiFetch } from "~/utils/$apiFetch";
import { useToast } from "primevue/usetoast";

const toast = useToast();

const columns = [
  { field: "id", header: "ID" },
  { field: "lastAccessedAt", header: "Last Accessed At" },
];
const resourcesPath = "/api_keys";

const resources = useResources();
const resource = computed(() => resources.state.apiKey ?? {});

const buttonDisabled = ref(false);
const createApiToken = async () => {
  buttonDisabled.value = true;

  const apiKey = await $apiFetch("/api_keys", {
    method: "POST",
  }).catch((error) => {
    console.error(error);
    toast.add({
      severity: "error",
      summary: `${error.response.status} ${error.response.statusText}`,
      detail: error.response?.data?.error?.message,
    });
  });

  buttonDisabled.value = false;

  if (!apiKey) {
    return;
  }

  navigateTo(`/api_keys/${apiKey.id}`);
};
</script>
