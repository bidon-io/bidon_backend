<template>
  <FormField
    label="Slot Uuid"
    :error="slotUuidErrorMessage"
    :required="required"
  >
    <InputText v-model="slotUuid" placeholder="Slot Uuid" />
  </FormField>
  <FormField
    v-if="adTypeWithFormat.adType === 'interstitial'"
    label="Is Video"
    :required="required"
  >
    <Checkbox v-model="isVideo" binary />
  </FormField>
</template>

<script setup>
import { useField } from "vee-validate";

const props = defineProps({
  adTypeWithFormat: {
    type: [Object, null],
    required: true,
  },
  required: {
    type: Boolean,
    default: false,
  },
});

const { value: slotUuid, errorMessage: slotUuidErrorMessage } =
  useField("extra.slotUuid");
const { value: format } = useField("extra.format");
const isVideo = ref(format.value === "VIDEO");

const selectFormat = ({ adType, format }, isVideo) => {
  switch (adType) {
    case "banner":
      return format === "MREC" ? "MREC" : "BANNER";
    case "interstitial":
      return isVideo ? "VIDEO" : "INTERSTITIAL";
    case "rewarded":
      return "REWARDED";
  }
};

watchEffect(
  () => (format.value = selectFormat(props.adTypeWithFormat, isVideo.value)),
);
</script>
