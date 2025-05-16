<template>
  <div class="mb-8">
    <!-- Bidding Networks Section -->
    <div v-if="biddingNetworks.length > 0" class="mb-6">
      <h2 class="text-xl font-semibold mb-4">Bidding Networks</h2>
      <div class="grid grid-cols-3 gap-4">
        <div
          v-for="network in biddingNetworks"
          :key="network.name"
          class="p-4 bg-gray-100 rounded-md flex items-center justify-start text-left border border-gray-200 min-h-[60px]"
        >
          <div class="flex items-center justify-center w-[30px] mr-2">
            <i :class="getNetworkIcon(true)" class="text-xl"></i>
          </div>
          <div>{{ formatNetworkName(network.name) }}</div>
        </div>
      </div>
    </div>

    <!-- Waterfall Networks Section -->
    <div v-if="cpmNetworks && cpmNetworks.length > 0" class="mt-8 mb-6">
      <h2 class="text-xl font-semibold mb-4">Waterfall Networks</h2>
      <div class="flex flex-col gap-2">
        <div
          v-for="(network, index) in cpmNetworks"
          :key="network.id"
          class="p-4 bg-gray-100 rounded-md flex items-center"
          :class="{
            'mb-2': index !== cpmNetworks.length - 1,
            'border border-gray-200': true,
          }"
        >
          <div class="flex items-center justify-center w-[30px] mr-3">
            <i :class="getNetworkIcon(network.isBidding)" class="text-xl"></i>
          </div>
          <div class="flex-grow">
            {{ formatNetworkName(network.name) }} -
            <span class="font-medium">${{ network.bidFloor }}</span>
          </div>
        </div>
      </div>
    </div>

    <div v-if="isLoading" class="mt-4 flex items-center">
      <i class="pi pi-spin pi-spinner mr-2"></i>
      <p>Loading networks data...</p>
    </div>

    <div
      v-if="
        !isLoading &&
        biddingNetworks.length === 0 &&
        (!cpmNetworks || cpmNetworks.length === 0)
      "
      class="mt-4"
    >
      <p class="text-gray-500">
        No networks configured for this auction configuration.
      </p>
    </div>
  </div>
</template>

<script setup>
import { computed } from "vue";
import { useAsyncData } from "#app";

const props = defineProps({
  appId: {
    type: Number,
    required: true,
  },
  adType: {
    type: String,
    required: true,
  },
  bidding: {
    type: Array,
    default: () => [],
  },
  adUnitIds: {
    type: Array,
    default: () => [],
  },
});

const biddingNetworks = computed(() => {
  if (!Array.isArray(props.bidding)) {
    return [];
  }

  const networks = props.bidding.map((networkKey) => {
    return { name: networkKey };
  });
  return networks;
});

// Helper method to get appropriate icon based on network type
const getNetworkIcon = (isBidding = false) => {
  // Use different icons for bidding vs waterfall networks
  if (isBidding) {
    return "pi pi-bolt text-blue-500"; // Bidding networks icon
  } else {
    return "pi pi-tag text-green-500"; // Waterfall networks icon
  }
};

// Helper method to format network names for display
const formatNetworkName = (networkName) => {
  // Special case formatting for known networks
  const specialCases = {
    admob: "AdMob",
    applovin: "AppLovin",
    bidmachine: "BidMachine",
    bigoads: "BigoAds",
    chartboost: "Chartboost",
    dtexchange: "DT Exchange",
    gam: "Google Ad Manager",
    mintegral: "Mintegral",
    unityads: "Unity Ads",
    ironsource: "ironSource",
    vkads: "VK Ads",
    vungle: "Vungle",
    yandex: "Yandex",
    amazon: "Amazon",
    meta: "Meta",
    facebook: "Facebook",
    mobilefuse: "MobileFuse",
    inmobi: "InMobi",
  };

  const lowerName = networkName.toLowerCase();
  if (specialCases[lowerName]) {
    return specialCases[lowerName];
  }

  // Default formatting: capitalize each word
  return networkName
    .split(/(?=[A-Z])/)
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(" ");
};

// Fetch line items using useAsyncData
const { data: cpmNetworks, status } = useAsyncData(
  `line-items-${props.appId}-${props.adType}-${props.adUnitIds?.join("-")}`,
  async () => {
    if (
      !props.appId ||
      !props.adType ||
      !Array.isArray(props.adUnitIds) ||
      props.adUnitIds.length === 0
    ) {
      return [];
    }

    try {
      const url = `/line_items?app_id=${props.appId}&ad_type=${props.adType}`;
      const lineItems = await $apiFetch(url);
      const adUnitIdsAsStrings = props.adUnitIds.map((id) => String(id));

      const filteredLineItems = lineItems.filter((item) => {
        const itemIdStr = String(item.id);
        return (
          adUnitIdsAsStrings.includes(itemIdStr) && item.isBidding === false
        );
      });

      const processedNetworks = filteredLineItems.map((item) => {
        const networkName = item.accountType.split("::")[1] || "Unknown";

        return {
          id: item.id,
          name: networkName,
          bidFloor: parseFloat(item.bidFloor).toFixed(2),
          isBidding: false,
        };
      });

      processedNetworks.sort(
        (a, b) => parseFloat(b.bidFloor) - parseFloat(a.bidFloor),
      );
      return processedNetworks;
    } catch (error) {
      console.error("Error processing line items:", error);
      return [];
    }
  },
);

const isLoading = computed(() => status.value === "pending");
</script>
