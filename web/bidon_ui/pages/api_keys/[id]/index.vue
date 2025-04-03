<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton :path="resourcesPath" />
      <DestroyButton
        v-if="resource._permissions.delete"
        :id="id"
        :path="resourcesPath"
      />
    </NavigationContainer>
    <ResourceCard
      title="API Credential"
      :fields="fields"
      :resource="resource"
    />
  </PageContainer>
</template>

<script setup lang="ts">
import { $apiFetch } from "~/utils/$apiFetch";

const route = useRoute();
const id = route.params.id;

const resourcesPath = "/api_keys";

const resource = await $apiFetch(`${resourcesPath}/${id}`);
const endpoint = `${window.location.origin}/api`;

const fields = [
  { label: "ID", key: "id" },
  { label: "Endpoint", type: "static", value: endpoint, copyable: true },
  { label: "API Key", key: "value", copyable: true },
  { label: "Last Accessed At", key: "lastAccessedAt" },
];
</script>
