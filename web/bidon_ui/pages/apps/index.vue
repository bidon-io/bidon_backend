<template>
  <CreateResourceButton label="New App" :resources-path="resourcesPath" />
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup>
const columns = [
  { field: "id", header: "Id", sortable: true },
  { field: "platform_id", header: "Platform Id" },
  { field: "human_name", header: "Human Name" },
  { field: "package_name", header: "Package Name" },
  { field: "user_id", header: "User" },
  { field: "app_key", header: "App Key" },
];
const resourcesPath = "/apps";
</script>
