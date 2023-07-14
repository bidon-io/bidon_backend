<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton />
    </NavigationContainer>
    <LineItemForm :value="resource" @submit="handleSubmit" />
  </PageContainer>
</template>

<script setup>
const resource = {};
const resourcesPath = "/line_items";
const handleSubmit = useCreateResource({
  path: resourcesPath,
  message: "Line Item created!",
});
</script>
