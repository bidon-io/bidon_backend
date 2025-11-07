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
      filter
    />
  </FormField>
</template>

<script setup>
import { computed } from "vue";
import { DEMAND_SOURCE_OPTIONS } from "@/constants/DemandSourceOptions";

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
const options = ref(DEMAND_SOURCE_OPTIONS);
</script>
