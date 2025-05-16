<template>
  <CreateResourceButton
    label="New Demand Source"
    :resources-path="resourcesPath"
  />
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup>
const columns = [
  { field: "humanName", header: "Name" },
  { field: "apiKey", header: "Api Key" },
];
const resourcesPath = "/demand_sources";
</script>
