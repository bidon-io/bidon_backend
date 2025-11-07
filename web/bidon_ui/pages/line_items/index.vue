<template>
  <CreateResourceButton label="New Line Item" :resources-path="resourcesPath" />
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
  ResourceTableFields.HumanName,
  ResourceTableFields.App,
  ResourceTableFields.BidType,
  ResourceTableFields.AdTypeWithFormat,
  ResourceTableFields.DemandSourceType,
  ResourceTableFields.DemandSourceAccount,
];
const resourcesPath = "/line_items";
const collectionPath = "/line_items_collection";
</script>
