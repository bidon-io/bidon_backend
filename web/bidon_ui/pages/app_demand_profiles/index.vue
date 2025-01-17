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
  ResourceTableFields.PublicUid,
  ResourceTableFields.App,
  ResourceTableFields.DemandSource,
  ResourceTableFields.DemandSourceAccount,
  { field: "accountType", header: "Account Type" },
];
const resourcesPath = "/app_demand_profiles";
const collectionPath = "/app_demand_profiles_collection";
</script>
