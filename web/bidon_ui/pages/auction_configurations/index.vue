<template>
  <DataTable
    v-model:selection="selectedConfiguration"
    :value="configurations"
    data-key="id"
    paginator
    :rows="12"
    :rows-per-page-options="[12, 24, 36, 48]"
    table-style="min-width: 50rem"
  >
    <Column selection-mode="multiple" header-style="width: 3rem"></Column>
    <Column field="id" header="Id" sortable></Column>
    <Column field="app" header="App"></Column>
    <Column field="name" header="Name"></Column>
    <Column field="adType" header="AdType"></Column>
    <Column field="priceFloor" header="PriceFloor"></Column>
    <Column style="width: 10%; min-width: 8rem" body-style="text-align:center">
      <template #body="slotProps">
        <div class="flex justify-between">
          <NuxtLink :key="slotProps.data.id" :to="`/auction_configurations/${slotProps.data.id}`">
            <i class="pi pi-eye" style="color: slateblue"></i>
          </NuxtLink>
          <NuxtLink :key="slotProps.data.id" :to="`/auction_configurations/${slotProps.data.id}/edit`">
            <i class="pi pi-pencil" style="color: green"></i>
          </NuxtLink>
          <NuxtLink :key="slotProps.data.id" :to="`/auction_configurations/${slotProps.data.id}`">
            <i class="pi pi-trash" style="color: red"></i>
          </NuxtLink>
        </div>
      </template>
    </Column>
  </DataTable>
</template>

<script setup>
import { ref } from "vue";

const sampleConfigurations = Array.from({ length: 100 }, (_, index) => ({
  id: index,
  app: `App${index}`,
  name: `Configuration${index}`,
  adType: "banner",
  priceFloor: 0.01 * index,
}));
const configurations = ref(sampleConfigurations);
const selectedConfiguration = ref();
</script>
