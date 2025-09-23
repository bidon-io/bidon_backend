<template>
  <div v-if="adUnitsStatus === 'pending'" class="flex items-center p-4">
    <i class="pi pi-spin pi-spinner mr-2"></i>
    <span>Loading networks...</span>
  </div>
  <Accordion v-else-if="isLoaded" :active-index="0" :multiple="true">
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
        <!-- For Waterfall Networks -->
        <template v-if="!network.isBidding">
          <div class="flex flex-col gap-2">
            <div
              v-for="(adUnit, index) in network.adUnits"
              :key="adUnit.id"
              class="flex items-center p-2"
              :class="{
                'border-b border-gray-200':
                  index !== network.adUnits.length - 1,
              }"
            >
              <Checkbox
                v-model="network.selectedAdUnitIds"
                :value="adUnit.id"
                class="mr-3"
              />
              <div class="flex-grow">
                {{ adUnit.label }} -
                <span class="font-medium"
                  >${{ adUnit.pricefloor.toFixed(2) }}</span
                >
              </div>
            </div>
          </div>
        </template>

        <!-- For Bidding Networks -->
        <template v-else>
          <div class="flex flex-col gap-2">
            <div
              v-for="(adUnit, index) in network.adUnits"
              :key="adUnit.id"
              class="flex items-center p-2"
              :class="{
                'border-b border-gray-200':
                  index !== network.adUnits.length - 1,
              }"
            >
              <Checkbox
                v-model="network.selectedAdUnitIds"
                :value="adUnit.id"
                class="mr-3"
              />
              <div class="flex-grow">
                {{ adUnit.label }}
              </div>
            </div>
          </div>
        </template>
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
  initialAdUnitIds: {
    type: Array as PropType<number[]>,
    default: () => [],
  },
});
const emit = defineEmits([
  "update:networkKeys",
  "update:adUnitIds",
  "network-enabled",
  "network-disabled",
]);

// TODO: Fetch from API instead of hardcoded
const networks = ref<Network[]>([
  // Waterfall networks
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
    label: "InMobi",
    key: "inmobi",
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
    label: "Moloco",
    key: "moloco",
    isBidding: true,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
  {
    label: "TaurusX",
    key: "taurusx",
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

// Fetch ad units using useAsyncData with proper caching
const {
  data: adUnitsData,
  status: adUnitsStatus,
  error: adUnitsError,
} = useAsyncData(
  `line-items-${props.appId}-${props.adType}-${props.isBidding}`,
  async () => {
    if (!props.appId || !props.adType) return [];

    const url = `/line_items?app_id=${props.appId}&ad_type=${props.adType}&is_bidding=${props.isBidding}`;
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
  },
  {
    default: () => [],
    server: false,
  },
);

// Watch for errors in ad units fetching
watch(adUnitsError, (error) => {
  if (error) {
    console.error("Failed to fetch ad units:", error);
  }
});

const updateNetworks = () => {
  if (!adUnitsData.value) return;

  const adUnits = adUnitsData.value;
  const bidTypeNetworks = networks.value.filter(
    (network) => network.isBidding === props.isBidding,
  );
  const updatedNetworks = bidTypeNetworks.map((network) => {
    const networkAdUnits = adUnits
      .filter((adUnit) => adUnit.networkKey === network.key)
      .sort((a, b) => a.pricefloor - b.pricefloor);
    // Use initialAdUnitIds for first load, then use adUnitIds for updates
    const idsToUse =
      props.adUnitIds.length > 0 ? props.adUnitIds : props.initialAdUnitIds;

    return {
      ...network,
      enabled: props.networkKeys.includes(network.key),
      adUnits: networkAdUnits,
      selectedAdUnitIds: idsToUse.filter((id) =>
        networkAdUnits.some((unit) => unit.id === id),
      ),
    };
  });
  isLoaded.value = true;
  networks.value = updatedNetworks;
};

watch(
  [
    adUnitsData,
    () => props.networkKeys,
    () => props.adUnitIds,
    () => props.initialAdUnitIds,
  ],
  () => {
    if (adUnitsData.value) {
      updateNetworks();
    }
  },
  { immediate: true },
);

const emitUpdates = () => {
  const enabledNetworks = networks.value.filter((network) => network.enabled);
  const selectedNetworkKeys = enabledNetworks.map((network) => network.key);
  const selectedAdUnitIds = enabledNetworks
    .map((network) => network.selectedAdUnitIds)
    .flat();

  emit("update:networkKeys", selectedNetworkKeys);
  emit("update:adUnitIds", selectedAdUnitIds);
};

watch(
  () =>
    networks.value.map((network) => ({
      key: network.key,
      enabled: network.enabled,
      selectedAdUnitIds: [...network.selectedAdUnitIds],
    })),
  (newNetworks, oldNetworks) => {
    if (!isLoaded.value || !oldNetworks) return;

    let hasChanges = false;

    newNetworks.forEach((newNetwork, index) => {
      const oldNetwork = oldNetworks[index];
      if (!oldNetwork) return;

      if (newNetwork.enabled !== oldNetwork.enabled) {
        const event = newNetwork.enabled
          ? "network-enabled"
          : "network-disabled";
        emit(event, newNetwork.key);
        hasChanges = true;
      }

      const oldIds = JSON.stringify(oldNetwork.selectedAdUnitIds.sort());
      const newIds = JSON.stringify(newNetwork.selectedAdUnitIds.sort());
      if (oldIds !== newIds) {
        hasChanges = true;
      }
    });

    if (hasChanges) {
      emitUpdates();
    }
  },
  { deep: true },
);
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
