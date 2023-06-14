<template>
  <a href="_" @:click.prevent="handler">
    <Button label="Delete" icon="pi pi pi-trash" severity="danger" />
  </a>
</template>

<script setup lang="ts">
defineProps<{
  handler: () => void;
}>();
</script>
