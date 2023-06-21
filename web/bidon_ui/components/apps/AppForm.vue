<template>
  <form @submit.prevent="emit('submit', resource)">
    <FormCard title="Appp">
      <PlatformIdDropdown v-model="resource.platform_id" />
      <FormField label="Human Name">
        <InputText v-model="resource.human_name" type="text" placeholder="Name" />
      </FormField>
      <FormField label="Package Name">
        <InputText v-model="resource.package_name" type="text" placeholder="Name" />
      </FormField>
      <UserDropdown v-model="resource.user_id" />
      <FormField label="App Key">
        <InputText v-model="resource.app_key" type="text" placeholder="Name" />
      </FormField>
      <FormField label="Settings">
        <TextareaJSON v-model="resource.settings" rows="5" cols="50" />
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
