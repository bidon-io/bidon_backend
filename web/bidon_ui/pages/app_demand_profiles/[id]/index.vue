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
      title="App Demand Profile"
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
const resourcesPath = "/app_demand_profiles";

const response = await axios.get(`${resourcesPath}/${id}`);
const resource = response.data;

const fields = [
  ResourceCardFields.PublicUid,
  ResourceCardFields.App,
  ResourceCardFields.DemandSource,
  ResourceCardFields.DemandSourceAccount,
  { label: "Data", key: "data", type: "textarea" },
  { label: "Account Type", key: "accountType" },
];
</script>
