<template>
  <FormField label="Segment" :error="error" :required="required">
    <Dropdown
      v-model="value"
      :options="segments"
      option-label="name"
      option-value="id"
      class="w-full md:w-14rem"
      placeholder="None"
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

const segments = ref([]);
axios
  .get("/segments")
  .then((response) => {
    segments.value = [{ name: "None", id: null }, ...response.data];
  })
  .catch((error) => {
    console.error(error);
  });
</script>
