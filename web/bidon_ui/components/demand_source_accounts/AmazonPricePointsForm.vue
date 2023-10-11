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
  </div>
</template>

<script setup>
import { useField } from "vee-validate";

const { value: pricePoints } = useField("extra.pricePointsMap");
const showUploads = ref(!pricePoints.value);

const bannerPricePoints = ref([]);
const mrecPricePoints = ref([]);
const interstitialPricePoints = ref([]);
const videoPricePoints = ref([]);

const isValid = computed(() =>
  [
    bannerPricePoints.value.length,
    mrecPricePoints.value.length,
    interstitialPricePoints.value.length,
    videoPricePoints.value.length,
  ].every((length) => length > 0),
);

watchEffect(() => {
  if (!isValid.value) return;

  pricePoints.value = Object.fromEntries(
    [
      ...bannerPricePoints.value,
      ...mrecPricePoints.value,
      ...interstitialPricePoints.value,
      ...videoPricePoints.value,
    ]
      .filter((el) => el.pricePoint)
      .map((el) => [el.pricePoint, el]),
  );
});
</script>
