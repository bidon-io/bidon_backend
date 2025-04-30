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
const route = useRoute();
const nuxtApp = useNuxtApp();

const goBack = () => {
  const previousRoute = nuxtApp.$routeHistory().getPrevious();

  if (
    previousRoute &&
    previousRoute.path !== "" &&
    previousRoute.path !== route.path
  ) {
    router.push(previousRoute);
  } else {
    router.push(props.path);
  }
};
</script>
