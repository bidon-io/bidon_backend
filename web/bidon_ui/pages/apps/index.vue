<template>
  <CreateResourceButton label="New App" :resources-path="resourcesPath" />
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup>
const columns = [
  { field: "id", header: "Id", sortable: true },
  { field: "platformId", header: "Platform Id" },
  { field: "humanName", header: "Human Name" },
  { field: "packageName", header: "Package Name" },
  { field: "userId", header: "User" },
  { field: "appKey", header: "App Key" },
];
const resourcesPath = "/apps";
</script>
