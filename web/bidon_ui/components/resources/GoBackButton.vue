<template>
  <NuxtLink :to="path">
    <Button label="Go back" icon="pi pi-arrow-left" severity="secondary" text />
  </NuxtLink>
</template>

<script setup lang="ts">
defineProps<{
  path: string;
}>();
</script>
