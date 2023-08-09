<template>
  <CreateResourceButton
    label="New Demand Source"
    :resources-path="resourcesPath"
  />
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup>
import { ResourceTableFields } from "@/constants";

const columns = [
  ResourceTableFields.Id,
  { field: "humanName", header: "Human Name" },
  { field: "apiKey", header: "Api Key" },
];
const resourcesPath = "/demand_sources";
</script>
