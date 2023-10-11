<template>
  <FormField v-if="appIdVisible" label="App Id" :error="appIdError" required>
    <InputText v-model="appId" type="text" placeholder="App ID" />
  </FormField>
  <FormField
    v-if="appSecretVisible"
    label="App Secret"
    :error="appSecretError"
    required
  >
    <InputText v-model="appSecret" type="text" placeholder="App Secret" />
  </FormField>
  <FormField v-if="gameIdVisible" label="Game Id" :error="gameIdError" required>
    <InputText v-model="gameId" type="text" placeholder="Game ID" />
  </FormField>
  <FormField v-if="appKeyVisible" label="App key" :error="appKeyError" required>
    <InputText v-model="appKey" type="text" placeholder="App Key" />
  </FormField>
</template>

<script setup>
import { useField } from "vee-validate";
import * as yup from "yup";

const props = defineProps({
  schema: {
    type: Object,
    required: true,
  },
  accountType: {
    type: String,
    required: true,
  },
});
const emit = defineEmits(["update:schema"]);

const dataSchemas = {
  "DemandSourceAccount::Admob": yup.object({
    appId: yup.string().required().label("App Id"),
  }),
  "DemandSourceAccount::BigoAds": yup.object({
    appId: yup.number().required().label("App Id"),
  }),
  "DemandSourceAccount::DtExchange": yup.object({
    appId: yup.number().required().label("App Id"),
  }),
  "DemandSourceAccount::Vungle": yup.object({
    appId: yup.string().required().label("App Id"),
  }),
  "DemandSourceAccount::Meta": yup.object({
    appId: yup.number().required().label("App Id"),
    appSecret: yup.string().required().label("App Secret"),
  }),
  "DemandSourceAccount::Mintegral": yup.object({
    appId: yup.number().required().label("App Id"),
  }),
  "DemandSourceAccount::UnityAds": yup.object({
    gameId: yup.number().required().label("Game Id"),
  }),
  "DemandSourceAccount::Inmobi": yup.object({
    appKey: yup.string().required().label("App Key"),
  }),
  "DemandSourceAccount::MobileFuse": yup.object({
    appKey: yup.string().required().label("App Key"),
  }),
};

const appIdVisible = computed(() =>
  [
    "DemandSourceAccount::Admob",
    "DemandSourceAccount::BigoAds",
    "DemandSourceAccount::DtExchange",
    "DemandSourceAccount::Vungle",
    "DemandSourceAccount::Meta",
    "DemandSourceAccount::Mintegral",
  ].includes(props.accountType),
);
const appSecretVisible = computed(
  () => props.accountType === "DemandSourceAccount::Meta",
);
const gameIdVisible = computed(
  () => props.accountType === "DemandSourceAccount::UnityAds",
);
const appKeyVisible = computed(() =>
  [
    "DemandSourceAccount::Inmobi",
    "DemandSourceAccount::MobileFuse",
    "DemandSourceAccount::Amazon",
  ].includes(props.accountType),
);

const { value: appId, errorMessage: appIdError } = useField("data.appId");
const { value: appSecret, errorMessage: appSecretError } =
  useField("data.appSecret");
const { value: gameId, errorMessage: gameIdError } = useField("data.gameId");
const { value: appKey, errorMessage: appKeyError } = useField("data.appKey");

const schema = computed(() => dataSchemas[props.accountType] || yup.object());

watchEffect(() => {
  emit("update:schema", schema.value);
});
</script>
