<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton />
      <DestroyButton :handler="() => deleteHandle(id)" />
      <EditButton :path="`${resourcesPath}/${id}/edit`" />
    </NavigationContainer>
    <ResourceCard title="App" :fields="fields" :resource="resource" />
  </PageContainer>
</template>

<script setup>
import axios from "@/services/ApiService.js";

const route = useRoute();
const id = route.params.id;
const resourcesPath = "/apps";
const deleteHandle = useDeleteResource({
  path: resourcesPath,
  hook: async () => await navigateTo(resourcesPath),
});

const response = await axios.get(`${resourcesPath}/${id}`);
const resource = response.data;

const fields = [
  { label: "ID", key: "id" },
  { label: "Platform Id", key: "platformId" },
  { label: "Human Name", key: "humanName" },
  { label: "Package Name", key: "packageName" },
  { label: "User", key: "userId" },
  { label: "App Key", key: "appKey" },
  { label: "Settings", key: "settings", type: "textarea" },
];
</script>
