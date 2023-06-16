<template>
  <Toast />
  <PageContainer>
    <NavigationContainer>
      <GoBackButton />
    </NavigationContainer>
    <AuctionConfigurationForm :value="resource" @submit="handleSubmit" />
  </PageContainer>
</template>

<script setup>
const resource = {};
const resourcesPath = "/auction_configurations";
const handleSubmit = useFormSubmit(resourcesPath, "post", "Auction configuration created!", async (id) =>
  navigateTo(`${resourcesPath}/${id}`)
);
</script>
