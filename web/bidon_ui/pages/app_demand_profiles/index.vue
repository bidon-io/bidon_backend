<template>
  <CreateResourceButton
    label="New App Demand Profile"
    :resources-path="resourcesPath"
  />
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup>
import { ResourceTableFields } from "@/constants";

const columns = [
  ResourceTableFields.Id,
  ResourceTableFields.App,
  ResourceTableFields.DemandSource,
  ResourceTableFields.DemandSourceAccount,
  { field: "accountType", header: "Account Type" },
];
const resourcesPath = "/app_demand_profiles";
</script>
