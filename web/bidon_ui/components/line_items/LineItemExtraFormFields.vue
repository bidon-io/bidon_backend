<template>
  <template v-if="apiKey === 'admob'">
    <VeeFormFieldWrapper field="extra.adUnitId" label="Ad Unit Id" required />
  </template>
  <template v-if="apiKey === 'applovin'">
    <VeeFormFieldWrapper field="extra.zoneId" label="Zone Id" required />
  </template>
  <template v-if="apiKey === 'amazon'">
    <LineItemsAmazonExtraFields
      :ad-type-with-format="adTypeWithFormat"
      required
    />
  </template>
  <template v-if="apiKey === 'bigoads'">
    <VeeFormFieldWrapper field="extra.slotId" label="Slot Id" required />
  </template>
  <template v-if="apiKey === 'dtexchange'">
    <VeeFormFieldWrapper field="extra.spotId" label="Spot Id" required />
  </template>
  <template v-if="apiKey === 'gam'">
    <VeeFormFieldWrapper field="extra.adUnitId" label="Ad Unit Id" required />
  </template>
  <template v-if="apiKey === 'inmobi'">
    <VeeFormFieldWrapper
      field="extra.placementId"
      label="Placement Id"
      required
    />
  </template>
  <template v-if="apiKey === 'ironsource'">
    <VeeFormFieldWrapper field="extra.instaceId" label="Instance Id" required />
  </template>
  <template v-if="apiKey === 'meta'">
    <VeeFormFieldWrapper
      field="extra.placementId"
      label="Placement Id"
      required
    />
  </template>
  <template v-if="apiKey === 'mintegral'">
    <VeeFormFieldWrapper
      field="extra.placementId"
      label="Placement Id"
      required
    />
    <VeeFormFieldWrapper field="extra.unitId" label="Unit Id" required />
  </template>
  <template v-if="apiKey === 'mobilefuse'">
    <VeeFormFieldWrapper
      field="extra.placementId"
      label="Placement Id"
      required
    />
  </template>
  <template v-if="apiKey === 'unityads'">
    <VeeFormFieldWrapper
      field="extra.placementId"
      label="Placement Id"
      required
    />
  </template>
  <template v-if="apiKey === 'vungle'">
    <VeeFormFieldWrapper
      field="extra.placementId"
      label="Placement Id"
      required
    />
  </template>
  <template v-if="apiKey === 'vkads'">
    <VeeFormFieldWrapper field="extra.slotId" label="Slot Id" required />
    <VeeFormFieldWrapper field="extra.mediation" label="Mediation" />
  </template>
  <template v-if="apiKey === 'yandex'">
    <VeeFormFieldWrapper field="extra.adUnitId" label="Ad Unit Id" />
  </template>
</template>

<script setup>
import * as yup from "yup";

const props = defineProps({
  schema: {
    type: Object,
    required: true,
  },
  apiKey: {
    type: String,
    required: true,
  },
  adType: {
    type: String,
    required: true,
  },
  adTypeWithFormat: {
    type: Object,
    required: true,
  },
});
const emit = defineEmits(["update:schema"]);

const dataSchemas = {
  admob: yup.object({
    adUnitId: yup.string().required().label("Ad Unit Id"),
  }),
  applovin: yup.object({
    zoneId: yup.string().required().label("Zone Id"),
  }),
  amazon: yup.object({
    slotUuid: yup.string().required().label("Slot Uuid"),
    isVideo: yup.boolean().label("Is Video"),
  }),
  bidmachine: yup.object({}),
  bigoads: yup.object({
    slotId: yup.string().required().label("Slot Id"),
  }),
  dtexchange: yup.object({
    spotId: yup.string().required().label("Spot Id"),
  }),
  gam: yup.object({
    adUnitId: yup.string().required().label("Ad Unit Id"),
  }),
  inmobi: yup.object({
    placementId: yup.string().required().label("Placement Id"),
  }),
  ironsource: yup.object({
    instaceId: yup.string().required().label("Instance Id"),
  }),
  meta: yup.object({
    placementId: yup.string().required().label("Placement Id"),
  }),
  mintegral: yup.object({
    placementId: yup.string().required().label("Placement Id"),
    unitId: yup.string().required().label("Ad Unit Id"),
  }),
  mobilefuse: yup.object({
    placementId: yup.string().required().label("Placement Id"),
  }),
  unityads: yup.object({
    placementId: yup.string().required().label("Placement Id"),
  }),
  vungle: yup.object({
    placementId: yup.string().required().label("Placement Id"),
  }),
  vkads: yup.object({
    slotId: yup.string().required().label("Slot Id"),
    mediation: yup.string().label("Mediation"),
  }),
  yandex: yup.object({
    adUnitId: yup.string().label("Ad Unit Id"),
  }),
};

const schema = computed(() => dataSchemas[props.apiKey] || yup.object());
watchEffect(() => {
  emit("update:schema", schema.value);
});
</script>
