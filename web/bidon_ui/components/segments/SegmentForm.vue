<template>
  <form @submit="onSubmit">
    <FormCard title="Segment">
      <FormField label="Name" :error="errors.name" required>
        <InputText v-model="name" type="text" placeholder="Name" />
      </FormField>
      <FormField label="Description" :error="errors.description">
        <Textarea v-model="description" rows="5" cols="50" />
      </FormField>
      <FormField label="Filters" :error="errors.filters">
        <InputJSON v-model="filters" placeholder="Filters" />
      </FormField>
      <FormField label="Enabled" :error="errors.enabled">
        <Checkbox v-model="enabled" :binary="true" />
      </FormField>
      <AppDropdown v-model="appId" :error="errors.appId" />
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
    name: yup.string().required().label("Name"),
    description: yup.string(),
    filters: yup.array().required().label("Filters"),
    enabled: yup.boolean(),
    appId: yup.number().required().label("App Id"),
  }),
  initialValues: {
    name: resource.value.name || "",
    description: resource.value.description || "",
    filters: resource.value.filters || [],
    enbled: resource.value.enabled || false,
    appId: resource.value.appId || null,
  },
});

const name = useFieldModel("name");
const description = useFieldModel("description");
const filters = useFieldModel("filters");
const enabled = useFieldModel("enabled");
const appId = useFieldModel("appId");

const onSubmit = handleSubmit((values) => emit("submit", values));
</script>
