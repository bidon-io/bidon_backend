<template>
  <Toast />
  <ConfirmDialog />
  <NavigationContainer>
    <NuxtLink to="/auction_configurations/new">
      <Button label="New Auction Configuration" icon="pi pi-plus" class="p-button-success" />
    </NuxtLink>
  </NavigationContainer>
  <DataTable
    v-model:selection="selectedConfigurations"
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
          <a :key="slotProps.data.id" href="_" @:click.prevent="deleteHandle(slotProps.data.id)">
            <i class="pi pi-trash" style="color: red"></i>
          </a>
        </div>
      </template>
    </Column>
  </DataTable>
</template>

<script setup>
import { ref } from "vue";
import axios from "@/services/ApiService.js";

const path = "/auction_configurations";
const response = await axios.get(path);
const configurations = ref(response.data);
const selectedConfigurations = ref([]);

const deleteHandle = useDeleteResource(
  path,
  (id) => (configurations.value = configurations.value.filter((item) => item.id !== id))
);
</script>
