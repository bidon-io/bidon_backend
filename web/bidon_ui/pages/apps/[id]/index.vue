<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton :path="resourcesPath" />
      <DestroyButton
        v-if="resource._permissions.delete"
        :id="id"
        :path="resourcesPath"
      />
      <EditButton
        v-if="resource._permissions.update"
        :id="id"
        :path="resourcesPath"
      />
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
  { ...ResourceCardFields.PublicUid, copyable: true },
  { label: "App Key", key: "appKey", copyable: true },
  ResourceCardFields.Owner,
  { label: "App Name", key: "humanName" },
  { label: "Platform", key: "platformId" },
  { label: "Package Name", key: "packageName" },
  { label: "Store ID", key: "storeId" },
  { label: "Store URL", key: "storeUrl" },
  { label: "Categories", key: "categories" },
  { label: "Blocked Advertiser Domains", key: "badv" },
  { label: "Blocked Categories", key: "bcat" },
  { label: "Blocked Apps", key: "bapp" },
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
  {
    key: "auctionConfigurations",
    label: "Auction Configurations",
    type: "associatedResourcesLink",
    associatedResourcesLink: {
      extractLinkData: ({ id }) => ({
        label: "Auction Configurations",
        path: `/v2/auction_configurations?appId=${id}`,
      }),
    },
  },
];
</script>
