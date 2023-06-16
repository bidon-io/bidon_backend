<template>
  <Toast />
  <PageContainer>
    <NavigationContainer> <GoBackButton" /> </NavigationContainer>
    <AuctionConfigurationForm v-if="isReady" :value="resource" @submit="handleSubmit" />
  </PageContainer>
</template>

<script setup>
import { useAsyncState } from "@vueuse/core";
import axios from "@/services/ApiService";

const route = useRoute();
const id = route.params.id;
const resourcesPath = "/auction_configurations";

const { state: resource, isReady } = useAsyncState(async () => {
  const response = await axios.get(`${resourcesPath}/${id}`);
  return response.data;
});

const handleSubmit = useFormSubmit(resourcesPath, "Auction Configuration Updated!");
</script>
