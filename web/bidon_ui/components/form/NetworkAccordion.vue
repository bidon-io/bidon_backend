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
      <Fieldset legend="Ad Units">
        <div class="flex flex-col gap-3">
          <div
            v-for="adUnit in network.adUnits"
            :key="adUnit.id"
            class="flex align-items-center"
          >
            <Checkbox v-model="network.selectedAdUnitIds" :value="adUnit.id" />
            <div class="flex gap-2 ml-4 text-sm">
              <span><b>Label:</b> {{ adUnit.label }}</span>
              <span><b>UID:</b> {{ adUnit.uid }}</span>
              <span
                ><b>Price Floor:</b>
                {{ `$${adUnit.pricefloor.toFixed(2)}` }}</span
              >
              <span><b>Account:</b> {{ adUnit.account }}</span>
              <span><b>Bidding:</b> {{ adUnit.isBidding }}</span>
            </div>
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
    label: "UnityAds",
    key: "unityads",
    isBidding: false,
    enabled: false,
    adUnits: [],
    selectedAdUnitIds: [],
  },
  // Bidding networks
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
]);
const isLoaded = ref(false);

const fetchAdUnits = async () => {
  if (!props.appId || !props.adType) return [];
  const url = `/line_items?appId=${props.appId}&adType=${props.adType}&isBidding=${props.isBidding}`;
  const { data, error } = await useAsyncData(url, async () => {
    const response = await axios.get(url);
    // TODO: Filters on API doesn't work, so filtering here
    const result = response.data.filter(
      (
        adUnit: any, // eslint-disable-line @typescript-eslint/no-explicit-any
      ) =>
        adUnit.adType === props.adType &&
        adUnit.appId === props.appId &&
        adUnit.isBidding === props.isBidding,
    );
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
      const networkAdUnits = adUnits.filter(
        (adUnit) => adUnit.networkKey === network.key,
      );
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
