<template>
  <div class="flex-1 p-6 mx-auto w-full">
    <Card>
      <template #title>Auction Config {{ id }}</template>
      <template #content>
        <div class="divide-y">
          <div v-for="field in fields" :key="field.key" class="flex flex-row py-2">
            <div class="w-1/4 px-6">
              <div class="font-semibold text-gray-500">{{ field.label }}</div>
            </div>
            <div class="px-6">
              <div v-if="!field.type" class="text-gray-900">{{ localResource[field.key] }}</div>
              <NuxtLink v-if="field.type === 'link'" :to="field.link">{{ localResource[field.key] }}</NuxtLink>
              <Textarea
                v-if="field.type === 'textarea'"
                v-model="localResource[field.key]"
                rows="5"
                cols="80"
                disabled
              />
            </div>
          </div>
        </div>
      </template>
    </Card>
  </div>
</template>

<script setup>
import { ref } from "vue";

const props = defineProps({
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
