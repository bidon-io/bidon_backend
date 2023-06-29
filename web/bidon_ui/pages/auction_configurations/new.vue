<template>
  <PageContainer>
    <NavigationContainer>
      <GoBackButton />
    </NavigationContainer>
    <AuctionConfigurationForm :value="resource" @submit="handleSubmit" />
  </PageContainer>
</template>

<script setup>
const resource = { rounds: [] };
const resourcesPath = "/auction_configurations";
const handleSubmit = useCreateResource({ path: resourcesPath, message: "Auction configuration created!" });
</script>
