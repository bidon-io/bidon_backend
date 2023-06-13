<template>
  <Toast />
  <PageContainer>
    <NavigationContainer>
      <NuxtLink to="/auction_configurations/">
        <Button label="Go back" icon="pi pi-arrow-left" severity="secondary" text />
      </NuxtLink>
    </NavigationContainer>
    <AuctionConfigurationForm :value="resource" @submit="handleSubmit" />
  </PageContainer>
</template>

<script setup>
import axios from "@/services/ApiService.js";
import { useToast } from "primevue/usetoast";
const resource = {};

const toast = useToast();
const handleSubmit = (event) => {
  axios
    .post("/auction_configurations", event)
    .then(async (response) => {
      const id = response.data.id;
      await navigateTo(`/auction_configurations/${id}`);

      toast.add({
        severity: "success",
        summary: "Success",
        detail: "Auction configuration created",
      });
    })
    .catch((error) => {
      console.error(error);
      toast.add({
        severity: "error",
        summary: "Error",
        detail: error.message,
      });
    });
};
</script>
