<template>
  <NuxtLink :to="path" class="text-blue-500">
    {{ label }}
  </NuxtLink>
</template>

<script setup lang="ts">
const props = defineProps<{
  link: AssociatedResourcesLink;
  data: object;
}>();
const { label, path } = props.link.extractLinkData(props.data);
</script>
