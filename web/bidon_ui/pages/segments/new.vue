<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton :path="resourcesPath" />
    </NavigationContainer>
    <SegmentForm :value="resource" @submit="handleSubmit" />
  </PageContainer>
</template>

<script setup>
const resource = {};
const resourcesPath = "/segments";
const handleSubmit = useCreateResource({
  path: resourcesPath,
  message: "Segment created!",
});
</script>
