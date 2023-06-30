<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton />
      <DestroyButton :handler="() => deleteHandle(id)" />
      <EditButton :path="`${resourcesPath}/${id}/edit`" />
    </NavigationContainer>
    <ResourceCard
      title="App Deman Profile"
      :fields="fields"
      :resource="resource"
    />
  </PageContainer>
</template>

<script setup>
import axios from "@/services/ApiService.js";

const route = useRoute();
const id = route.params.id;
const resourcesPath = "/app_demand_profiles";
const deleteHandle = useDeleteResource({
  path: resourcesPath,
  hook: async () => await navigateTo(resourcesPath),
});

const response = await axios.get(`${resourcesPath}/${id}`);
const resource = response.data;

const fields = [
  { label: "ID", key: "id" },
  { label: "App", key: "appId", type: "link", link: `/apps/${resource.appId}` },
  { label: "Demand Source", key: "demandSourceId" },
  {
    label: "Account",
    key: "accountId",
    link: `/demand_source_accounts/${resource.accountId}`,
  },
  { label: "Data", key: "data", type: "textarea" },
  { label: "Account Type", key: "accountType" },
];
</script>
