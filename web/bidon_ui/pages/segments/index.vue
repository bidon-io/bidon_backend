<template>
  <CreateResourceButton label="New Segment" :resources-path="resourcesPath" />
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup>
import { ResourceTableFields } from "@/constants";

const columns = [
  ResourceTableFields.Id,
  ResourceTableFields.App,
  { field: "name", header: "Name" },
  { field: "filters", header: "Filters" },
  { field: "enabled", header: "Enabled" },
];
const resourcesPath = "/segments";
</script>
