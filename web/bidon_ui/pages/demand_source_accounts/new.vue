<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton :path="resourcesPath" />
    </NavigationContainer>
    <DemandSourceAccountForm :value="resource" @submit="handleSubmit" />
  </PageContainer>
</template>

<script setup>
const resource = {};
const resourcesPath = "/demand_source_accounts";
const handleSubmit = useCreateResource({
  path: resourcesPath,
  message: "Demand source account created!",
});
</script>
