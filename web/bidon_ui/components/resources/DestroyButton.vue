<template>
  <a href="_" @:click.prevent="() => deleteHandle(id)">
    <Button label="Delete" icon="pi pi pi-trash" severity="danger" />
  </a>
</template>

<script setup lang="ts">
const props = defineProps<{
  id: string;
  path: string;
}>();

const deleteHandle = useDeleteResource({
  path: props.path,
  hook: async () => await navigateTo(props.path),
});
</script>
