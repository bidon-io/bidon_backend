<template>
  <nav class="mt-6">
    <NuxtLink
      v-for="resource in filteredResources"
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

    <div class="mt-4 px-6">
      <span class="text-gray-400 uppercase text-sm">Settings</span>
      <NuxtLink
        to="/settings/security"
        :class="[
          'flex items-center mt-4 px-6 py-2 text-gray-600 hover:bg-gray-700 hover:bg-opacity-25 hover:text-gray-100',
          route.path === '/settings/security'
            ? 'bg-gray-700 bg-opacity-25 text-blue-100'
            : '',
        ]"
      >
        <span>Security</span>
      </NuxtLink>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { pluralize, titleize } from "inflection";

interface User {
  isAdmin?: boolean;
}

const resources = useResources();
const route = useRoute();
const authStore = useAuthStore();
const currentUser = computed<User | null>(() => authStore.currentUser);

// Filter resources based on permissions
const filteredResources = computed(() => {
  if (!resources.state) return [];

  return Object.values(resources.state).filter((resource) => {
    // Always hide country resource
    if (resource.key === "country") {
      return false;
    }

    // Only show demand_source to admin users
    if (
      resource.key === "demand_source" &&
      (!currentUser.value || currentUser.value?.isAdmin !== true)
    ) {
      return false;
    }

    return true;
  });
});

function title(key: string) {
  if (key === "auction_configuration_v2") {
    return "Auction Configurations";
  }

  if (key === "api_key") {
    return "API Credentials";
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
