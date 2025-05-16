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
  ResourceTableFields.DemandSource,
  ResourceTableFields.Owner,
  ResourceTableFields.Label,
];
const resourcesPath = "/demand_source_accounts";
</script>
