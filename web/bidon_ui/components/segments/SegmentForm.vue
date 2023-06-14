<template>
  <form @submit.prevent="emit('submit', resource)">
    <FormCard title="Segment">
      <FormField lable="Name">
        <InputText v-model="resource.name" type="text" placeholder="Name" />
      </FormField>
      <FormField lable="Description">
        <Textarea v-model="resource.description" rows="5" cols="50" />
      </FormField>
      <FormField lable="Filters">
        <InputJSON v-model="resource.filters" placeholder="Filters" style="min-width: 400px" />
      </FormField>
      <FormField lable="Enabled">
        <Checkbox v-model="resource.enabled" :binary="true" />
      </FormField>
      <AppDropdown v-model="resource.app_id" />
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
