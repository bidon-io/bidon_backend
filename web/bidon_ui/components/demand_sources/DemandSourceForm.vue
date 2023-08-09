<template>
  <form @submit="onSubmit">
    <FormCard title="Demand source">
      <FormField label="Human Name" :error="errors.humanName" required>
        <InputText v-model="humanName" type="text" placeholder="Name" />
      </FormField>
      <FormField label="Api Key" :error="errors.apiKey" required>
        <InputText v-model="apiKey" type="text" placeholder="Api Key" />
      </FormField>
      <FormSubmitButton />
    </FormCard>
  </form>
</template>

<script setup>
import * as yup from "yup";

const props = defineProps({
  value: {
    type: Object,
    required: true,
  },
});
const emit = defineEmits(["submit"]);
const resource = ref(props.value);

const { errors, useFieldModel, handleSubmit } = useForm({
  validationSchema: yup.object({
    humanName: yup.string().required().label("Human Name"),
    apiKey: yup.string().required().label("Api Key"),
  }),
  initialValues: {
    humanName: resource.value.humanName || "",
    apiKey: resource.value.apiKey || "",
  },
});

const humanName = useFieldModel("humanName");
const apiKey = useFieldModel("apiKey");

const onSubmit = handleSubmit((values) => emit("submit", values));
</script>
