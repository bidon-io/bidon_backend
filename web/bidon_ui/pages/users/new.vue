<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton />
    </NavigationContainer>
    <UserForm :value="resource" @submit="handleSubmit" />
  </PageContainer>
</template>

<script setup>
const resource = { settings: {} };
const resourcesPath = "/users";
const handleSubmit = useCreateResource({
  path: resourcesPath,
  message: "User created!",
});
</script>
