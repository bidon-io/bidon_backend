<template>
  <CreateResourceButton
    label="New App Demand Profile"
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
  ResourceTableFields.App,
  ResourceTableFields.DemandSource,
  ResourceTableFields.DemandSourceAccount,
];
const resourcesPath = "/app_demand_profiles";
const collectionPath = "/app_demand_profiles_collection";
</script>
