<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton :path="resourcePath" />
    </NavigationContainer>
    <DemandSourceAccountForm
      v-if="isReady"
      :value="resource"
      :submit-error="error"
      @submit="handleSubmit"
    />
  </PageContainer>
</template>

<script setup>
import { useAsyncState } from "@vueuse/core";
import axios from "@/services/ApiService";

const route = useRoute();
const id = route.params.id;
const resourcePath = `/demand_source_accounts/${id}`;

const { state: resource, isReady } = useAsyncState(async () => {
  const response = await axios.get(resourcePath);
  return response.data;
});

const error = ref(null);
const handleSubmit = useUpdateResource({
  path: resourcePath,
  message: "Demand source account updated!",
  onError: async (e) => (error.value = e),
});
</script>
