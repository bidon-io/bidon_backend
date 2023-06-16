<template>
  <Toast />
  <PageContainer>
    <NavigationContainer> <GoBackButton /> </NavigationContainer>
    <AuctionConfigurationForm v-if="isReady" :value="resource" @submit="handleSubmit" />
  </PageContainer>
</template>

<script setup>
import { useAsyncState } from "@vueuse/core";
import axios from "@/services/ApiService";

const route = useRoute();
const id = route.params.id;
const resourcePath = `/auction_configurations/${id}`;

const { state: resource, isReady } = useAsyncState(async () => {
  const response = await axios.get(resourcePath);
  return response.data;
});

const handleSubmit = useFormSubmit(resourcePath, "patch", "Auction Configuration Updated!");
</script>
