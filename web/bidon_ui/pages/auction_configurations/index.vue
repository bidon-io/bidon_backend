<template>
  <CreateResourceButton label="New Auction Configuration" :resources-path="resourcesPath" />
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup>
const columns = [
  { field: "id", header: "Id", sortable: true },
  { field: "app_id", header: "App" },
  { field: "name", header: "Name" },
  { field: "ad_type", header: "Ad Type" },
  { field: "pricefloor", header: "Price Floor" },
  { field: "segment_id", header: "Segment" },
];
const resourcesPath = "/auction_configurations";
</script>
