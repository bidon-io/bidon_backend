<template>
  <NuxtLink v-if="permissions.update" :to="`${path}/${id}/edit`">
    <Button label="Edit" icon="pi pi-pencil" />
  </NuxtLink>
</template>

<script setup lang="ts">
import { useAuthStore } from "@/stores/AuthStore";
const props = defineProps<{
  id: string;
  path: string;
}>();

const { getResourcePermissionsByPath } = useAuthStore();
const permissions: ResourcePermissions = await getResourcePermissionsByPath(
  props.path,
);
</script>
