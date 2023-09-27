<template>
  <template v-if="apiKey === 'applovin'">
    <VeeFormFieldWrapper field="extra.sdkKey" label="SDK Key" required />
  </template>
  <template v-if="apiKey === 'bidmachine'">
    <VeeFormFieldWrapper field="extra.sellerId" label="Seller ID" required />
    <VeeFormFieldWrapper field="extra.endpoint" label="Endpoint" required />
    <VeeFormFieldWrapper
      field="extra.mediationConfig"
      label="Mediation Config"
      type="array"
      required
    />
  </template>
  <template v-if="apiKey === 'bigoads'">
    <VeeFormFieldWrapper
      field="extra.publisherId"
      label="Publisher ID"
      required
    />
    <VeeFormFieldWrapper field="extra.endpoint" label="Endpoint" required />
  </template>
  <template v-if="apiKey === 'inmobi'">
    <VeeFormFieldWrapper field="extra.accountId" label="Account ID" required />
  </template>
  <template v-if="apiKey === 'mintegral'">
    <VeeFormFieldWrapper field="extra.appKey" label="App Key" required />
    <VeeFormFieldWrapper
      field="extra.publisherId"
      label="Publisher ID"
      required
    />
  </template>
  <template v-if="apiKey === 'mobilefuse'">
    <VeeFormFieldWrapper
      field="extra.publisherId"
      label="Publisher ID"
      required
    />
  </template>
  <template v-if="apiKey === 'vungle'">
    <VeeFormFieldWrapper field="extra.accountId" label="Account ID" required />
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
});
const emit = defineEmits(["update:schema"]);

const dataSchemas = {
  applovin: yup.object({
    sdkKey: yup.string().required().label("SDK Key"),
  }),
  bidmachine: yup.object({
    sellerId: yup.string().required().label("Seller Id"),
    endpoint: yup.string().url().required().label("Endpoint"),
    mediationConfig: yup
      .array()
      .of(yup.string())
      .min(1)
      .required()
      .label("Mediation Config"),
  }),
  bigoads: yup.object({
    publisherId: yup.string().required().label("Publisher Id"),
    endpoint: yup.string().url().required().label("Endpoint"),
  }),
  inmobi: yup.object({
    accountId: yup.string().required().label("Account Id"),
  }),
  mintegral: yup.object({
    appKey: yup.string().required().label("App Key"),
    publisherId: yup.string().required().label("Publisher Id"),
  }),
  mobilefuse: yup.object({
    publisherId: yup.string().required().label("Publisher Id"),
  }),
  vungle: yup.object({
    accountId: yup.string().required().label("Account Id"),
  }),
};

const schema = computed(() => dataSchemas[props.apiKey] || yup.object());
watchEffect(() => {
  emit("update:schema", schema.value);
});
</script>
