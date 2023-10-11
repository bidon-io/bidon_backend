<template>
  <FormField label="Upload Price Points">
    <Checkbox v-model="showUploads" :binary="true" />
  </FormField>
  <div v-if="showUploads" class="grid grid-cols-2 gap-5 pt-4">
    <AmazonPricePointsUpload
      v-model="bannerPricePoints"
      label="Banner Price Points"
    />
    <AmazonPricePointsUpload
      v-model="mrecPricePoints"
      label="MREC Price Points"
    />
    <AmazonPricePointsUpload
      v-model="interstitialPricePoints"
      label="Interstitial Price Points"
    />
    <AmazonPricePointsUpload
      v-model="videoPricePoints"
      label="Video Price Points"
    />
    <AmazonPricePointsUpload
      v-model="rewardedPricePoints"
      label="Rewarded Price Points"
    />
  </div>
</template>

<script setup>
import { useField } from "vee-validate";

const { value: pricePoints } = useField("extra.pricePointsMap");
const showUploads = ref(!pricePoints.value);
const persistedPricePoints = pricePoints.value;

const bannerPricePoints = ref([]);
const mrecPricePoints = ref([]);
const interstitialPricePoints = ref([]);
const videoPricePoints = ref([]);
const rewardedPricePoints = ref([]);

const isValid = computed(() =>
  [
    bannerPricePoints.value.length,
    mrecPricePoints.value.length,
    interstitialPricePoints.value.length,
    videoPricePoints.value.length,
    rewardedPricePoints.value.length,
  ].every((length) => length > 0),
);

// reset pricePoints for new uploads (to make submit btn disabled), or restore persisted value
watch(showUploads, () => {
  if (showUploads.value) pricePoints.value = null;
  else pricePoints.value = persistedPricePoints;
});

watchEffect(() => {
  if (!isValid.value) return;

  pricePoints.value = Object.fromEntries(
    [
      ...bannerPricePoints.value,
      ...mrecPricePoints.value,
      ...interstitialPricePoints.value,
      ...videoPricePoints.value,
      ...rewardedPricePoints.value,
    ]
      .filter((el) => el.pricePoint)
      .map((el) => [el.pricePoint, el]),
  );
});
</script>
