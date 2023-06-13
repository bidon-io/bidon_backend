<template>
  <Toast />
  <ConfirmDialog />
  <PageContainer>
    <NavigationContainer>
      <NuxtLink to="/auction_configurations/">
        <Button label="Go back" icon="pi pi-arrow-left" severity="secondary" text />
      </NuxtLink>
      <a href="_" @:click.prevent="deleteHandle(id)">
        <Button label="Delete" icon="pi pi pi-trash" severity="danger" />
      </a>
      <NuxtLink :to="`/auction_configurations/${id}/edit`">
        <Button label="Edit" icon="pi pi-pencil" />
      </NuxtLink>
    </NavigationContainer>
    <ResourceCard :fields="fields" :resource="resource" />
  </PageContainer>
</template>

<script setup>
import axios from "@/services/ApiService.js";
const route = useRoute();
const id = route.params.id;
const deleteHandle = useDeleteResource(
  "auction_configurations",
  async () => await navigateTo("/auction_configurations")
);

const response = await axios.get(`auction_configurations/${id}`);
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
