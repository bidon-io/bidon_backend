<template>
  <Message v-if="error" severity="error" closable>{{ error?.message }}</Message>
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
    :total-records
    :loading="status === 'pending'"
    lazy
    @filter="onFilter"
    @page="onPage"
    @update:rows="onLimit"
  >
    <template v-if="status == 'success'" #empty> No data found. </template>
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
      :body-class="column.bodyClass"
      :body-style="column.bodyStyle"
      :header-style="column.headerStyle"
    >
      <template
        v-if="
          column.link ||
          column.associatedResourcesLink ||
          column.copyable ||
          column.customBody
        "
        #body="{ data, field }"
      >
        <div v-if="column.copyable">
          <button @click="copyField(data[field])">
            <i class="pi pi-copy" style="color: slateblue"></i>
          </button>
          <span>{{ data[field] }}</span>
        </div>
        <span v-else-if="column.customBody">{{ column.customBody(data) }}</span>
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
          v-model.lazy="filterModel.value"
          type="text"
          class="p-column-filter"
          :placeholder="column.filter.placeholder"
          @input="debouncedFilter(filterCallback)"
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
          @focus="loadFilterOptions(column.filter.field)"
        />
      </template>
    </Column>
    <Column
      style="width: 12%; min-width: 10rem"
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
          <NuxtLink
            v-if="resourcePermissions?.create"
            :key="slotProps.data.id"
            :to="`${resourcesPath}/new?clone=${slotProps.data.id}`"
          >
            <i class="pi pi-copy" style="color: orange"></i>
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
import { decamelize } from "humps";
import { camelize, singularize } from "inflection";
import { debounce } from "~/utils/debounce";
import { buildQueryParams as buildFilterQueryParams } from "~/utils/filterUtils";

// Define a Resource interface for the table rows
interface Resource {
  id: string | number;
  _permissions: {
    update: boolean;
    delete: boolean;
  };
  [key: string]: unknown;
}

interface Filter {
  field: string;
  value?: string;
  matchMode: keyof FilterMatchModeOptions;
  type: "input" | "select" | "select-filter";
  placeholder: string;
}

interface SelectFilter extends Filter {
  type: "select";
  loadOptions?: () => Promise<{ label: string; value: string }[]>;
}

interface Column {
  header: string;
  field: string;
  filter?: Filter;
  sortable?: boolean;
  copyable?: boolean;
  link?: ResourceLink;
  associatedResourcesLink?: AssociatedResourcesLink;
  customBody?: (rowData: Resource) => string;

  bodyClass?: string;
  bodyStyle?: string;
  headerStyle?: string;
}

interface PageEvent {
  page: number;
}

const props = defineProps<{
  resourcesPath: string;
  collectionPath: string;
  columns: Column[];
}>();

const selectedResources = ref([]);

const route = useRoute();
const router = useRouter();

// Get resource permissions for clone functionality
const resourceKey = computed(() => {
  if (props.resourcesPath === "/v2/auction_configurations") {
    return "auctionConfigurationV2";
  }
  return camelize(singularize(props.resourcesPath.replace("/", "")), true);
});

const resourcesStore = useResources();
const resourcePermissions = computed(
  () => resourcesStore.state[resourceKey.value]?.permissions ?? {},
);

// read query params
const page = ref(+(route.query.page ?? 1));
const limit = ref(+(route.query.limit ?? 12));

const buildQueryParams = () => {
  return buildFilterQueryParams(filters.value, page.value, limit.value);
};

type FilterValue = string | undefined | { adType: string; format: string };

const debouncedFilter = debounce(
  (filterCallback: () => void) => filterCallback(),
  500,
);

// Initialize filters from URL query parameters
const filters = ref<
  Record<
    string,
    { matchMode: keyof FilterMatchModeOptions; value: FilterValue }
  >
>(
  props.columns
    .filter((column) => column.filter)
    .map((column) => column.filter as Filter)
    .reduce((result, filter) => {
      let value;

      // Special handling for AdTypeWithFormat
      if (filter.field === "adTypeWithFormat") {
        const adType = route.query["ad_type"] as string;
        const format = route.query["format"] as string;

        if (adType) {
          value = { adType, format: format || "" };
        }
      } else {
        // Regular handling for other fields
        value =
          (route.query[decamelize(filter.field)] as string | undefined) ||
          (route.query[filter.field] as string | undefined);
      }

      return {
        ...result,
        [filter.field]: {
          matchMode: filter.matchMode,
          value,
        },
      };
    }, {}),
);

const filtersOptions = ref<Record<string, { label: string; value: string }[]>>(
  {},
);

const loadFilterOptions = async (field: string) => {
  if (filtersOptions.value[field]) {
    return filtersOptions.value[field];
  }

  const column = props.columns.find((column) => column.filter?.field === field);
  const filter = column?.filter as SelectFilter;
  const options = await filter?.loadOptions?.();

  if (options) {
    filtersOptions.value = {
      ...filtersOptions.value,
      [field]: options,
    };

    const currentFilterValue = filters.value[field]?.value;
    const matchedOption = options.find(
      (opt) => opt.value == currentFilterValue,
    );

    if (matchedOption) {
      filters.value = {
        ...filters.value,
        [field]: {
          ...filters.value[field],
          value: matchedOption.value,
        },
      };
    }
  }

  return options;
};

// Preload filter options for filters with values in URL
const preloadFilterOptions = async () => {
  const filtersWithValues = Object.entries(filters.value)
    .filter(([, filter]) => filter?.value)
    .map(([field]) => field);

  for (const field of filtersWithValues) {
    await loadFilterOptions(field);
  }
};

// Fetch resources
const {
  data: collection,
  status,
  error,
  execute: fetchData,
} = useAsyncData(
  "fetch-resources-collection",
  async () => {
    await preloadFilterOptions();
    const params = buildQueryParams();
    return await $apiFetch(props.collectionPath, { params });
  },
  {
    default: () => [],
    immediate: true,
  },
);

const resources = computed(() => collection.value?.items ?? []);
const totalRecords = computed(() => collection.value?.meta?.totalCount ?? 0);

const onFilter = async () => {
  await fetchData();
};

const onPage = async (event: PageEvent) => {
  page.value = event.page + 1;
  await fetchData();
};

const onLimit = async (value: number) => {
  limit.value = value;
  await fetchData();
};

// Update URL when filters, page, or limit change
watch(
  [page, limit, filters],
  () => {
    const query = buildQueryParams();
    router.push({ query: query as Record<string, string | number> });
  },
  { deep: true },
);

const deleteHandle = useDeleteResource({
  path: props.resourcesPath,
  hook: () => {
    fetchData();
  },
});

const copyField = (field: string) => {
  navigator.clipboard.writeText(field);
};
</script>
