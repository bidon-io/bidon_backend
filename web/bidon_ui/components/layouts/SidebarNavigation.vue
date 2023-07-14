<template>
  <nav class="mt-10">
    <NuxtLink
      v-for="item in menuItems"
      :key="item.name"
      :to="item.link"
      class="flex items-center mt-4 py-2 px-6 text-gray-600 hover:bg-gray-700 hover:bg-opacity-25 hover:text-gray-100"
    >
      <span class="mx-3">{{ item.name }}</span>
    </NuxtLink>
  </nav>
</template>

<script setup>
const menuItems = [
  { name: "App Demand Profile", link: "/app_demand_profiles" },
  { name: "Apps", link: "/apps" },
  { name: "Auction Configurations", link: "/auction_configurations" },
  { name: "Demand Source Accounts", link: "/demand_source_accounts" },
  { name: "Line Items", link: "/line_items" },
  { name: "Segments", link: "/segments" },
];
</script>
