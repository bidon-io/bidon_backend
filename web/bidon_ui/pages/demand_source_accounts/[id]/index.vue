<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton />
      <DestroyButton :handler="() => deleteHandle(id)" />
      <EditButton :path="`${resourcesPath}/${id}/edit`" />
    </NavigationContainer>
    <ResourceCard
      title="Demand Source Account"
      :fields="fields"
      :resource="resource"
    />
  </PageContainer>
</template>

<script setup>
import axios from "@/services/ApiService.js";
import { ResourceCardFields } from "@/constants";

const route = useRoute();
const id = route.params.id;
const resourcesPath = "/demand_source_accounts";
const deleteHandle = useDeleteResource({
  path: resourcesPath,
  hook: async () => await navigateTo(resourcesPath),
});

const response = await axios.get(`${resourcesPath}/${id}`);
const resource = response.data;

const fields = [
  ResourceCardFields.Id,
  ResourceCardFields.User,
  { label: "Type", key: "type" },
  ResourceCardFields.DemandSource,
  { label: "IsBidding", key: "isBidding" },
  { label: "Extra", key: "extra", type: "textarea" },
];
</script>
