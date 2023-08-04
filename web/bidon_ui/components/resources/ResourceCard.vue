<template>
  <FormCard :title="title">
    <FormField v-for="field in fields" :key="field.key" :label="field.label">
      <div v-if="!field.type" class="text-gray-900">
        {{ localResource[field.key] }}
      </div>
      <ResourceLink
        v-if="field.type === 'link'"
        :link="field.link"
        :data="localResource"
      />
      <AssociatedResourcesLink
        v-if="field.type === 'associatedResourcesLink'"
        :link="field.associatedResourcesLink"
        :data="localResource"
      />
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
