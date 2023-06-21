<template>
  <DataTable
    v-model:selection="selectedResources"
    :value="resources"
    data-key="id"
    paginator
    :rows="12"
    :rows-per-page-options="[12, 24, 36, 48]"
    table-style="min-width: 50rem"
  >
    <Column selection-mode="multiple" header-style="width: 3rem"></Column>
    <Column
      v-for="column in columns"
      :key="column.field"
      :field="column.field"
      :header="column.header"
      :sortable="column.sortable"
    />
    <Column style="width: 10%; min-width: 8rem" body-style="text-align:center">
      <template #body="slotProps">
        <div class="flex justify-between">
          <NuxtLink :key="slotProps.data.id" :to="`${resourcesPath}/${slotProps.data.id}`">
            <i class="pi pi-eye" style="color: slateblue"></i>
          </NuxtLink>
          <NuxtLink :key="slotProps.data.id" :to="`${resourcesPath}/${slotProps.data.id}/edit`">
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

<script setup lang="ts">
import { ref } from "vue";
import axios from "@/services/ApiService.js";
import useDeleteResource from "@/composables/useDeleteResource";

interface Column {
  header: string;
  field: string;
  sortable?: boolean;
}

const props = defineProps<{
  resourcesPath: string;
  columns: Column[];
}>();

const response = await axios.get(props.resourcesPath);
const resources = ref(response.data);
const selectedResources = ref([]);

const deleteHandle = useDeleteResource({
  path: props.resourcesPath,
  hook: (id: number) => (resources.value = resources.value.filter((item: { id: number }) => item.id !== id)),
});
</script>
