<template>
  <Toast />
  <PageContainer>
    <NavigationContainer>
      <GoBackButton />
    </NavigationContainer>
    <SegmentForm :value="resource" @submit="handleSubmit" />
  </PageContainer>
</template>

<script setup>
const resource = {};
const resourcesPath = "/segments";
const handleSubmit = useFormSubmit(resourcesPath, "post", "Segment created!", async (id) =>
  navigateTo(`${resourcesPath}/${id}`)
);
</script>
