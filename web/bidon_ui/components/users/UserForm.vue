<template>
  <form @submit="onSubmit">
    <FormCard title="User">
      <FormField label="Email" :error="errors.email" required>
        <InputText v-model="email" type="string" placeholder="Email" />
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
    email: yup.string().required().label("Email"),
  }),
  initialValues: {
    email: resource.value.email || "",
  },
});

const email = useFieldModel("email");

const onSubmit = handleSubmit((values) => emit("submit", values));
</script>
