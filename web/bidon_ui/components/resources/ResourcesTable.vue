<template>
  <DataTable
    v-model:selection="selectedResources"
    v-model:filters="filters"
    :value="resources"
    data-key="id"
    paginator
    :rows="12"
    :rows-per-page-options="[12, 24, 36, 48]"
    filter-display="row"
    class="whitespace-nowrap"
    @filter="onFilter"
  >
    <template #empty> No data found. </template>
    <Column selection-mode="multiple" header-style="width: 3rem"></Column>
    <Column
      v-for="column in columns"
      :key="column.field"
      :field="column.field"
      :header="column.header"
      :sortable="column.sortable"
      :copyable="column.copyable"
      :filter-field="column.filter?.field"
      :show-filter-menu="false"
    >
      <template
        v-if="column.link || column.associatedResourcesLink || column.copyable"
        #body="{ data, field }"
      >
        <div v-if="column.copyable">
          <button @click="copyField(data[field])">
            <i class="pi pi-copy" style="color: slateblue"></i>
          </button>
          <span>{{ data[field] }}</span>
        </div>
        <ResourceLink v-if="column.link" :link="column.link" :data="data" />
        <AssociatedResourcesLink
          v-if="column.associatedResourcesLink"
          :link="column.associatedResourcesLink"
          :data="data"
        />
      </template>
      <template v-if="column.filter" #filter="{ filterModel, filterCallback }">
        <InputText
          v-if="column.filter.type === 'input'"
          v-model="filterModel.value"
          type="text"
          class="p-column-filter"
          :placeholder="column.filter.placeholder"
          @input="filterCallback()"
        />
        <Dropdown
          v-if="['select', 'select-filter'].includes(column.filter.type)"
          v-model="filterModel.value"
          :options="
            filtersOptions[column.filter.field as keyof typeof filtersOptions]
          "
          option-label="label"
          option-value="value"
          :placeholder="column.filter.placeholder"
          class="p-column-filter"
          :show-clear="true"
          :filter="column.filter.type === 'select-filter'"
          @change="filterCallback()"
        />
      </template>
    </Column>
    <Column
      style="width: 10%; min-width: 8rem"
      body-style="text-align:center; position: sticky; right: 0; background-color: white"
    >
      <template #body="slotProps">
        <div class="flex justify-between">
          <NuxtLink
            :key="slotProps.data.id"
            :to="`${resourcesPath}/${slotProps.data.id}`"
          >
            <i class="pi pi-eye" style="color: slateblue"></i>
          </NuxtLink>
          <NuxtLink
            v-if="slotProps.data._permissions.update"
            :key="slotProps.data.id"
            :to="`${resourcesPath}/${slotProps.data.id}/edit`"
          >
            <i class="pi pi-pencil" style="color: green"></i>
          </NuxtLink>
          <a
            v-if="slotProps.data._permissions.delete"
            :key="slotProps.data.id"
            href="_"
            @:click.prevent="deleteHandle(slotProps.data.id)"
          >
            <i class="pi pi-trash" style="color: red"></i>
          </a>
        </div>
      </template>
    </Column>
  </DataTable>
</template>

<script setup lang="ts">
import type { FilterMatchModeOptions } from "primevue/api";

import useDeleteResource from "@/composables/useDeleteResource";

interface Filter {
  field: string;
  value?: string;
  matchMode: keyof FilterMatchModeOptions;
  type: "input" | "select" | "select-filter";
  placeholder: string;
}

interface SelectFilter extends Filter {
  type: "select";
  extractOptions: (records: object[]) => { label: string; value: string }[];
}

interface Column {
  header: string;
  field: string;
  filter?: Filter;
  sortable?: boolean;
  copyable?: boolean;
  link?: ResourceLink;
  associatedResourcesLink?: AssociatedResourcesLink;
}

const props = defineProps<{
  resourcesPath: string;
  columns: Column[];
}>();

// eslint-disable-next-line  @typescript-eslint/no-explicit-any
const response = await $apiFetch<any[]>(props.resourcesPath);

const resources = ref(response);
const selectedResources = ref([]);

const route = useRoute();
const router = useRouter();

const filters = ref(
  props.columns
    .filter((column) => column.filter)
    .map((column) => column.filter as Filter)
    .reduce(
      (result, filter) => ({
        ...result,
        [filter.field]: {
          matchMode: filter.matchMode,
          value: route.query[filter.field],
        },
      }),
      {},
    ),
);

const getFilterOptions = (filteredResources: object[]) =>
  props.columns
    .filter(
      (column) =>
        column.filter &&
        ["select", "select-filter"].includes(column.filter.type),
    )
    .map((column) => column.filter as SelectFilter)
    .reduce(
      (result, filter) => ({
        ...result,
        [filter.field || ""]: filter.extractOptions(filteredResources) || [],
      }),
      {},
    );

const filtersOptions = ref(getFilterOptions(resources.value));

const onFilter = (event: { filteredValue: object[] }) => {
  filtersOptions.value = getFilterOptions(event.filteredValue);
};

watch(filters, () => {
  const query = Object.entries(filters.value)
    .filter(([, filter]) => (filter as Filter).value)
    .reduce(
      (result, [key, filter]) => ({
        ...result,
        [key]: (filter as Filter).value,
      }),
      {},
    );

  router.push({ query });
});

const deleteHandle = useDeleteResource({
  path: props.resourcesPath,
  hook: (id: number) =>
    (resources.value = resources.value.filter(
      (item: { id: number }) => item.id !== id,
    )),
});

// eslint-disable-next-line  @typescript-eslint/no-explicit-any
const copyField = (field: any) => {
  navigator.clipboard.writeText(field);
};
</script>
