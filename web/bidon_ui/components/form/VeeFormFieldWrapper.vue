<template>
  <FormField :label="label" :error="errorMessage" :required="required">
    <InputText v-model="value" :type="type" :placeholder="label" />
  </FormField>
</template>

<script setup>
import { useField } from "vee-validate";

const props = defineProps({
  label: {
    type: String,
    required: true,
  },
  field: {
    type: String,
    required: true,
  },
  required: {
    type: Boolean,
    default: false,
  },
  type: {
    type: String,
    default: "text",
  },
});

const { value, errorMessage } = useField(props.field);
</script>
