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
      <FormField label="Store ID" :error="errors.storeId">
        <InputText
          v-model="storeId"
          type="text"
          placeholder="e.g., com.example.app or 123456789"
        />
        <small class="p-text-secondary">
          App store identifier (bundle ID for iOS, package name for Android)
        </small>
      </FormField>
      <FormField label="Store URL" :error="errors.storeUrl">
        <InputText
          v-model="storeUrl"
          type="url"
          placeholder="https://apps.apple.com/app/id123456789"
        />
        <small class="p-text-secondary">
          Direct link to the app's store page
        </small>
      </FormField>
      <FormField label="Categories" :error="errors.categories">
        <InputText
          v-model="categoriesText"
          type="text"
          placeholder="IAB1, IAB9-30, IAB14 (comma-separated)"
        />
        <small class="p-text-secondary">
          IAB content categories for better ad targeting (comma-separated)
        </small>
      </FormField>
      <FormSubmitButton />
    </FormCard>
  </form>
</template>

<script setup>
import * as yup from "yup";
import { computed } from "vue";

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
  storeId: yup.string().label("Store ID"),
  storeUrl: yup.string().label("Store URL"),
  categories: yup.array().of(yup.string()).label("Categories"),
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
    storeId: resource.value.storeId || "",
    storeUrl: resource.value.storeUrl || "",
    categories: resource.value.categories || [],
  },
});

const platformId = useFieldModel("platformId");
const humanName = useFieldModel("humanName");
const packageName = useFieldModel("packageName");
const userId = useFieldModel("userId");
const storeId = useFieldModel("storeId");
const storeUrl = useFieldModel("storeUrl");
const categories = useFieldModel("categories");

// Convert categories array to/from comma-separated string
const categoriesText = computed({
  get: () => {
    return Array.isArray(categories.value) ? categories.value.join(", ") : "";
  },
  set: (value) => {
    categories.value = value
      ? value
          .split(",")
          .map((s) => s.trim())
          .filter((s) => s)
      : [];
  },
});

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
