<template>
  <FormField label="Format">
    <Dropdown
      v-model="value"
      :options="options"
      option-label="label"
      option-value="value"
      class="w-full md:w-14rem"
      placeholder="N/A"
    />
  </FormField>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
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
  { label: "N/A", value: null },
  { label: "BANNER", value: "BANNER" },
  { label: "LEADERBOARD", value: "LEADERBOARD" },
  { label: "MREC", value: "MREC" },
  { label: "ADAPTIVE", value: "ADAPTIVE" },
]);
</script>
