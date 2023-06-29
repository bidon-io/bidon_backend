<template>
  <div class="flex flex-row py-2">
    <div class="w-1/4 px-6">
      <div class="font-semibold text-gray-500">
        <span class="font-semibold text-gray-500">{{ props.label }}</span>
        <span v-if="required" class="text-red-600 ml-1">*</span>
      </div>
    </div>
    <div class="px-6 w-3/4">
      <slot></slot>
      <small v-if="error" class="p-error">{{ error }}</small>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  error: {
    type: String,
    default: "",
  },
  label: {
    type: String,
    required: true,
  },
  required: {
    type: Boolean,
    default: false,
  },
});
</script>

<style>
input,
textarea {
  min-width: 100%;
}
</style>
