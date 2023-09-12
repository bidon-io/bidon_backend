<template>
  <CreateResourceButton label="New User" :resources-path="resourcesPath" />
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup>
import { ResourceTableFields } from "@/constants";

const columns = [
  ResourceTableFields.Id,
  ResourceTableFields.Email,
  { field: "isAdmin", header: "Is Admin" },
];
const resourcesPath = "/users";
</script>
