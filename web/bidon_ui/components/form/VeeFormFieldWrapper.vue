<template>
  <FormField :label="label" :error="errorMessage" :required="required">
    <InputText
      v-if="type === 'text'"
      v-model="value"
      :type="type"
      :placeholder="label"
    />
    <Checkbox
      v-if="type === 'bool'"
      v-model="value"
      :placeholder="label"
      binary
    />
    <TextareaJSON
      v-if="type === 'array'"
      v-model="value"
      :placeholder="label"
    />
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
    validator: (value) => ["text", "bool", "json"].includes(value),
    default: "text",
  },
});

const { value, errorMessage } = useField(props.field);
if (props.type === "array" && !value.value) value.value = [];
</script>
