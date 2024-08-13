<template>
  <FormField label="Owner" :error="error" :required="required">
    <Dropdown
      v-model="value"
      :options="users"
      option-label="email"
      option-value="id"
      class="w-full md:w-14rem"
      placeholder="Select Owner"
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

const users = ref([]);
axios
  .get("/users")
  .then((response) => {
    users.value = response.data;
  })
  .catch((error) => {
    console.error(error);
  });
</script>
