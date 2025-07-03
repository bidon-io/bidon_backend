<template>
  <form @submit="onSubmit">
    <FormCard title="Auction Configuration">
      <div
        v-if="showCopySettings"
        class="mb-4 p-4 bg-blue-50 border border-blue-200 rounded-lg"
      >
        <h3 class="text-sm font-semibold text-blue-900 mb-3">
          Clone Settings from Another Configuration
        </h3>
        <div class="flex gap-3 items-start">
          <div class="flex-1">
            <label class="block text-xs font-medium text-gray-700 mb-1"
              >Auction Key</label
            >
            <InputText
              v-model="copyAuctionKey"
              type="text"
              placeholder="Enter auction key (e.g., 1HVR32MFO0400)"
              class="w-full"
              size="small"
            />
            <div v-if="copyError" class="text-xs text-red-600 mt-1">
              {{ copyError }}
            </div>
          </div>
          <div class="pt-5">
            <Button
              type="button"
              label="Clone"
              icon="pi pi-copy"
              size="small"
              :loading="copyLoading"
              :disabled="!copyAuctionKey || copyLoading"
              @click="copySettings"
            />
          </div>
        </div>
      </div>

      <FormField label="Name" :error="errors.name" required>
        <InputText v-model="name" type="text" placeholder="Name" />
      </FormField>
      <AppDropdown v-model="appId" :error="errors.appId" required />
      <AdTypeDropdown v-model="adType" :error="errors.adType" required />

      <SegmentDropdown v-model="segmentId" :error="errors.segmentId" />
      <FormField label="Price floor" :error="errors.pricefloor" required>
        <InputNumber
          v-model="pricefloor"
          input-id="pricefloor"
          :min-fraction-digits="2"
          :max-fraction-digits="5"
          placeholder="Price floor"
        />
      </FormField>
      <FormField label="Is Default" :error="errors.isDefault">
        <Checkbox v-model="isDefault" :binary="true" />
      </FormField>
      <FormField
        label="External Win Notification"
        :error="errors.externalWinNotifications"
      >
        <Checkbox v-model="externalWinNotifications" :binary="true" />
      </FormField>
      <FormField label="Timeout" :error="errors.timeout" required>
        <InputNumber
          v-model="timeout"
          input-id="timeout"
          placeholder="Timeout (ms)"
        />
      </FormField>
      <!-- App Demand Profile Validation Warnings -->
      <div v-if="hasValidationWarnings" class="mb-4">
        <transition-group name="p-message" tag="div">
          <Message
            v-for="(warning, index) in validationWarnings"
            :key="index"
            :severity="warning.severity"
            class="mb-2"
          >
            {{ warning.message }}
          </Message>
        </transition-group>
      </div>

      <FormField v-if="showNetworks" label="CPM Networks">
        <NetworkAccordion
          v-model:network-keys="demands"
          v-model:ad-unit-ids="demandAdUnitIds"
          :ad-type="adType"
          :app-id="appId"
          :is-bidding="false"
          @network-enabled="onNetworkEnabled"
          @network-disabled="onNetworkDisabled"
        />
      </FormField>
      <FormField v-if="showNetworks" label="Bidding Networks">
        <NetworkAccordion
          v-model:network-keys="bidding"
          v-model:ad-unit-ids="biddingAdUnitIds"
          :ad-type="adType"
          :app-id="appId"
          :is-bidding="true"
          @network-enabled="onNetworkEnabled"
          @network-disabled="onNetworkDisabled"
        />
      </FormField>
      <FormSubmitButton />
    </FormCard>
  </form>
</template>

<script setup>
import * as yup from "yup";
import { useToast } from "primevue/usetoast";
import axios from "@/services/ApiService.js";
import { useAppDemandProfileValidation } from "@/composables/useAppDemandProfileValidation.js";

const props = defineProps({
  value: {
    type: Object,
    required: true,
  },
});
const emit = defineEmits(["submit"]);
const resource = ref(props.value);

