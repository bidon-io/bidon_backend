<template>
  <a
    v-if="permissions.delete"
    href="_"
    @:click.prevent="() => deleteHandle(id)"
  >
    <Button label="Delete" icon="pi pi pi-trash" severity="danger" />
  </a>
</template>

<script setup lang="ts">
import { useAuthStore } from "@/stores/AuthStore";

const props = defineProps<{
  id: string;
  path: string;
}>();

const deleteHandle = useDeleteResource({
  path: props.path,
  hook: async () => await navigateTo(props.path),
});

const { getResourcePermissionsByPath } = useAuthStore();
const permissions: ResourcePermissions = await getResourcePermissionsByPath(
  props.path,
);
</script>
