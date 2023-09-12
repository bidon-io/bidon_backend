<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton :path="resourcesPath" />
      <DestroyButton :id="id" :path="resourcesPath" />
      <EditButton :id="id" :path="resourcesPath" />
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

const fields = [
  ResourceCardFields.Id,
  ResourceCardFields.HumanName,
  ResourceCardFields.App,
  ResourceCardFields.BidFloor,
  ResourceCardFields.AdType,
  { key: "format", label: "Format" },
  ResourceCardFields.DemandSourceAccount,
  ResourceCardFields.AccountType,
  { key: "code", label: "Code" },
  { key: "extra", label: "Extra", type: "textarea" },
];
</script>
