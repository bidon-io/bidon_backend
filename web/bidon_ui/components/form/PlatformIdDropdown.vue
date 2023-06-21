<template>
  <FormField label="Platform Id">
    <Dropdown v-model="value" :options="platforms" class="w-full md:w-14rem" placeholder="Select Platform" />
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
const platforms = ref(["ios", "android"]);
</script>
