<template>
  <transition-group name="p-message" tag="div">
    <Message v-for="(msg, index) in errorMsgs" :key="index" severity="error">{{
      msg
    }}</Message>
  </transition-group>
  <form @submit="onSubmit">
    <FormCard title="Line Item">
      <AppDropdown v-model="appId" :error="errors.appId" required />
      <AdTypeWithFormatDropdown
        v-model="adTypeWithFormat"
        :error="errors.adType"
        required
      />
      <DemandSourceTypeDropdown
        v-model="accountType"
        label="Demand Source"
        :error="errors.accountType"
        required
      />
      <DemandSourceAccountDropdown
        v-model="accountId"
        :error="errors.accountId"
        :accounts="demandSourceAccounts"
        :disabled="!accountType"
        required
      />
      <FormField label="Label" :error="errors.humanName" required>
        <InputText v-model="humanName" type="text" placeholder="Label" />
      </FormField>
      <FormField label="Auction Type">
        <div class="flex flex-wrap gap-3 my-2">
          <div class="flex align-items-center">
            <RadioButton
              v-model="auctionType"
              input-id="biddingAuction"
              name="auctionType"
              value="bidding"
            />
            <label for="biddingAuction" class="ml-2">Bidding</label>
          </div>
          <div class="flex align-items-center">
            <RadioButton
              v-model="auctionType"
              input-id="defaultAuction"
              name="auctionType"
              value="default"
            />
            <label for="defaultAuction" class="ml-2">Default</label>
          </div>
        </div>
      </FormField>
      <FormField
        v-if="auctionType === 'default'"
        label="Bid Floor"
        :error="errors.bidFloor"
        required
      >
        <InputNumber
          v-model="bidFloor"
          input-id="bidFloor"
          :min-fraction-digits="2"
          :max-fraction-digits="5"
          placeholder="Bid Floor"
        />
      </FormField>
      <LineItemExtraFormFields v-model:schema="extraSchema" :api-key="apiKey" />
      <FormSubmitButton :disabled="!meta.valid" />
    </FormCard>
  </form>
</template>

<script setup>
import axios from "@/services/ApiService";
import * as yup from "yup";

const props = defineProps({
  value: {
    type: Object,
    required: true,
  },
  submitError: {
    type: [Error, null],
    default: null,
  },
});
const emit = defineEmits(["submit"]);
const resource = ref(props.value);

const extraSchema = ref(yup.object());
const auctionType = ref(resource.value.isBidding ? "bidding" : "default");

const adTypeWithFormat = ref({
  adType: resource.value.adType,
  format: resource.value.format,
});
const { errors, meta, useFieldModel, handleSubmit } = useForm({
  validationSchema: computed(() =>
    yup.object({
      humanName: yup.string().required().label("Label"),
      appId: yup.number().required().label("App Id"),
      bidFloor:
        auctionType.value !== "bidding"
          ? yup.number().positive().required().label("Bid Floor")
          : yup.number().nullable(true).positive().label("Bid Floor"),
      adType: yup.string().required().label("AdType"),
      format: yup
        .string()
        .nullable(true)
        .when("adType", {
          is: "banner",
          then: (schema) =>
            schema.required("Format is required for Banner Ad Type"),
        }),
      accountId: yup.number().required().label("Account Id"),
      accountType: yup.string().required().label("Demand Source"),
      isBidding: yup.boolean(),
      extra: extraSchema.value,
    })
  ),
  initialValues: {
    humanName: resource.value.humanName || "",
    appId: resource.value.appId || null,
    bidFloor: resource.value.bidFloor
      ? parseFloat(resource.value.bidFloor)
      : null,
    adType: resource.value.adType || "",
    format: resource.value.format || null,
    accountId: resource.value.accountId || null,
    accountType: resource.value.accountType || "",
    isBidding: resource.value.isBidding || false,
    extra: resource.value.extra || {},
  },
});

const humanName = useFieldModel("humanName");
const appId = useFieldModel("appId");
const bidFloor = useFieldModel("bidFloor");
const adType = useFieldModel("adType");
const format = useFieldModel("format");
const accountId = useFieldModel("accountId");
const accountType = useFieldModel("accountType");
const isBidding = useFieldModel("isBidding");

// filter demand source accounts by account type
const response = await axios.get("/demand_source_accounts");
const demandSourceAccountsAll = response.data;
const demandSourceAccounts = computed(() =>
  demandSourceAccountsAll.filter(
    (account) => account.type === accountType.value
  )
);

// compute demand source api key from account type (e.g. "DemandSource::Admob" => "admob")
// in order to fetch extra fields schema specific to the demand source
const apiKey = computed(() =>
  accountType.value ? accountType.value.split("::")[1].toLowerCase() : ""
);

// reset accountId when accountType changes
watch(accountType, () => (accountId.value = null));

// track adType and format together
watchEffect(() => {
  if (!adTypeWithFormat.value.adType) return;

  adType.value = adTypeWithFormat.value.adType;
  format.value = adTypeWithFormat.value.format;
});

watchEffect(() => (isBidding.value = auctionType.value === "bidding"));

// push submit error to error messages
const errorMsgs = ref([]);
watch(
  () => props.submitError,
  () => {
    if (!props.submitError) return;

    const error = props.submitError.response.data.error;
    const errorMessage = error
      ? `Status Code ${error.code} ${error.message}`
      : `Status Code ${props.submitError.status} ${props.submitError.statusText}`;
    errorMsgs.value.push(errorMessage);
  }
);

const onSubmit = handleSubmit((values) => emit("submit", values));
</script>
