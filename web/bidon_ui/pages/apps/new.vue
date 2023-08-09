<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton :path="resourcesPath" />
    </NavigationContainer>
    <AppForm :value="resource" @submit="handleSubmit" />
  </PageContainer>
</template>

<script setup>
const resource = { settings: {} };
const resourcesPath = "/apps";
const handleSubmit = useCreateResource({
  path: resourcesPath,
  message: "App created!",
});
</script>
