<template>
  <div class="my-4 py-4">
    <Button
      type="submit"
      label="Save"
      icon="pi pi-save"
      class="p-button-success"
      :disabled="disabled"
    />
  </div>
</template>

<script setup>
defineProps({
  disabled: {
    type: Boolean,
    default: false,
  },
});
</script>
