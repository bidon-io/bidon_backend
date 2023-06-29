<template>
  <form @submit="onSubmit">
    <FormCard title="App">
      <PlatformIdDropdown v-model="platformId" :error="errors.platformId" required />
      <FormField label="Human Name" required>
        <InputText v-model="humanName" :error="errors.humanId" reuired type="text" placeholder="Name" />
      </FormField>
      <FormField label="Package Name" :error="errors.packageName" required>
        <InputText v-model="packageName" type="text" placeholder="Name" />
      </FormField>
      <UserDropdown v-model="userId" :error="errors.userId" required />
      <FormField label="App Key" :error="errors.appKey" required>
        <InputText v-model="appKey" type="text" placeholder="Name" />
      </FormField>
      <FormField label="Settings">
        <TextareaJSON v-model="settings" :error="errors.settings" rows="5" />
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
    platformId: yup.string().required().label("Platform Id"),
    humanName: yup.string().required().label("Human Name"),
    packageName: yup.string().required().label("Package Name"),
    userId: yup.number().required().label("User Id"),
    appKey: yup.string().required().label("App Key"),
    settings: yup.object(),
  }),
  initialValues: {
    platformId: resource.value.platformId || "",
    humanName: resource.value.humanName || "",
    packageName: resource.value.packageName || "",
    userId: resource.value.userId || null,
    appKey: resource.value.appKey || "",
    settings: resource.value.settings || {},
  },
});

const platformId = useFieldModel("platformId");
const humanName = useFieldModel("humanName");
const packageName = useFieldModel("packageName");
const userId = useFieldModel("userId");
const appKey = useFieldModel("appKey");
const settings = useFieldModel("settings");

const onSubmit = handleSubmit((values) => emit("submit", values));
</script>
