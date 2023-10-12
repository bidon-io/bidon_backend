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
      title="Demand Source Account"
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
const resourcesPath = "/demand_source_accounts";

const response = await axios.get(`${resourcesPath}/${id}`);
const resource = response.data;

const fields = [
  ResourceCardFields.PublicUid,
  ResourceCardFields.DemandSource,
  ResourceCardFields.Owner,
  { label: "Label", key: "label" },
  { label: "Extra", key: "extra", type: "textarea" },
];
</script>
