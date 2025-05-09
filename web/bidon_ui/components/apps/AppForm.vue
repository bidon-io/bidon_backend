<template>
  <transition-group name="p-message" tag="div">
    <Message v-for="(msg, index) in errorMsgs" :key="index" severity="error">{{
      msg
    }}</Message>
  </transition-group>
  <form @submit="onSubmit">
    <FormCard title="App">
      <OwnerDropdown
        v-if="currentUser.isAdmin"
        v-model="userId"
        :error="errors.userId"
        required
      />
      <FormField label="App Name" :error="errors.humanName" required>
        <InputText v-model="humanName" type="text" placeholder="Name" />
      </FormField>
      <PlatformIdDropdown
        v-model="platformId"
        label="Platform"
        :error="errors.platformId"
        required
      />
      <FormField label="Package Name" :error="errors.packageName" required>
        <InputText v-model="packageName" type="text" placeholder="Name" />
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
  submitError: {
    type: [Error, null],
    default: null,
  },
});
const emit = defineEmits(["submit"]);
const resource = ref(props.value);

const { currentUser } = useAuthStore();
let validationFields = {
  platformId: yup.string().required().label("Platform"),
  humanName: yup.string().required().label("Owner Name"),
  packageName: yup.string().required().label("Package Name"),
};

if (currentUser.isAdmin) {
  validationFields.userId = yup.number().required().label("User Id");
}

const validationSchema = yup.object(validationFields);

const { errors, useFieldModel, handleSubmit } = useForm({
  validationSchema,
  initialValues: {
    platformId: resource.value.platformId || "",
    humanName: resource.value.humanName || "",
    packageName: resource.value.packageName || "",
    userId: resource.value.userId || null,
  },
});

const platformId = useFieldModel("platformId");
const humanName = useFieldModel("humanName");
const packageName = useFieldModel("packageName");
const userId = useFieldModel("userId");

// push submit error to error messages
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
