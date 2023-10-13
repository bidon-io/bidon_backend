<template>
  <form @submit="onSubmit">
    <Card class="p-6 max-w-lg mx-auto">
      <template #header>
        <div class="flex justify-center">
          <h1 class="text-2xl font-bold">Log In</h1>
        </div>
      </template>
      <template #content>
        <div class="flex flex-col mb-4">
          <label class="font-semibold text-gray-500 mb-2" for="emailInput"
            >Email</label
          >
          <InputText
            id="emailInput"
            v-model="email"
            type="text"
            placeholder="Email"
          />
          <small v-if="errors.email" class="p-error">{{ errors.email }}</small>
        </div>
        <div class="flex flex-col mb-4">
          <label class="font-semibold text-gray-500 mb-2" for="passwordInput"
            >Password</label
          >
          <InputText
            id="passwordInput"
            v-model="password"
            type="password"
            placeholder="Pasword"
          />
          <small v-if="errors.password" class="p-error">{{
            errors.password
          }}</small>
        </div>
        <Button
          type="submit"
          label="Log In"
          class="p-button-primary w-full block"
        />
      </template>
    </Card>
  </form>
</template>

<script setup>
import * as yup from "yup";

definePageMeta({
  layout: "auth",
});

const { errors, useFieldModel, handleSubmit } = useForm({
  validationSchema: yup.object({
    email: yup.string().required().label("Email"),
    password: yup.string().required().label("Password"),
  }),
});

const email = useFieldModel("email");
const password = useFieldModel("password");

const authStore = useAuthStore();

const onSubmit = handleSubmit(
  async ({ email, password }) => await authStore.login(email, password),
);
</script>
