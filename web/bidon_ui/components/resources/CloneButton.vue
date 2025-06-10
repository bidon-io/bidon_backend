<template>
  <NuxtLink
    v-if="resource?.permissions?.create"
    :to="`${path}/new?clone=${id}`"
  >
    <Button label="Clone" icon="pi pi-copy" severity="secondary" />
  </NuxtLink>
</template>

<script setup lang="ts">
import { camelize, singularize } from "inflection";

const props = defineProps<{
  id: string;
  path: string;
}>();

const key = computed(() => {
  if (props.path === "/v2/auction_configurations") {
    return "auctionConfigurationV2";
  }
  return camelize(singularize(props.path.replace("/", "")), true);
});

const resources = useResources();
const resource = computed(() => resources.state[key.value] ?? {});
</script>
