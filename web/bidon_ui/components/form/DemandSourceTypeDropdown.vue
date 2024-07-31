<template>
  <FormField :label="label" :error="error" :required="required">
    <Dropdown
      v-model="value"
      :options="options"
      option-label="label"
      option-value="value"
      class="w-full md:w-14rem"
      placeholder="Select Demand source type"
      :disabled="disabled"
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
  label: {
    type: String,
    default: "Demand Source Type",
  },
  modelValue: {
    type: [String, null],
    default: null,
  },
  disabled: {
    type: Boolean,
    default: false,
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
  { label: "Admob", value: "DemandSourceAccount::Admob" },
  { label: "Amazon", value: "DemandSourceAccount::Amazon" },
  { label: "Applovin", value: "DemandSourceAccount::Applovin" },
  { label: "Bidmachine", value: "DemandSourceAccount::BidMachine" },
  { label: "BigoAds", value: "DemandSourceAccount::BigoAds" },
  { label: "DTExchange", value: "DemandSourceAccount::DtExchange" },
  { label: "Google Ad Manager", value: "DemandSourceAccount::GAM" },
  { label: "Inmobi", value: "DemandSourceAccount::Inmobi" },
  { label: "Meta", value: "DemandSourceAccount::Meta" },
  { label: "Mintegral", value: "DemandSourceAccount::Mintegral" },
  { label: "MobileFuse", value: "DemandSourceAccount::MobileFuse" },
  { label: "UnityAds", value: "DemandSourceAccount::UnityAds" },
  { label: "Vungle", value: "DemandSourceAccount::Vungle" },
  { label: "VKAds", value: "DemandSourceAccount::VKAds" },
  { label: "Yandex", value: "DemandSourceAccount::Yandex" },
]);
</script>
