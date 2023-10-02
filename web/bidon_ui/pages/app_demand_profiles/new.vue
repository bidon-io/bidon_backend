<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton :path="resourcesPath" />
    </NavigationContainer>
    <AppDemandProfileForm
      :value="resource"
      :submit-error="error"
      @submit="handleSubmit"
    />
  </PageContainer>
</template>

<script setup>
const resource = {};
const resourcesPath = "/app_demand_profiles";

const error = ref(null);
const handleSubmit = useCreateResource({
  path: resourcesPath,
  message: "App demand profile created!",
  onError: async (e) => (error.value = e),
});
</script>
