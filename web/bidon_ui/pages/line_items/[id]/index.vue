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
    <ResourceCard title="Line Item" :fields="fields" :resource="resource" />
  </PageContainer>
</template>

<script setup>
import axios from "@/services/ApiService.js";
import { ResourceCardFields } from "@/constants";

const route = useRoute();
const id = route.params.id;
const resourcesPath = "/line_items";

const response = await axios.get(`${resourcesPath}/${id}`);
const resource = response.data;

const jsonFields = jsonToFields(resource.extra, "extra", "static", true);
const fields = [
  ResourceCardFields.PublicUid,
  ResourceCardFields.HumanName,
  ResourceCardFields.App,
  ResourceCardFields.BidFloor,
  ResourceCardFields.AdType,
  ...(resource.format ? [{ key: "format", label: "Format" }] : []),
  ResourceCardFields.DemandSourceAccount,
  { key: "isBidding", label: "Bidding" },
  ...jsonFields,
];
</script>
