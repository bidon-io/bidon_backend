<template>
  <template v-if="apiKey === 'amazon'">
    <AmazonPricePointsForm />
  </template>
  <template v-if="apiKey === 'applovin'">
    <VeeFormFieldWrapper
      field="extra.sdkKey"
      label="SDK Key"
      placeholder="Returned as 'app_key' in SDK API"
      required
    />
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
  <template v-if="apiKey === 'gam'">
    <VeeFormFieldWrapper
      field="extra.networkCode"
      label="Network Code"
      required
    />
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
  <template v-if="apiKey === 'moloco'">
    <!-- Moloco configuration is handled via environment variables -->
    <!-- No additional form fields required -->
  </template>
  <template v-if="apiKey === 'startio'">
    <VeeFormFieldWrapper field="extra.account" label="Account" required />
  </template>
  <template v-if="apiKey === 'vungle'">
    <VeeFormFieldWrapper field="extra.accountId" label="Account ID" required />
  </template>
  <template v-if="apiKey === 'yandex'">
    <VeeFormFieldWrapper
      field="extra.oauth_token"
      label="Oauth Token"
      required
    />
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
  amazon: yup.object({
    pricePointsMap: yup.object().required().label("Price Points"),
  }),
  applovin: yup.object({
    sdkKey: yup.string().required().label("SDK Key"),
  }),
  bidmachine: yup.object({
    sellerId: yup.string().required().label("Seller Id"),
    endpoint: yup
      .string()
      .required()
      .test(
        "is-url-or-host",
        "Must be a valid URL or host",
        (value) =>
          yup.string().url().isValidSync(value) ||
          (value.replaceAll(".", "").length <= 255 &&
            yup
              .string()
              .matches(
                /^([a-zA-Z0-9_][a-zA-Z0-9_-]{0,62})(\.[a-zA-Z0-9_][a-zA-Z0-9_-]{0,62})*[._]?$/,
              )
              .isValidSync(value)), // regexp is from https://github.com/asaskevich/govalidator/blob/a9d515a09cc289c60d55064edec5ef189859f172/patterns.go#L33
      )
      .label("Endpoint"),
    mediationConfig: yup
      .array()
      .of(yup.string())
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
  moloco: yup.object({
    // No additional validation required - API key handled via environment
  }),
  startio: yup.object({
    account: yup.string().required().label("Account"),
  }),
  vungle: yup.object({
    accountId: yup.string().required().label("Account Id"),
  }),
  yandex: yup.object({
    oauth_token: yup.string().required().label("Oauth Token"),
  }),
};

const schema = computed(() => dataSchemas[props.apiKey] || yup.object());
watchEffect(() => {
  emit("update:schema", schema.value);
});
</script>
