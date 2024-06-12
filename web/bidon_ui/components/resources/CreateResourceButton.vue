<template>
  <NavigationContainer v-if="resource?.permissions?.create">
    <NuxtLink :to="`${resourcesPath}/new`">
      <Button :label="label" icon="pi pi-plus" class="p-button-success" />
    </NuxtLink>
  </NavigationContainer>
</template>

<script setup lang="ts">
import { camelize, singularize } from "inflection";

const props = defineProps<{
  resourcesPath: string;
  label: string;
}>();

const key = computed(() => {
  if (props.resourcesPath === "/v2/auction_configurations") {
    return "auctionConfigurationV2";
  }
  return camelize(singularize(props.resourcesPath.replace("/", "")), true);
});

const resources = useResources();
const resource = computed(() => resources.state[key.value] ?? {});
</script>
