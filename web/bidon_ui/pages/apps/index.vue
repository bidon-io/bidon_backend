<template>
  <CreateResourceButton label="New App" :resources-path="resourcesPath" />
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup>
import { ResourceTableFields } from "@/constants";
import { FilterMatchMode } from "primevue/api";

const columns = [
  ResourceTableFields.AppName,
  {
    field: "platformId",
    header: "Platform",
    filter: {
      field: "platformId",
      type: "select",
      matchMode: FilterMatchMode.EQUALS,
      placeholder: "Platform",
      loadOptions: async () => [
        {
          label: "iOS",
          value: "ios",
        },
        {
          label: "Android",
          value: "android",
        },
      ],
      extractOptions: (records) => [
        ...new Map(
          records.map(({ platformId }) => [
            platformId,
            { label: platformId, value: platformId },
          ]),
        ).values(),
      ],
    },
  },
  {
    field: "packageName",
    header: "Package Name",
    filter: {
      field: "packageName",
      type: "input",
      matchMode: FilterMatchMode.CONTAINS,
      placeholder: "Search by package name",
    },
  },
  ResourceTableFields.Owner,
  {
    field: "",
    header: "Actions",
    associatedResourcesLink: {
      extractLinkData: ({ id }) => ({
        label: "Actions",
        path: `/v2/auction_configurations?appId=${id}`,
        dropdown: [
          {
            label: "Auction Configurations",
            path: `/v2/auction_configurations?appId=${id}`,
          },
          {
            label: "Line Items",
            path: `/line_items?appId=${id}`,
          },
          {
            label: "App Demand Profiles",
            path: `/app_demand_profiles?appId=${id}`,
          },
        ],
      }),
    },
  },
];
const resourcesPath = "/apps";
</script>
