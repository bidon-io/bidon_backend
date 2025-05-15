<template>
  <form @submit="onSubmit">
    <FormCard title="Auction Configuration">
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
      <FormField v-if="showNetworks" label="CPM Networks">
        <NetworkAccordion
          v-model:network-keys="demands"
          v-model:ad-unit-ids="demandAdUnitIds"
          :ad-type="adType"
          :app-id="appId"
          :is-bidding="false"
        />
      </FormField>
      <FormField v-if="showNetworks" label="Bidding Networks">
        <NetworkAccordion
          v-model:network-keys="bidding"
          v-model:ad-unit-ids="biddingAdUnitIds"
          :ad-type="adType"
          :app-id="appId"
          :is-bidding="true"
        />
      </FormField>
      <FormSubmitButton />
    </FormCard>
  </form>
</template>

<script setup>
import * as yup from "yup";

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

const onSubmit = handleSubmit((values) =>
  emit("submit", {
    ...values,
    adUnitIds: adUnitIds.value,
    demands: demands.value,
    bidding: bidding.value,
  }),
);
</script>
