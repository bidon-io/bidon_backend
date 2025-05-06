<template>
  <Button
    label="Go back"
    icon="pi pi-arrow-left"
    severity="secondary"
    text
    @click.prevent="goBack"
  />
</template>

<script setup>
const router = useRouter();

const goBack = () => {
  router.back();
};
</script>
