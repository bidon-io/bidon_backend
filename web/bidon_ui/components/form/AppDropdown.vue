<template>
  <FormField label="App" :required="required">
    <Dropdown
      v-model="value"
      :options="apps"
      option-label="label"
      option-value="id"
      class="w-full md:w-14rem"
      placeholder="Select App"
      filter
    />
  </FormField>
</template>

<script setup>
import { computed } from "vue";
import axios from "@/services/ApiService";

const props = defineProps({
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

const apps = ref([]);
axios
  .get("/apps")
  .then((response) => {
    apps.value = response.data.map((app) => ({
      ...app,
      label: `${app.packageName} (${app.platformId})`,
    }));
  })
  .catch((error) => {
    console.error(error);
  });
</script>
