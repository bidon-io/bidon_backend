<template>
  <CreateResourceButton
    label="New Auction Configuration"
    :resources-path="resourcesPath"
  />
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup>
const columns = [
  { field: "id", header: "Id", sortable: true },
  { field: "appId", header: "App" },
  { field: "name", header: "Name" },
  { field: "adType", header: "Ad Type" },
  { field: "pricefloor", header: "Price Floor" },
  { field: "segmentId", header: "Segment" },
];
const resourcesPath = "/auction_configurations";
</script>
