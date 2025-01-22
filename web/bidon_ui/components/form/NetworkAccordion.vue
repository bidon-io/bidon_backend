<template>
  <Accordion v-if="isLoaded" :active-index="0" :multiple="true">
    <AccordionTab v-for="network in networks" :key="network.key">
      <ToggleButton
        v-model="network.enabled"
        class="w-6rem mb-2"
        on-label="On"
        off-label="Off"
        on-icon="pi pi-check"
        off-icon="pi pi-times"
      />
      <template #header>
        <span class="flex align-items-center gap-2 w-full">
          {{ network.label }}
          <Badge
            :value="network.selectedAdUnitIds?.length"
            :severity="network.enabled ? 'success' : 'danger'"
            class="ml-auto mr-2"
          />
        </span>
      </template>
      <Fieldset legend="Ad Units" class="p-fieldset">
        <div class="p-datatable p-datatable-striped">
          <div
            class="p-datatable-header p-grid"
            style="grid-template-columns: 5% 45% 25% 10% 15%"
          >
            <div class="p-text-center">#</div>
            <div>Label</div>
            <div>UID</div>
            <div class="p-text-center">Price</div>
            <div class="p-text-center">Account</div>
          </div>
          <div
            v-for="adUnit in network.adUnits"
            :key="adUnit.id"
            class="p-datatable-row p-grid"
            style="
              grid-template-columns: 5% 45% 25% 10% 15%;
              align-items: center;
            "
          >
            <div class="p-text-center">
              <Checkbox
                v-model="network.selectedAdUnitIds"
                :value="adUnit.id"
              />
            </div>
            <div>{{ adUnit.label }}</div>
            <div>{{ adUnit.uid }}</div>
            <div class="p-text-center">${{ adUnit.pricefloor.toFixed(2) }}</div>
            <div class="p-text-center">{{ adUnit.account }}</div>
          </div>
        </div>
      </Fieldset>
    </AccordionTab>
  </Accordion>
</template>

<script lang="ts" setup>
import axios from "@/services/ApiService";

type Network = {
  label: string;
  key: string;
  enabled: boolean;
  adUnits: AdUnit[];
  isBidding: boolean;
  selectedAdUnitIds: number[];
};

type AdUnit = {
  id: number;
  uid: string;
  label: string;
  networkKey: string;
  isBidding: boolean;
  pricefloor: number;
  account: string;
};

const props = defineProps({
  appId: {
    type: Number as PropType<number>,
    default: null,
  },
  adType: {
    type: String as PropType<string>,
    default: "",
  },
  isBidding: {
    type: Boolean as PropType<boolean>,
    default: false,
  },
  networkKeys: {
    type: Array as PropType<string[]>,
    default: () => [],
  },
  adUnitIds: {
    type: Array as PropType<number[]>,
    default: () => [],
  },
});
const emit = defineEmits(["update:networkKeys", "update:adUnitIds"]);

// TODO: Fetch from API instead of hardcoded
const networks = ref<Network[]>([
  // CPM networks
  {
    label: "Admob",
    key: "admob",
    isBidding: false,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
  {
    label: "Applovin",
    key: "applovin",
    isBidding: false,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
  {
    label: "BidMachine",
    key: "bidmachine",
    isBidding: false,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
  {
    label: "Bigoads",
    key: "bigoads",
    isBidding: false,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
  {
    label: "Chartboost",
    key: "chartboost",
    isBidding: false,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
  {
    label: "DtExchange",
    key: "dtexchange",
    isBidding: false,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
  {
    label: "Google Ad Manager",
    key: "gam",
    isBidding: false,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
  {
    label: "Mintegral",
    key: "mintegral",
    isBidding: false,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
  {
    label: "UnityAds",
    key: "unityads",
    isBidding: false,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
  {
    label: "IronSource",
    key: "ironsource",
    isBidding: false,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
  {
    label: "VK Ads",
    key: "vkads",
    isBidding: false,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
  {
    label: "Vungle",
    key: "vungle",
    isBidding: false,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
  {
    label: "Yandex",
    key: "yandex",
    isBidding: false,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
  // Bidding networks
  {
    label: "Amazon",
    key: "amazon",
    isBidding: true,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
  {
    label: "BidMachine",
    key: "bidmachine",
    isBidding: true,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
  {
    label: "Bigoads",
    key: "bigoads",
    isBidding: true,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
  {
    label: "Meta",
    key: "meta",
    isBidding: true,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
  {
    label: "Mintegral",
    key: "mintegral",
    isBidding: true,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
  {
    label: "MobileFuse",
    key: "mobilefuse",
    isBidding: true,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
  {
    label: "Vungle",
    key: "vungle",
    isBidding: true,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
  {
    label: "VK Ads",
    key: "vkads",
    isBidding: true,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
]);
const isLoaded = ref(false);

const fetchAdUnits = async () => {
  if (!props.appId || !props.adType) return [];
  const url = `/line_items?app_id=${props.appId}&ad_type=${props.adType}&is_bidding=${props.isBidding}`;
  const { data, error } = await useAsyncData(url, async () => {
    const result = (await axios.get(url)).data;
    return result.map(
      (
        adUnit: any, // eslint-disable-line @typescript-eslint/no-explicit-any
      ) => ({
        id: adUnit.id,
        label: adUnit.humanName,
        networkKey: adUnit.accountType.split("::")[1].toLowerCase(),
        uid: adUnit.publicUid,
        pricefloor: parseFloat(adUnit.bidFloor),
        account: `${adUnit.accountType.split("::")[1].toLowerCase()} (${
          adUnit.accountId
        })`,
        isBidding: adUnit.isBidding,
      }),
    ) as AdUnit[];
  });
  if (error.value !== null) {
    console.error("Failed to fetch ad units:", error);
    return [];
  }
  return data.value || [];
};

watch(
  () => [props.appId, props.adType, props.isBidding],
  async () => {
    const adUnits = await fetchAdUnits();
    const bidTypeNetworks = networks.value.filter(
      (network) => network.isBidding === props.isBidding,
    );
    const updatedNetworks = bidTypeNetworks.map((network) => {
      const networkAdUnits = adUnits
        .filter((adUnit) => adUnit.networkKey === network.key)
        .sort((adUnit) => adUnit.pricefloor);
      return {
        ...network,
        enabled: props.networkKeys.includes(network.key),
        adUnits: networkAdUnits,
        selectedAdUnitIds: props.adUnitIds.filter((id) =>
          networkAdUnits.some((unit) => unit.id === id),
        ),
      };
    });
    isLoaded.value = true;
    networks.value = updatedNetworks;
  },
  { immediate: true },
);

watchEffect(() => {
  if (!isLoaded.value) return;

  const enabledNetworks = networks.value.filter((network) => network.enabled);
  const selectedNetworkKeys = enabledNetworks.map((network) => network.key);
  const selectedAdUnitIds = enabledNetworks
    .map((network) => network.selectedAdUnitIds)
    .flat();

  emit("update:networkKeys", selectedNetworkKeys);
  emit("update:adUnitIds", selectedAdUnitIds);
});
</script>
<style scoped>
.p-datatable-header,
.p-datatable-row {
  display: grid;
}
.p-datatable-header {
  font-weight: bold;
  background-color: var(--surface-b);
  border-bottom: 1px solid var(--surface-d);
}
.p-datatable-row {
  border-bottom: 1px solid var(--surface-d);
}
</style>
