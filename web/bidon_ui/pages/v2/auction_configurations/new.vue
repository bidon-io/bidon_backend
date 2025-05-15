<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton :path="resourcesPath" />
    </NavigationContainer>
    <AuctionConfigurationsForm :value="resource" @submit="handleSubmit" />
  </PageContainer>
</template>

<script setup>
const resource = {};
const resourcesPath = "/v2/auction_configurations";
const handleSubmit = useCreateResource({
  path: resourcesPath,
  message: "Auction configuration created!",
});
</script>
