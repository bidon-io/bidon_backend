<template>
  <CreateResourceButton label="New Segment" :resources-path="resourcesPath" />
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup>
const columns = [
  { field: "id", header: "Id", sortable: true },
  { field: "app_id", header: "App" },
  { field: "name", header: "Name" },
  { field: "filters", header: "Filters" },
  { field: "enabled", header: "Enabled" },
];
const resourcesPath = "/segments";
</script>
