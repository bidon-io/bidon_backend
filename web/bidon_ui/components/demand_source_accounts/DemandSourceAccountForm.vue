<template>
  <transition-group name="p-message" tag="div">
    <Message v-for="(msg, index) in errorMsgs" :key="index" severity="error">{{
      msg
    }}</Message>
  </transition-group>
  <form @submit="onSubmit">
    <FormCard title="Demand source account">
      <DemandSourceTypeDropdown
        v-model="type"
        label="Demand Source"
        :error="errors.type"
        required
      />
      <DemandSourceDropdown
        v-model="demandSourceId"
        :error="errors.demandSourceId"
        :show="false"
        :selected-api-key="apiKey"
        required
      />
      <OwnerWithSharedDropdown
        v-if="currentUser.isAdmin"
        v-model="userId"
        :error="errors.userId"
      />
      <FormField label="Label" :error="errors.label" required>
        <InputText v-model="label" type="text" placeholder="Label" />
      </FormField>
      <DemandSourceAccountExtraFormFields
        v-model:schema="extraSchema"
        :api-key="apiKey"
      />
      <FormSubmitButton :disabled="!meta.valid" />
    </FormCard>
  </form>
</template>

<script setup>
import * as yup from "yup";
import { useAuthStore } from "@/stores/AuthStore";
import OwnerWithSharedDropdown from "~/components/form/OwnerWithSharedDropdown.vue";

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
const extraSchema = ref(yup.object());

const { user: currentUser } = useAuthStore();

const { errors, meta, useFieldModel, handleSubmit } = useForm({
  validationSchema: computed(() => {
    const validationFields = {
      label: yup.string().required().label("Label"),
      type: yup.string().required().label("Demand Source Type"),
      demandSourceId: yup.number().required().label("Demand Source Id"),
      extra: extraSchema.value,
    };

    if (currentUser.isAdmin) {
      validationFields.userId = yup.number().nullable(true).label("User Id");
    }
    return yup.object(validationFields);
  }),
  initialValues: {
    userId: resource.value.userId || null,
    label: resource.value.label || "",
    type: resource.value.type || "",
    demandSourceId: resource.value.demandSourceId || null,
    extra: resource.value.extra || {},
  },
});

const userId = useFieldModel("userId");
const label = useFieldModel("label");
const type = useFieldModel("type");
const demandSourceId = useFieldModel("demandSourceId");

// compute demand source api key from account type (e.g. "DemandSource::Admob" => "admob")
// in order to fetch extra fields schema specific to the demand source
const apiKey = computed(() =>
  type.value ? type.value.split("::")[1].toLowerCase() : ""
);

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
  }
);

const onSubmit = handleSubmit((values) => emit("submit", values));
</script>
