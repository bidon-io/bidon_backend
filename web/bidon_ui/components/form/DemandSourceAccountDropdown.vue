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

const options = ref([]);
if (props.accounts.length > 0) {
  options.value = buildOptions(props.accounts);
} else {
  axios
    .get("/demand_source_accounts")
    .then((response) => {
      options.value = buildOptions(response.data);
    })
    .catch((error) => {
      console.error(error);
    });
}
</script>
