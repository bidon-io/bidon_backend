<template>
  <NuxtLink :to="path">
    <Button label="Edit" icon="pi pi-pencil" />
  </NuxtLink>
</template>

<script setup lang="ts">
defineProps<{
  path: string;
}>();
</script>
