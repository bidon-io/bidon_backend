<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton :path="resourcesPath" />
      <DestroyButton :id="id" :path="resourcesPath" />
      <EditButton :id="id" :path="resourcesPath" />
    </NavigationContainer>
    <ResourceCard title="App" :fields="fields" :resource="resource" />
  </PageContainer>
</template>

<script setup>
import axios from "@/services/ApiService.js";
import { ResourceCardFields } from "@/constants";

const route = useRoute();
const id = route.params.id;
const resourcesPath = "/apps";

const response = await axios.get(`${resourcesPath}/${id}`);
const resource = response.data;

const fields = [
  ResourceCardFields.PublicUid,
  { label: "App Key", key: "appKey" },
  ResourceCardFields.Owner,
  { label: "App Name", key: "humanName" },
  { label: "Platform", key: "platformId" },
  { label: "Bundle ID / Package Name", key: "packageName" },
  {
    key: "lineItems",
    label: "Line Items",
    type: "associatedResourcesLink",
    associatedResourcesLink: {
      extractLinkData: ({ id }) => ({
        label: "Line Items",
        path: `/line_items?appId=${id}`,
      }),
    },
  },
  {
    key: "appDemandProfiles",
    label: "App Demand Profiles",
    type: "associatedResourcesLink",
    associatedResourcesLink: {
      extractLinkData: ({ id }) => ({
        label: "App Demand Profiles",
        path: `/app_demand_profiles?appId=${id}`,
      }),
    },
  },
];
</script>
