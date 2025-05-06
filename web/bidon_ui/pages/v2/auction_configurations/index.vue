<template>
  <CreateResourceButton
    label="New Auction Configuration"
    :resources-path="resourcesPath"
  />
  <LazyResourcesTable
    :columns="columns"
    :resources-path="resourcesPath"
    :collection-path="collectionPath"
  />
</template>

<script setup>
import { ResourceTableFields } from "@/constants";
import LazyResourcesTable from "~/components/resources/LazyResourcesTable.vue";

const columns = [
  ResourceTableFields.PublicUid,
  ResourceTableFields.App,
  ResourceTableFields.Name,
  ResourceTableFields.AuctionKey,
  ResourceTableFields.AdType,
  { field: "pricefloor", header: "Price Floor" },
  ResourceTableFields.IsDefault,
  ResourceTableFields.Segment,
];
const resourcesPath = "/v2/auction_configurations";
const collectionPath = "/v2/auction_configurations_collection";
</script>
