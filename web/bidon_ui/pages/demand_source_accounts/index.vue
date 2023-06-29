<template>
  <CreateResourceButton label="New Demand Source Account" :resources-path="resourcesPath" />
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup>
const columns = [
  { field: "id", header: "Id", sortable: true },
  { field: "userId", header: "User" },
  { field: "type", header: "Type" },
  { field: "demandSourceId", header: "Demand Source" },
  { field: "isBidding", header: "Bidding" },
];
const resourcesPath = "/demand_source_accounts";
</script>
