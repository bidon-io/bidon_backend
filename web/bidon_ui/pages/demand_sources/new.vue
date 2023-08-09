<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton />
    </NavigationContainer>
    <DemandSourceForm :value="resource" @submit="handleSubmit" />
  </PageContainer>
</template>

<script setup>
const resource = { settings: {} };
const resourcesPath = "/demand_sources";
const handleSubmit = useCreateResource({
  path: resourcesPath,
  message: "Demand Source created!",
});
</script>
