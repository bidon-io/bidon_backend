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
      <span>{{ pluralize(titleize(resource.key)) }}</span>
    </NuxtLink>
  </nav>
</template>

<script setup lang="ts">
import { pluralize, titleize } from "inflection";

const resources = useResources();
const route = useRoute();

function resourcePath(key: string) {
  return `/${pluralize(key)}`;
}
</script>
