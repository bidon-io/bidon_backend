<template>
  <CreateResourceButton label="New User" :resources-path="resourcesPath" />
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup>
import { ResourceTableFields } from "@/constants";

const columns = [ResourceTableFields.Id, { field: "email", header: "Email" }];
const resourcesPath = "/users";
</script>
