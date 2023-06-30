<template>
  <FormCard :title="title">
    <FormField v-for="field in fields" :key="field.key" :label="field.label">
      <div v-if="!field.type" class="text-gray-900">
        {{ localResource[field.key] }}
      </div>
      <NuxtLink v-if="field.type === 'link'" :to="field.link">{{
        localResource[field.key]
      }}</NuxtLink>
      <Textarea
        v-if="field.type === 'textarea'"
        :value="JSON.stringify(localResource[field.key])"
        rows="5"
        cols="50"
        disabled
      />
    </FormField>
  </FormCard>
</template>

<script setup>
import { ref } from "vue";

const props = defineProps({
  title: {
    type: String,
    required: true,
  },
  resource: {
    type: Object,
    required: true,
  },
  fields: {
    type: Array,
    required: true,
  },
});
const localResource = ref(props.resource);
</script>
