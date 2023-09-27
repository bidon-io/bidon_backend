<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton :path="resourcesPath" />
    </NavigationContainer>
    <DemandSourceAccountForm
      :value="resource"
      :submit-error="error"
      @submit="handleSubmit"
    />
  </PageContainer>
</template>

<script setup>
const resource = {};
const resourcesPath = "/demand_source_accounts";

const error = ref(null);
const handleSubmit = useCreateResource({
  path: resourcesPath,
  message: "Demand source account created!",
  onError: async (e) => (error.value = e),
});
</script>
