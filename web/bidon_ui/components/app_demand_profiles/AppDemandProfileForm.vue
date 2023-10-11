<template>
  <transition-group name="p-message" tag="div">
    <Message v-for="(msg, index) in errorMsgs" :key="index" severity="error">{{
      msg
    }}</Message>
  </transition-group>
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
        :accounts="accounts"
        :error="errors.accountId"
        :disabled="!demandSourceId"
        required
      />
      <AppDemandProfileDataFormFields
        v-model:schema="dataSchema"
        :account-type="accountType"
      />
      <FormSubmitButton />
    </FormCard>
  </form>
</template>

<script setup>
import axios from "@/services/ApiService";
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

const dataSchema = ref(yup.object());
const { errors, useFieldModel, handleSubmit } = useForm({
  validationSchema: computed(() =>
    yup.object({
      appId: yup.number().required().label("App Id"),
      demandSourceId: yup.number().required().label("Deamand Source Id"),
      accountId: yup.number().required().label("Account Id"),
      data: dataSchema.value,
    }),
  ),
  initialValues: {
    appId: resource.value.appId || null,
    demandSourceId: resource.value.demandSourceId || null,
    accountId: resource.value.accountId || null,
    data: resource.value.data || {},
    accountType: resource.value.accountType || "",
  },
});

const appId = useFieldModel("appId");
const demandSourceId = useFieldModel("demandSourceId");
const accountId = useFieldModel("accountId");

// filter demand source accounts by demand source id
const response = await axios.get("/demand_source_accounts");
const accountsAll = response.data;
const accounts = computed(() =>
  accountsAll.filter(
    (account) => account.demandSourceId === demandSourceId.value,
  ),
);

const accountToTypeMapping = new Map(
  accountsAll.map((account) => [account.id, account.type]),
);

const accountType = computed(
  () => accountToTypeMapping.get(accountId.value) || "",
);

watch(demandSourceId, () => (accountId.value = null));

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

const onSubmit = handleSubmit((values) =>
  emit("submit", { ...values, accountType: accountType.value }),
);
</script>
