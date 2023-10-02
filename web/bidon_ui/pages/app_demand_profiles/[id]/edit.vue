<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton :path="resourcePath" />
    </NavigationContainer>
    <AppDemandProfileForm
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
const resourcePath = `/app_demand_profiles/${id}`;

const { state: resource, isReady } = useAsyncState(async () => {
  const response = await axios.get(resourcePath);
  return response.data;
});

const error = ref(null);
const handleSubmit = useUpdateResource({
  path: resourcePath,
  message: "App Demand Profile Updated!",
  onError: async (e) => (error.value = e),
});
</script>
