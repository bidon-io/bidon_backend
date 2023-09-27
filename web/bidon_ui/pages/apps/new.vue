<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton :path="resourcesPath" />
    </NavigationContainer>
    <AppForm :value="resource" :submit-error="error" @submit="handleSubmit" />
  </PageContainer>
</template>

<script setup>
const resource = { settings: {} };
const resourcesPath = "/apps";

const error = ref(null);
const handleSubmit = useCreateResource({
  path: resourcesPath,
  message: "App created!",
  onError: async (e) => (error.value = e),
});
</script>
