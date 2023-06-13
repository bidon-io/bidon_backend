<template>
  <Card class="p-6">
    <template #title>{{ props.title }}</template>
    <template #content>
      <div class="divide-y">
        <div></div>
        <slot></slot>
      </div>
    </template>
  </Card>
</template>

<script setup>
import { defineProps } from "vue";
const props = defineProps({
  title: {
    type: String,
    required: true,
  },
});
</script>
