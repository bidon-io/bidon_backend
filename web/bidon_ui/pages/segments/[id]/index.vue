<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton :path="resourcesPath" />
      <DestroyButton :id="id" :path="resourcesPath" />
      <EditButton :id="id" :path="resourcesPath" />
    </NavigationContainer>
    <ResourceCard title="Segment" :fields="fields" :resource="resource" />
  </PageContainer>
</template>

<script setup>
import axios from "@/services/ApiService.js";
import { ResourceCardFields } from "@/constants";

const route = useRoute();
const id = route.params.id;
const resourcesPath = "/segments";

const response = await axios.get(`${resourcesPath}/${id}`);
const resource = response.data;

const fields = [
  ResourceCardFields.Id,
  { label: "Public UID", key: "publicUid" },
  { label: "Name", key: "name" },
  { label: "Description", key: "description" },
  { label: "Filters", key: "filters" },
  { label: "Enabled", key: "enabled" },
  { label: "Priority", key: "priority" },
  ResourceCardFields.App,
];
</script>
