<template>
  <Toast />
  <PageContainer>
    <NavigationContainer>
      <GoBackButton :path="resourcesPath" />
    </NavigationContainer>
    <SegmentForm v-if="isReady" :value="resource" @submit="handleSubmit" />
  </PageContainer>
</template>

<script setup>
import { useAsyncState } from "@vueuse/core";
import axios from "@/services/ApiService";

const route = useRoute();
const id = route.params.id;
const resourcesPath = "/segments";

const { state: resource, isReady } = useAsyncState(async () => {
  const response = await axios.get(`${resourcesPath}/${id}`);
  return response.data;
});

const handleSubmit = useFormSubmit(resourcesPath, "Segment Updated!");
</script>
