<template>
  <Button
    label="Go back"
    icon="pi pi-arrow-left"
    severity="secondary"
    text
    @click.prevent="goBack"
  />
</template>

<script setup>
const props = defineProps({
  path: {
    type: String,
    default: "/",
  },
});

const router = useRouter();
const goBack = () => {
  router.push(props.path);
};
</script>
