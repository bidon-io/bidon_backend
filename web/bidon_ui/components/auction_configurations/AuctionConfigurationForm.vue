<template>
  <form @submit="onSubmit">
    <FormCard title="Auction Configuration">
      <FormField label="Name" :error="errors.name" required>
        <InputText v-model="name" type="text" placeholder="Name" />
      </FormField>
      <AppDropdown v-model="appId" :error="errors.appId" required />
      <AdTypeDropdown v-model="adType" :error="errors.adType" required />
      <FormField label="Price floor" :error="errors.pricefloor" required>
        <InputNumber
          v-model="pricefloor"
          input-id="pricefloor"
          :min-fraction-digits="2"
          :max-fraction-digits="5"
          placeholder="Price floor"
        />
      </FormField>
      <FormField label="Rounds" :error="errors.rounds" required>
        <TextareaJSON v-model="rounds" rows="5" />
      </FormField>
      <SegmentDropdown v-model="segmentId" :error="errors.segmentId" />
      <FormField
        label="External Win Notification"
        :error="errors.externalWinNotifications"
      >
        <Checkbox v-model="externalWinNotifications" :binary="true" />
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
    name: yup.string().required().label("Name"),
    appId: yup.number().required().label("App Id"),
    adType: yup.string().required().label("AdType"),
    pricefloor: yup.number().positive().required().label("Pricefloor"),
    rounds: yup.array().required(),
    appId: yup.number().required().label("Segment Id"),
    externalWinNotifications: yup.boolean(),
  }),
  initialValues: {
    name: resource.value.name || "",
    appId: resource.value.appId || null,
    adType: resource.value.adType || "",
    pricefloor: resource.value.pricefloor || null,
    rounds: resource.value.rounds || [],
    segmentId: resource.value.segmentId || null,
    externalWinNotifications: resource.value.externalWinNotifications || false,
  },
});

const name = useFieldModel("name");
const appId = useFieldModel("appId");
const adType = useFieldModel("adType");
const pricefloor = useFieldModel("pricefloor");
const rounds = useFieldModel("rounds");
const segmentId = useFieldModel("segmentId");
const externalWinNotifications = useFieldModel("externalWinNotifications");

const onSubmit = handleSubmit((values) => emit("submit", values));
</script>
