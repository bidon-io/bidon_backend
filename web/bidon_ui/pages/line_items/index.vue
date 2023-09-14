<template>
  <CreateResourceButton label="New Line Item" :resources-path="resourcesPath" />
  <ResourcesTable :columns="columns" :resources-path="resourcesPath" />
</template>

<script setup>
import { ResourceTableFields } from "@/constants";

const columns = [
  ResourceTableFields.Id,
  { field: "publicUid", header: "Public UID" },
  ResourceTableFields.HumanName,
  ResourceTableFields.App,
  ResourceTableFields.BidFloor,
  ResourceTableFields.AdType,
  { field: "format", header: "Format" },
  ResourceTableFields.DemandSourceAccount,
  ResourceTableFields.AccountType,
  { field: "code", header: "Code" },
];
const resourcesPath = "/line_items";
</script>
