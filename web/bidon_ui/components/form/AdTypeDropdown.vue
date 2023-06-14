<template>
  <FormField lable="Ad Type">
    <Dropdown v-model="value" :options="adTypes" class="w-full md:w-14rem" placeholder="Select Ad Type" />
  </FormField>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  modelValue: {
    type: String,
    required: true,
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
const adTypes = ref(["banner", "interstitial", "rewarded"]);
</script>
