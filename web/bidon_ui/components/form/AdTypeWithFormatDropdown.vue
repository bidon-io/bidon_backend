<template>
  <FormField label="Ad Type">
    <Dropdown
      v-model="value"
      option-label="label"
      option-value="value"
      :options="options"
      class="w-full md:w-14rem"
      placeholder="Select Ad Format"
    />
  </FormField>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  modelValue: {
    type: [Object, null],
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
  { label: "Adaptive banner", value: { adType: "banner", format: "ADAPTIVE" } },
  { label: "Banner", value: { adType: "banner", format: "BANNER" } },
  { label: "MREC", value: { adType: "banner", format: "MREC" } },
  { label: "Interstitial", value: { adType: "interstitial", format: "" } },
  { label: "Rewarded", value: { adType: "rewarded", format: "" } },
]);
</script>
