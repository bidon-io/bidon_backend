<template>
  <CreateResourceButton label="New Segment" :resources-path="resourcesPath" />
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup>
import { ResourceTableFields } from "@/constants";

const columns = [
  ResourceTableFields.App,
  { field: "name", header: "Name" },
  { field: "enabled", header: "Enabled" },
  { field: "priority", header: "Priority" },
];
const resourcesPath = "/segments";
</script>
