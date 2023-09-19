<template>
  <CreateResourceButton label="New App" :resources-path="resourcesPath" />
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup>
import { ResourceTableFields } from "@/constants";

const columns = [
  ResourceTableFields.PublicUid,
  { field: "platformId", header: "Platform Id" },
  ResourceTableFields.HumanName,
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
  { field: "packageName", header: "Package Name" },
  ResourceTableFields.User,
  { field: "appKey", header: "App Key" },
];
const resourcesPath = "/apps";
</script>
