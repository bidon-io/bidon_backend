<template>
  <form @submit="onSubmit">
    <FormCard title="Demand source account">
      <UserDropdown v-model="userId" :error="errors.userId" required />
      <DemandSourceTypeDropdown v-model="type" :error="errors.type" required />
      <DemandSourceDropdown v-model="demandSourceId" :error="errors.demandSourceId" required />
      <FormField label="Bidding">
        <Checkbox v-model="isBidding" :binary="true" />
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
    userId: yup.number().required().label("User Id"),
    type: yup.string().required().label("Demand Source Type"),
    demandSourceId: yup.number().required().label("Deamand Source Id"),
    isBidding: yup.boolean(),
    extra: yup.object(),
  }),
  initialValues: {
    userId: resource.value.userId || null,
    type: resource.value.type || "",
    demandSourceId: resource.value.demandSourceId || null,
    isBidding: resource.value.isBidding || false,
    extra: resource.value.extra || {},
  },
});

const userId = useFieldModel("userId");
const type = useFieldModel("type");
const demandSourceId = useFieldModel("demandSourceId");
const isBidding = useFieldModel("isBidding");
const extra = useFieldModel("extra");

const onSubmit = handleSubmit((values) => emit("submit", values));
</script>
