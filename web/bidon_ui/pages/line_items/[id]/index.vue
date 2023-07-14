<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton />
      <DestroyButton :handler="() => deleteHandle(id)" />
      <EditButton :path="`${resourcesPath}/${id}/edit`" />
    </NavigationContainer>
    <ResourceCard title="Line Item" :fields="fields" :resource="resource" />
  </PageContainer>
</template>

<script setup>
import axios from "@/services/ApiService.js";
import { ResourceCardFields } from "@/constants";

const route = useRoute();
const id = route.params.id;
const resourcesPath = "/line_items";
const deleteHandle = useDeleteResource({
  path: resourcesPath,
  hook: async () => await navigateTo(resourcesPath),
});

const response = await axios.get(`${resourcesPath}/${id}`);
const resource = response.data;

const fields = [
  ResourceCardFields.Id,
  ResourceCardFields.HumanName,
  ResourceCardFields.App,
  ResourceCardFields.BidFloor,
  ResourceCardFields.AdType,
  { key: "format", label: "Format" },
  ResourceCardFields.DemandSourceAccount,
  ResourceCardFields.AccountType,
  { key: "code", label: "Code" },
  { key: "extra", label: "Extra", type: "textarea" },
];
</script>
