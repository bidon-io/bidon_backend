<template>
  <CreateResourceButton
    label="New Demand Source Account"
    :resources-path="resourcesPath"
  />
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup>
import { ResourceTableFields } from "@/constants";

const columns = [
  ResourceTableFields.Id,
  ResourceTableFields.User,
  ResourceTableFields.Label,
  { field: "type", header: "Type" },
  ResourceTableFields.DemandSource,
  { field: "isBidding", header: "Bidding" },
];
const resourcesPath = "/demand_source_accounts";
</script>
