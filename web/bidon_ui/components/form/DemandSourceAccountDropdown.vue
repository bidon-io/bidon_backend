<template>
  <FormField label="Demand Source Account" :error="error" :required="required">
    <Dropdown
      v-model="value"
      :options="options"
      option-label="label"
      option-value="id"
      class="w-full md:w-14rem"
      :disabled="disabled"
      placeholder="Select Demand source account"
      filter
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
  disabled: {
    type: Boolean,
    default: false,
  },
  accounts: {
    type: Array,
    default: () => [],
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

const buildOptions = (accounts) =>
  accounts.map(({ id, type, label }) => ({
    id,
    label: `(${type.split("::")[1]}) ${label ? label : `#${id}`}`,
  }));

const fetchAccounts = async () => {
  const response = await axios.get("/demand_source_accounts");
  return response.data;
};

// if accounts are passed as props, use them, otherwise fetch them
const options = ref([]);
if (props.accounts.length > 0) {
  options.value = buildOptions(props.accounts);
} else {
  const accounts = await fetchAccounts();
  options.value = buildOptions(accounts);
}

// if accounts are passed as props, watch them for changes
watch(
  () => props.accounts,
  (accounts) => {
    options.value = buildOptions(accounts);
  },
);
</script>
