<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton :path="resourcesPath" />
    </NavigationContainer>
    <LineItemForm
      :value="resource"
      :submit-error="error"
      @submit="handleSubmit"
    />
  </PageContainer>
</template>

<script setup>
const resource = {};
const resourcesPath = "/line_items";

const error = ref(null);
const handleSubmit = useCreateResource({
  path: resourcesPath,
  message: "Line Item created!",
  onError: async (e) => (error.value = e),
});
</script>
