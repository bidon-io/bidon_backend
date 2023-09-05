<template>
  <FormField label="Demand Source Type" :error="error" :required="required">
    <Dropdown
      v-model="value"
      :options="options"
      option-label="label"
      option-value="value"
      class="w-full md:w-14rem"
      placeholder="Select Demand source type"
    />
  </FormField>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  error: {
    type: String,
    default: "",
  },
  required: {
    type: Boolean,
    default: false,
  },
  modelValue: {
    type: [String, null],
    default: null,
  },
});
const emit = defineEmits(["update:modelValue"]);

const value = computed({
  get() {
    return props.modelValue;
  },
  set(value) {
    emit("update:modelValue", value);
  },
});
const options = ref([
  { label: "BidmachineAccount", value: "DemandSourceAccount::BidMachine" },
  { label: "AdmobAccount", value: "DemandSourceAccount::Admob" },
  { label: "ApplovinAccount", value: "DemandSourceAccount::Applovin" },
  { label: "DTExchangeAccount", value: "DemandSourceAccount::DtExchange" },
  { label: "UnityAdsAccount", value: "DemandSourceAccount::UnityAds" },
  { label: "VungleAccount", value: "DemandSourceAccount::Vungle" },
  { label: "MintegralAccount", value: "DemandSourceAccount::Mintegral" },
  { label: "BigoAdsAccount", value: "DemandSourceAccount::BigoAds" },
  { label: "MetaAccount", value: "DemandSourceAccount::Meta" },
  { label: "InmobiAccount", value: "DemandSourceAccount::Inmobi" },
]);
</script>
