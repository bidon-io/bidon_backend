<template>
  <Button
    label="Go back"
    icon="pi pi-arrow-left"
    severity="secondary"
    text
    @click.prevent="router.back()"
  />
</template>

<script setup>
const router = useRouter();
</script>
