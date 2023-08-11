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
        :accounts="accounts"
        :error="errors.accountId"
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
    })
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

const response = await axios.get("/demand_source_accounts");
const accounts = ref(response.data);

const accountToTypeMapping = new Map(
  accounts.value.map((account) => [account.id, account.type])
);
const accountType = computed(
  () => accountToTypeMapping.get(accountId.value) || ""
);

const onSubmit = handleSubmit((values) =>
  emit("submit", { ...values, accountType: accountType.value })
);
</script>
