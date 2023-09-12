<template>
  <form @submit="onSubmit">
    <FormCard title="User">
      <FormField label="Email" :error="errors.email" required>
        <InputText v-model="email" type="string" placeholder="Email" />
      </FormField>
      <FormField label="Is Admin">
        <Checkbox v-model="isAdmin" :binary="true" />
      </FormField>
      <FormField label="Password" :error="errors.password" required>
        <InputText v-model="password" type="password" placeholder="Password" />
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
    isAdmin: yup.boolean(),
    password: yup.string().required().label("Password"),
  }),
  initialValues: {
    email: resource.value.email || "",
    password: "",
    isAdmin: resource.value.isAdmin || false,
  },
});

const email = useFieldModel("email");
const isAdmin = useFieldModel("isAdmin");
const password = useFieldModel("password");

const onSubmit = handleSubmit((values) => emit("submit", values));
</script>
