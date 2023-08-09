<template>
  <form @submit="onSubmit">
    <FormCard title="Country">
      <FormField label="Human Name" :error="errors.humanName" required>
        <InputText v-model="humanName" type="text" placeholder="Name" />
      </FormField>
      <FormField label="Alpha 2 Code" :error="errors.alpha2Code" required>
        <InputText
          v-model="alpha2Code"
          type="text"
          placeholder="Alpha 2 Code"
        />
      </FormField>
      <FormField label="Alpha 3 Code" :error="errors.alpha3Code" required>
        <InputText
          v-model="alpha3Code"
          type="text"
          placeholder="Alpha 3 Code"
        />
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
    alpha2Code: yup.string().required().label("Alpha 2 Code"),
    alpha3Code: yup.string().required().label("Alpha 3 Code"),
  }),
  initialValues: {
    humanName: resource.value.humanName || "",
    alpha2Code: resource.value.alpha2Code || "",
    alpha3Code: resource.value.alpha3Code || "",
  },
});

const humanName = useFieldModel("humanName");
const alpha2Code = useFieldModel("alpha2Code");
const alpha3Code = useFieldModel("alpha3Code");

const onSubmit = handleSubmit((values) => emit("submit", values));
</script>
