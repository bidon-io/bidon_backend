<template>
  <form @submit.prevent="emit('submit', resource)">
    <FormCard title="Auction Configuration">
      <FormField label="Name">
        <InputText v-model="resource.name" type="text" placeholder="Name" />
      </FormField>
      <AppDropdown v-model="resource.app_id" />
      <AdTypeDropdown v-model="resource.ad_type" />
      <FormField label="Price floor">
        <InputNumber
          v-model="resource.pricefloor"
          input-id="pricefloor"
          :min-fraction-digits="2"
          :max-fraction-digits="5"
          placeholder="Price floor"
        />
      </FormField>
      <FormField label="Rounds">
        <TextareaJSON v-model="resource.rounds" rows="5" cols="50" />
      </FormField>
      <SegmentDropdown v-model="resource.segment_id" />
      <FormField label="External Win Notification">
        <Checkbox v-model="resource.external_win_notifications" :binary="true" />
      </FormField>
      <FormSubmitButton />
    </FormCard>
  </form>
</template>

<script setup>
const props = defineProps({
  value: {
    type: Object,
    required: true,
  },
});
const emit = defineEmits(["submit"]);
const resource = ref(props.value);
</script>
