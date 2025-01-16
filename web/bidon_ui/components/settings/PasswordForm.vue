<template>
  <transition-group name="p-message" tag="div">
    <Message v-for="(msg, index) in errorMsgs" :key="index" severity="error">{{
      msg
    }}</Message>
  </transition-group>
  <form @submit="onSubmit">
    <FormCard title="Password">
      <FormField
        label="Current Password"
        :error="errors.currentPassword"
        required
      >
        <InputText
          v-model="currentPassword"
          type="password"
          placeholder="Current Password"
        />
      </FormField>
      <FormField label="New Password" :error="errors.newPassword" required>
        <InputText
          v-model="newPassword"
          type="password"
          placeholder="New Password"
        />
      </FormField>
      <FormField
        label="Confirm New Password"
        :error="errors.newPasswordConfirmation"
        required
      >
        <InputText
          v-model="newPasswordConfirmation"
          type="password"
          placeholder="Confirm New Password"
        />
      </FormField>
      <FormSubmitButton />
    </FormCard>
  </form>
</template>

<script setup>
import * as yup from "yup";

const props = defineProps({
  submitError: {
    type: [Error, null],
    default: null,
  },
});
const emit = defineEmits(["submit"]);

const { errors, useFieldModel, handleSubmit } = useForm({
  validationSchema: computed(() =>
    yup.object({
      currentPassword: yup.string().required().label("Current Password"),
      newPassword: yup
        .string()
        .required("New Password is required")
        .min(8, "Password must be at least 8 characters")
        .max(50, "Password must be at most 50 characters")
        .matches(/[A-Z]/, "Password must include at least one uppercase letter")
        .matches(/[a-z]/, "Password must include at least one lowercase letter")
        .matches(/\d/, "Password must include at least one number")
        .label("New Password"),
      newPasswordConfirmation: yup
        .string()
        .oneOf([yup.ref("newPassword"), null], "Passwords must match")
        .required()
        .label("Confirm New Password"),
    }),
  ),
});

const currentPassword = useFieldModel("currentPassword");
const newPassword = useFieldModel("newPassword");
const newPasswordConfirmation = useFieldModel("newPasswordConfirmation");

const errorMsgs = ref([]);
watch(
  () => props.submitError,
  () => {
    if (!props.submitError) return;

    const error = props.submitError.response.data.error;
    const errorMessage = error
      ? `Status Code ${error.code} ${error.message}`
      : `Status Code ${props.submitError.status} ${props.submitError.statusText}`;
    errorMsgs.value.push(errorMessage);
  },
);

const onSubmit = handleSubmit((values) => emit("submit", values));
</script>
