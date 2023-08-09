<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton />
      <DestroyButton :handler="() => deleteHandle(id)" />
      <EditButton :path="`${resourcesPath}/${id}/edit`" />
    </NavigationContainer>
    <ResourceCard title="Country" :fields="fields" :resource="resource" />
  </PageContainer>
</template>

<script setup>
import axios from "@/services/ApiService.js";
import { ResourceCardFields } from "@/constants";

const route = useRoute();
const id = route.params.id;
const resourcesPath = "/countries";
const deleteHandle = useDeleteResource({
  path: resourcesPath,
  hook: async () => await navigateTo(resourcesPath),
});

const response = await axios.get(`${resourcesPath}/${id}`);
const resource = response.data;

const fields = [
  ResourceCardFields.Id,
  { label: "Human Name", key: "humanName" },
  { label: "Alpha 2 Code", key: "alpha2Code" },
  { label: "Alpha 3 Code", key: "alpha3Code" },
];
</script>
