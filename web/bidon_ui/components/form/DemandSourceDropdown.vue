<template>
  <FormField
    v-show="show"
    label="Demand source"
    :error="error"
    :required="required"
  >
    <Dropdown
      v-model="value"
      :options="options"
      option-label="humanName"
      option-value="id"
      class="w-full md:w-14rem"
      placeholder="Select Demand source"
    />
  </FormField>
</template>

<script setup>
import { computed } from "vue";
import axios from "@/services/ApiService";

const props = defineProps({
  error: {
    type: String,
    default: "",
  },
  required: {
    type: Boolean,
    default: false,
  },
  show: {
    type: Boolean,
    default: true,
  },
  selectedApiKey: {
    type: String,
    default: "",
  },
  modelValue: {
    type: [Number, null],
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

watch(
  () => props.selectedApiKey,
  (apiKey) => {
    if (!apiKey) {
      return;
    }
    const demandSource = options.value.find((o) => o.apiKey === apiKey);
    emit("update:modelValue", demandSource.id);
  }
);

const options = ref([]);
axios
  .get("/demand_sources")
  .then((response) => {
    options.value = response.data;
  })
  .catch((error) => {
    console.error(error);
  });
</script>
