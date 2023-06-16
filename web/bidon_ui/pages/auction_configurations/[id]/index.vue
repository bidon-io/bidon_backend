<template>
  <Toast />
  <ConfirmDialog />
  <PageContainer>
    <NavigationContainer>
      <GoBackButton />
      <DestroyButton :handler="() => deleteHandle(id)" />
      <EditButton :path="`${resourcesPath}/${id}/edit`" />
    </NavigationContainer>
    <ResourceCard title="Auction Configuration" :fields="fields" :resource="resource" />
  </PageContainer>
</template>

<script setup>
import axios from "@/services/ApiService.js";

const route = useRoute();
const id = route.params.id;
const resourcesPath = "/auction_configurations";
const deleteHandle = useDeleteResource(resourcesPath, async () => await navigateTo(resourcesPath));

const response = await axios.get(`${resourcesPath}/${id}`);
const resource = response.data;

const fields = [
  { label: "ID", key: "id" },
  { label: "App", key: "app_id", type: "link", link: `/apps/${resource.app_id}` },
  { label: "Name", key: "name" },
  { label: "Ad type", key: "ad_type" },
  { label: "Price floor", key: "pricefloor" },
  { label: "Rounds", key: "rounds", type: "textarea" },
];
</script>
