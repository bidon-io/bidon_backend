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
    <ResourceCard
      title="Auction Configuration"
      :fields="fields"
      :resource="resource"
    />
  </PageContainer>
</template>

<script setup>
import axios from "@/services/ApiService.js";
import { ResourceCardFields } from "@/constants";

const route = useRoute();
const id = route.params.id;
const resourcesPath = "/v2/auction_configurations";

const response = await axios.get(`${resourcesPath}/${id}`);
const resource = response.data;

const fields = [
  ResourceCardFields.PublicUid,
  ResourceCardFields.App,
  { label: "Name", key: "name" },
  { label: "Auction Key", key: "auctionKey", copyable: true },
  { label: "Ad type", key: "adType" },
  { label: "Price floor", key: "pricefloor" },
  { label: "Demands", key: "demands" },
  { label: "Bidding", key: "bidding" },
  { label: "Ad Unit IDs", key: "adUnitIds" },
  { label: "Timeout", key: "timeout" },
  { label: "Default", key: "isDefault" },
  { label: "External Win Notification", key: "externalWinNotifications" },
  ResourceCardFields.Segment,
];
</script>
