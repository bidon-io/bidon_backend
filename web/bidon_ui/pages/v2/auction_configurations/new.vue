<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton :path="resourcesPath" />
    </NavigationContainer>
    <AuctionConfigurationsForm
      v-if="isReady"
      :value="resource"
      @submit="handleSubmit"
    />
  </PageContainer>
</template>

<script setup>
import { useAsyncState } from "@vueuse/core";
import axios from "@/services/ApiService";

const route = useRoute();
const resourcesPath = "/v2/auction_configurations";
const cloneId = route.query.clone;

const { state: resource, isReady } = useAsyncState(async () => {
  if (cloneId) {
    try {
      const response = await axios.get(`${resourcesPath}/${cloneId}`);
      const sourceResource = response.data;

      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      const { id, publicUid, auctionKey, ...clonedResource } = sourceResource;
      return {
        ...clonedResource,
        name: `${sourceResource.name} (Copy)`,
        isDefault: false, // Don't clone as default
      };
    } catch (error) {
      console.error("Failed to fetch source auction configuration:", error);
      return {};
    }
  }
  return {};
}, {});

const handleSubmit = useCreateResource({
  path: resourcesPath,
  message: "Auction configuration created!",
});
</script>
