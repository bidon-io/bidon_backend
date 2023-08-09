<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton :path="resourcesPath" />
    </NavigationContainer>
    <AppDemandProfileForm :value="resource" @submit="handleSubmit" />
  </PageContainer>
</template>

<script setup>
const resource = {};
const resourcesPath = "/app_demand_profiles";
const handleSubmit = useCreateResource({
  path: resourcesPath,
  message: "App demand profile created!",
});
</script>