const { errors, useFieldModel, handleSubmit } = useForm({
  validationSchema: yup.object({
    name: yup.string().required().label("Name"),
    appId: yup.number().required().label("App Id"),
    adType: yup.string().required().label("AdType"),
    pricefloor: yup.number().positive().required().label("Pricefloor"),
    segmentId: yup.number().nullable(true).label("Segment Id"),
    isDefault: yup.boolean(),
    externalWinNotifications: yup.boolean(),
    timeout: yup.number().positive().required().label("Timeout"),
    settings: yup.object(),
  }),
  initialValues: {
    name: resource.value.name || "",
    appId: resource.value.appId || null,
    adType: resource.value.adType || "",
    pricefloor: resource.value.pricefloor || null,
    segmentId: resource.value.segmentId || null,
    isDefault: resource.value.isDefault || false,
    externalWinNotifications:
      resource.value.externalWinNotifications !== undefined
        ? resource.value.externalWinNotifications
        : true,
    timeout: resource.value.timeout || null,
    settings: resource.value.settings || {},
  },
});

const name = useFieldModel("name");
const appId = useFieldModel("appId");
const adType = useFieldModel("adType");
const pricefloor = useFieldModel("pricefloor");
const segmentId = useFieldModel("segmentId");
const isDefault = useFieldModel("isDefault");
const externalWinNotifications = useFieldModel("externalWinNotifications");
const timeout = useFieldModel("timeout");

const demands = ref(resource.value.demands || []);
const bidding = ref(resource.value.bidding || []);
const demandAdUnitIds = ref(resource.value.adUnitIds || []);
const biddingAdUnitIds = ref(resource.value.adUnitIds || []);

const showNetworks = computed(() => appId.value && adType.value);
const adUnitIds = computed(() => [
  ...new Set(demandAdUnitIds.value.concat(biddingAdUnitIds.value)),
]);

// App Demand Profile Validation
const {
  validateNetworkEnabled,
  validateNetworkDisabled,
  clearAllWarnings,
  warningMessages: validationWarnings,
  hasWarnings: hasValidationWarnings,
} = useAppDemandProfileValidation(appId);

const onNetworkEnabled = async (networkApiKey) => {
  await validateNetworkEnabled(networkApiKey);
};

const onNetworkDisabled = (networkApiKey) => {
  validateNetworkDisabled(networkApiKey);
};

// Clear warnings when app changes
watch(appId, () => {
  clearAllWarnings();
});

const toast = useToast();
const copyAuctionKey = ref("");
const copyLoading = ref(false);
const copyError = ref("");

const showCopySettings = computed(() => !!resource.value.id);

const fetchSourceConfig = async (auctionKey) => {
  const response = await axios.get("/v2/auction_configurations_collection", {
    params: { auction_key: auctionKey, limit: 1 },
  });
  return response.data?.items?.[0];
};

const validateSourceConfig = (sourceConfig) => {
  if (!sourceConfig) {
    throw new Error("No auction configuration found with this auction key");
  }
  if (sourceConfig.appId !== appId.value) {
    throw new Error("Source configuration must belong to the same app");
  }
  if (sourceConfig.adType !== adType.value) {
    throw new Error("Source configuration must have the same ad type");
  }
};

const updateFormFields = (sourceConfig) => {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const { name, isDefault, ...settings } = sourceConfig;

  pricefloor.value = settings.pricefloor;
  segmentId.value = settings.segmentId;
  externalWinNotifications.value = settings.externalWinNotifications;
  timeout.value = settings.timeout;

  demands.value = settings.demands || [];
  bidding.value = settings.bidding || [];
  demandAdUnitIds.value = settings.adUnitIds || [];
  biddingAdUnitIds.value = settings.adUnitIds || [];
};

const copySettings = async () => {
  if (!copyAuctionKey.value) return;

  copyLoading.value = true;
  copyError.value = "";

  try {
    const sourceConfig = await fetchSourceConfig(copyAuctionKey.value);
    validateSourceConfig(sourceConfig);
    updateFormFields(sourceConfig);

    copyAuctionKey.value = "";
    toast.add({
      severity: "success",
      summary: "Success",
      detail: "Settings cloned successfully!",
      life: 3000,
    });
  } catch (error) {
    copyError.value = error.message || "Failed to fetch source configuration";
  } finally {
    copyLoading.value = false;
  }
};

const onSubmit = handleSubmit((values) =>
  emit("submit", {
    ...values,
    adUnitIds: adUnitIds.value,
    demands: demands.value,
    bidding: bidding.value,
  }),
);
</script>
