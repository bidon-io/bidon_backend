<template>
  <CreateResourceButton label="New App" :resources-path="resourcesPath" />
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup>
import { ResourceTableFields } from "@/constants";

const columns = [
  ResourceTableFields.Id,
  { field: "platformId", header: "Platform Id" },
  { field: "humanName", header: "Human Name" },
  { field: "packageName", header: "Package Name" },
  ResourceTableFields.User,
  { field: "appKey", header: "App Key" },
];
const resourcesPath = "/apps";
</script>
