<template>
  <nav class="mt-6">
    <NuxtLink
      v-for="resource in resources.state"
      :key="resource.key"
      :to="resourcePath(resource.key)"
      :class="[
        'flex items-center mt-4 px-6 py-2 text-gray-600 hover:bg-gray-700 hover:bg-opacity-25 hover:text-gray-100',
        route.path === resourcePath(resource.key)
          ? 'bg-gray-700 bg-opacity-25 text-blue-100'
          : '',
      ]"
    >
      <span>{{ title(resource.key) }}</span>
    </NuxtLink>
  </nav>
</template>

<script setup lang="ts">
import { pluralize, titleize } from "inflection";

const resources = useResources();
const route = useRoute();

function title(key: string) {
  if (key === "auction_configuration_v2") {
    return "Auction Configurations V2";
  }

  return pluralize(titleize(key));
}

function resourcePath(key: string) {
  if (key === "auction_configuration_v2") {
    return "/v2/auction_configurations";
  }

  return `/${pluralize(key)}`;
}
</script>
