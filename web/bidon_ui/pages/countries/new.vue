<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton />
    </NavigationContainer>
    <CountryForm :value="resource" @submit="handleSubmit" />
  </PageContainer>
</template>

<script setup>
const resource = { settings: {} };
const resourcesPath = "/countries";
const handleSubmit = useCreateResource({
  path: resourcesPath,
  message: "Country created!",
});
</script>
