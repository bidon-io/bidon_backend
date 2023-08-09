<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton :path="resourcesPath" />
    </NavigationContainer>
    <AuctionConfigurationForm :value="resource" @submit="handleSubmit" />
  </PageContainer>
</template>

<script setup>
const resource = { rounds: [] };
const resourcesPath = "/auction_configurations";
const handleSubmit = useCreateResource({
  path: resourcesPath,
  message: "Auction configuration created!",
});
</script>
