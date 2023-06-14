<template>
  <InputText v-model="value" type="text" :placeholder="placeholder" />
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  placeholder: {
    type: String,
    default: "",
  },
  modelValue: {
    type: Object,
    required: true,
  },
});
const emit = defineEmits(["update:modelValue"]);
const value = computed({
  get: () => JSON.stringify(props.modelValue),
  set: (newValue) => {
    try {
      emit("update:modelValue", JSON.parse(newValue));
    } catch {}
  },
});
</script>
