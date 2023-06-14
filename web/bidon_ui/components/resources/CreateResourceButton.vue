<template>
  <NavigationContainer>
    <NuxtLink :to="`${resourcesPath}/new`">
      <Button :label="label" icon="pi pi-plus" class="p-button-success" />
    </NuxtLink>
  </NavigationContainer>
</template>

<script setup lang="ts">
defineProps<{
  resourcesPath: string;
  label: string;
}>();
</script>
