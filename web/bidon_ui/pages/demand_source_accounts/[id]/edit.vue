<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton />
    </NavigationContainer>
    <DemandSourceAccountForm v-if="isReady" :value="resource" @submit="handleSubmit" />
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

const handleSubmit = useUpdateResource({ path: resourcePath, message: "Demand source account updated!" });
</script>
