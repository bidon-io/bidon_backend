<template>
  <form @submit="onSubmit">
    <FormCard title="App Demand Profile">
      <AppDropdown v-model="appId" :error="errors.appId" required />
      <DemandSourceDropdown
        v-model="demandSourceId"
        :error="errors.demandSourceId"
        required
      />
      <DemandSourceAccountDropdown
        v-model="accountId"
        :error="errors.accountId"
        required
      />
      <DemandSourceTypeDropdown
        v-model="accountType"
        :error="errors.accountType"
        required
      />
      <FormField label="Data">
        <TextareaJSON v-model="data" rows="5" />
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
    appId: yup.number().required().label("App Id"),
    demandSourceId: yup.number().required().label("Deamand Source Id"),
    accountId: yup.number().required().label("Account Id"),
    data: yup.object(),
    accountType: yup.string().required().label("Demand Source Type"),
  }),
  initialValues: {
    appId: resource.value.appId || null,
    demandSourceId: resource.value.demandSourceId || null,
    accountId: resource.value.accountId || null,
    data: resource.value.extra || {},
    accountType: resource.value.accountType || "",
  },
});

const appId = useFieldModel("appId");
const demandSourceId = useFieldModel("demandSourceId");
const accountId = useFieldModel("accountId");
const data = useFieldModel("data");
const accountType = useFieldModel("accountType");

const onSubmit = handleSubmit((values) => emit("submit", values));
</script>
