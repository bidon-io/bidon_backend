<template>
  <CreateResourceButton
    label="New Auction Configuration"
    :resources-path="resourcesPath"
  />
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup>
import { ResourceTableFields } from "@/constants";

const columns = [
  ResourceTableFields.Id,
  ResourceTableFields.App,
  { field: "name", header: "Name" },
  ResourceTableFields.AdType,
  { field: "pricefloor", header: "Price Floor" },
  ResourceTableFields.Segment,
];
const resourcesPath = "/auction_configurations";
</script>
