<template>
  <div class="flex flex-col p-4 border rounded shadow-sm">
    <FileUpload
      mode="basic"
      :name="label"
      accept=".csv"
      :auto="true"
      :choose-label="label"
      @upload="onUpload"
    />
    <ScrollPanel
      v-if="content.length"
      style="width: 100%; height: 200px"
      class="mt-4"
    >
      <table>
        <thead>
          <tr>
            <th v-for="header in headers" :key="header">{{ header }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(row, rowIndex) in content" :key="rowIndex">
            <td v-for="cell in row" :key="cell">{{ cell }}</td>
          </tr>
        </tbody>
      </table>
    </ScrollPanel>
  </div>
</template>

<script setup>
defineProps({
  label: {
    type: String,
    required: true,
  },
  valueModel: {
    type: Array,
    default: () => [],
  },
});
const emit = defineEmits(["update:modelValue"]);
const value = computed({
  get() {
    return props.valueModel;
  },
  set(value) {
    emit("update:modelValue", value);
  },
});

const headers = ref([]);
const content = ref([]);

const contentToJson = (content) =>
  content.map((row) => ({
    name: row[0],
    price: parseFloat(row[1]),
    pricePoint: row[2],
  }));

const onUpload = (event) => {
  const file = event.files[0];
  const reader = new FileReader();

  reader.onload = function (e) {
    const data = e.target.result;
    const rows = data.split("\n");
    const [headersValue, ...contentValue] = rows.map((row) => row.split(","));

    headers.value = headersValue;
    content.value = contentValue;

    value.value = contentToJson(contentValue);
  };
  reader.readAsText(file);
};
</script>
