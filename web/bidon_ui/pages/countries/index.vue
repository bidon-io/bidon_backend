<template>
  <CreateResourceButton label="New Country" :resources-path="resourcesPath" />
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup>
import { ResourceTableFields } from "@/constants";

const columns = [
  ResourceTableFields.Id,
  { field: "humanName", header: "Human Name" },
  { field: "alpha2Code", header: "Alpha 2 Code" },
  { field: "alpha3Code", header: "Alpha 3 Code" },
];
const resourcesPath = "/countries";
</script>
