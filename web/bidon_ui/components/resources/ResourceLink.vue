<template>
  <NuxtLink v-if="isValid" :to="`${link.basePath}/${id}`" class="text-blue-500">
    {{ linkText }}
  </NuxtLink>
</template>

<script setup lang="ts">
const props = defineProps<{
  link: ResourceLink;
  data: object;
}>();
const { isValid, linkText, id } = props.link.extractLinkData(props.data);
</script>
