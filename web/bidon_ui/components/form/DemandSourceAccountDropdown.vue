<template>
  <FormField label="Demand Source Account" :error="error" :required="required">
    <Dropdown
      v-model="value"
      :options="options"
      option-label="label"
      option-value="id"
      class="w-full md:w-14rem"
      placeholder="Select Demand source account"
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

const options = ref([]);
axios
  .get("/demand_source_accounts")
  .then((response) => {
    options.value = response.data.map(({ id, type }) => ({
      id,
      label: `${type}:${id}`,
    }));
  })
  .catch((error) => {
    console.error(error);
  });
</script>