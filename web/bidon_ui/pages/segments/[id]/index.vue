<template>
  <Toast />
  <ConfirmDialog />
  <PageContainer>
    <NavigationContainer>
      <GoBackButton :path="resourcesPath" />
      <DestroyButton :handler="() => deleteHandle(id)" />
      <EditButton :path="`${resourcesPath}/${id}/edit`" />
    </NavigationContainer>
    <ResourceCard :fields="fields" :resource="resource" />
  </PageContainer>
</template>

<script setup>
import axios from "@/services/ApiService.js";

const route = useRoute();
const id = route.params.id;
const resourcesPath = "/segments";
const deleteHandle = useDeleteResource(resourcesPath, async () => await navigateTo(resourcesPath));

const response = await axios.get(`${resourcesPath}/${id}`);
const resource = response.data;

const fields = [
  { label: "ID", key: "id" },
  { label: "Name", key: "name" },
  { label: "Description", key: "description" },
  { label: "Filters", key: "filters" },
  { label: "Enabled", key: "enabled" },
  { label: "App", key: "app_id", type: "link", link: `/apps/${resource.app_id}` },
];
</script>
