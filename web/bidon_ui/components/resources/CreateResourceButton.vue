<template>
  <NavigationContainer v-if="permissions.create">
    <NuxtLink :to="`${resourcesPath}/new`">
      <Button :label="label" icon="pi pi-plus" class="p-button-success" />
    </NuxtLink>
  </NavigationContainer>
</template>

<script setup lang="ts">
import { useAuthStore } from "@/stores/AuthStore";

const props = defineProps<{
  resourcesPath: string;
  label: string;
}>();
const { getResourcePermissionsByPath } = useAuthStore();
const permissions: ResourcePermissions = await getResourcePermissionsByPath(
  props.resourcesPath,
);
</script>
