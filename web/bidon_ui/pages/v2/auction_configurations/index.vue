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
  ResourceTableFields.PublicUid,
  ResourceTableFields.App,
  { field: "name", header: "Name" },
  { field: "auctionKey", header: "Auction Key", copyable: true },
  ResourceTableFields.AdType,
  { field: "pricefloor", header: "Price Floor" },
  ResourceTableFields.IsDefault,
  ResourceTableFields.Segment,
];
const resourcesPath = "/v2/auction_configurations";
</script>
