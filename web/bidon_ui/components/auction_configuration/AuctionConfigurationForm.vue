<template>
  <form @submit.prevent="handleSubmit">
    <FormCard title="Auction Configuration">
      <FormField lable="Name">
        <InputText v-model="resource.name" type="text" placeholder="Name" />
      </FormField>
      <FormField lable="App">
        <Dropdown
          v-model="resource.app_id"
          :options="apps"
          option-label="package_name"
          option-value="id"
          class="w-full md:w-14rem"
          placeholder="Select App"
        />
      </FormField>
      <FormField lable="AD Type">
        <Dropdown
          v-model="resource.ad_type"
          :options="adTypes"
          class="w-full md:w-14rem"
          placeholder="Select Ad Type"
        />
      </FormField>
      <FormField lable="Price floor">
        <InputNumber
          v-model="resource.pricefloor"
          input-id="pricefloor"
          :min-fraction-digits="2"
          :max-fraction-digits="5"
          placeholder="Price floor"
        />
      </FormField>
      <FormField lable="Rounds">
        <Textarea v-model="rounds" rows="10" cols="80" />
      </FormField>
      <FormSubmitButton />
    </FormCard>
  </form>
</template>

<script setup>
import { defineProps, defineEmits } from "vue";
import axios from "@/services/ApiService.js";

const props = defineProps({
  value: {
    type: Object,
    required: true,
  },
});
const emit = defineEmits(["submit"]);

const resource = ref(props.value);

const apps = ref([]);
axios
  .get("/apps")
  .then((response) => {
    apps.value = response.data;
  })
  .catch((error) => {
    console.error(error);
  });

const adTypes = ref(["banner", "interstitial", "rewarded"]);

const rounds = computed({
  get: () => JSON.stringify(resource.value.rounds, null, 2),
  set: (newValue) => {
    try {
      resource.value.rounds = JSON.parse(newValue);
    } catch {}
  },
});

function handleSubmit() {
  emit("submit", resource.value);
}
</script>
