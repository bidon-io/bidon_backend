<template>
  <div class="flex flex-row py-2">
    <div class="w-1/4 px-6">
      <div class="font-semibold text-gray-500">{{ props.label }}</div>
    </div>
    <div class="px-6">
      <slot></slot>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  label: {
    type: String,
    required: true,
  },
});
</script>
