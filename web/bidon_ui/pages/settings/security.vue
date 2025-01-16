<template>
  <PageContainer>
    <PasswordForm :submit-error="error" @submit="handleSubmit" />
  </PageContainer>
</template>

<script setup>
const resourcePath = "/settings/password";
const error = ref(null);

const handleSubmit = useUpdateResource({
  path: resourcePath,
  message: "Password Updated!",
  onError: async (e) => (error.value = e),
});
</script>
