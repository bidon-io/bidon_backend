<template>
  <CreateResourceButton
    label="New App Demand Profile"
    :resources-path="resourcesPath"
  />
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup>
const columns = [
  { field: "id", header: "Id", sortable: true },
  { field: "appId", header: "App" },
  { field: "demandSourceId", header: "Demand Source" },
  { field: "accountId", header: "Account" },
  { field: "accountType", header: "Account Type" },
];
const resourcesPath = "/app_demand_profiles";
</script>
