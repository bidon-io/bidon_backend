<template>
  <ResourceCard :fields="fields" :resource="resource" />
</template>

<script setup>
const resource = {
  id: 1,
  app: "SampleApp",
  name: "SampleAuction",
  adType: "rewarded",
  priceFloor: 0.1,
  rounds:
    '[{"id"=>"ROUND_REWARDED_1", "demands"=>["admob", "unityads", "applovin", "dtexchange"], "timeout"=>15000}, {"id"=>"ROUND_REWARDED_2", "demands"=>["bidmachine"], "timeout"=>15000}]',
};
const fields = [
  { label: "ID", key: "id" },
  { label: "App", key: "app", type: "link", link: "/apps/1" },
  { label: "Name", key: "name" },
  { label: "Ad type", key: "adType" },
  { label: "Price floor", key: "priceFloor" },
  { label: "Rounds", key: "rounds", type: "textarea" },
];
</script>
