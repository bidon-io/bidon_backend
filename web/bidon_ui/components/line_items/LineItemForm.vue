<template>
  <form @submit="onSubmit">
    <FormCard title="Line Item">
      <FormField label="Human Name" :error="errors.humanName" required>
        <InputText v-model="humanName" type="text" placeholder="Human Name" />
      </FormField>
      <AppDropdown v-model="appId" :error="errors.appId" required />
      <FormField label="Bid Floor" :error="errors.bidFloor" required>
        <InputNumber
          v-model="bidFloor"
          input-id="bidFloor"
          :min-fraction-digits="2"
          :max-fraction-digits="5"
          placeholder="Bid Floor"
        />
      </FormField>
      <AdTypeDropdown v-model="adType" :error="errors.adType" required />
      <AdFormatDropdown v-model="format" :error="errors.format" required />
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
      <FormField label="Code" :error="errors.code" required>
        <InputText v-model="code" type="text" placeholder="Code" />
      </FormField>
      <FormField label="Extra">
        <TextareaJSON v-model="extra" rows="5" />
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
    humanName: yup.string().required().label("Human Name"),
    appId: yup.number().required().label("App Id"),
    bidFloor: yup.number().positive().required().label("Bid Floor"),
    adType: yup.string().required().label("AdType"),
    format: yup
      .string()
      .nullable(true)
      .when("adType", {
        is: "banner",
        then: (schema) =>
          schema.required("Format is required for Banner Ad Type"),
      }),
    accountId: yup.number().required().label("Account Id"),
    accountType: yup.string().required().label("Demand Source Type"),
    code: yup.string().required().label("Code"),
    extra: yup.object(),
  }),
  initialValues: {
    humanName: resource.value.humanName || "",
    appId: resource.value.appId || null,
    bidFloor: resource.value.bidFloor || null,
    adType: resource.value.adType || "",
    adFormat: resource.value.format || "",
    accountId: resource.value.accountId || null,
    accountType: resource.value.accountType || "",
    code: resource.value.code || "",
    extra: resource.value.extra || {},
  },
});

const humanName = useFieldModel("humanName");
const appId = useFieldModel("appId");
const bidFloor = useFieldModel("bidFloor");
const adType = useFieldModel("adType");
const format = useFieldModel("format");
const accountId = useFieldModel("accountId");
const accountType = useFieldModel("accountType");
const code = useFieldModel("code");
const extra = useFieldModel("extra");

const onSubmit = handleSubmit((values) => emit("submit", values));
</script>
