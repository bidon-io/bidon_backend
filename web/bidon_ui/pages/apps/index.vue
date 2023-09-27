<template>
  <CreateResourceButton label="New App" :resources-path="resourcesPath" />
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup>
import { ResourceTableFields } from "@/constants";

const columns = [
  ResourceTableFields.PublicUid,
  { field: "appKey", header: "App Key" },
  ResourceTableFields.Owner,
  ResourceTableFields.AppName,
  { field: "platformId", header: "Platform" },
  { field: "packageName", header: "Bundle ID / Package Name" },
  {
    field: "",
    header: "Line Items",
    associatedResourcesLink: {
      extractLinkData: ({ id }) => ({
        label: "Line Items",
        path: `/line_items?appId=${id}`,
      }),
    },
  },
  {
    field: "",
    header: "App Demand Profiles",
    associatedResourcesLink: {
      extractLinkData: ({ id }) => ({
        label: "App Demand Profiles",
        path: `/app_demand_profiles?appId=${id}`,
      }),
    },
  },
];
const resourcesPath = "/apps";
</script>
