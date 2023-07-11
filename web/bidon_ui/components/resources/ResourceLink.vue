<template>
  <NuxtLink v-if="isValid" :to="`${link.basePath}/${id}`" class="text-blue-500">
    {{ linkText }}
  </NuxtLink>
</template>

<script setup lang="ts">
/* eslint-disable @typescript-eslint/no-explicit-any */
const props = defineProps<{
  link: ResourceLink;
  data: any;
}>();
const { isValid, linkText, id } = props.link.extractLinkData(props.data);
/* eslint-enable @typescript-eslint/no-explicit-any */
</script>
