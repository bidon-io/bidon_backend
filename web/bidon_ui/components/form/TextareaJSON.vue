<template>
  <Textarea
    v-model="value"
    type="text"
    :placeholder="placeholder"
    :rows="rows"
    :cols="cols"
  />
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  modelValue: {
    type: [Object, null],
    default: null,
  },
  placeholder: {
    type: String,
    default: "",
  },
  rows: {
    type: [Number, String],
    required: false,
    default: 5,
  },
  cols: {
    type: [Number, String],
    required: false,
    default: 50,
  },
});
const emit = defineEmits(["update:modelValue"]);
const value = computed({
  get: () => JSON.stringify(props.modelValue),
  set: (newValue) => {
    try {
      emit("update:modelValue", JSON.parse(newValue, null, 2));
    } catch {}
  },
});
</script>
