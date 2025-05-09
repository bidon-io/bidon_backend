<template>
  <div v-if="hasDropdown" class="relative flex justify-center">
    <Button
      icon="pi pi-ellipsis-v"
      class="p-button-rounded p-button-outlined"
      aria-haspopup="true"
      aria-controls="dropdown-menu"
      @click="toggleDropdown"
    />
    <Menu ref="menu" :model="dropdownItems" :popup="true" />
  </div>
  <NuxtLink v-else :to="path" class="text-blue-500">
    {{ label }}
  </NuxtLink>
</template>

<script setup lang="ts">
import { ref } from "vue";
import type { ComponentPublicInstance } from "vue";

interface DropdownItem {
  label: string;
  path: string;
}

interface LinkData {
  label: string;
  path: string;
  dropdown?: DropdownItem[];
}

interface ResourceData {
  id: string | number;
  [key: string]: unknown;
}

interface AssociatedResourcesLinkType {
  extractLinkData: (data: ResourceData) => LinkData;
}

const props = defineProps<{
  link: AssociatedResourcesLinkType;
  data: ResourceData;
}>();

const linkData = props.link.extractLinkData(props.data);
const { label, path, dropdown } = linkData;
const hasDropdown = !!dropdown && dropdown.length > 0;

// Interface for PrimeVue Menu component
interface MenuInstance extends ComponentPublicInstance {
  toggle: (event: Event) => void;
}

const menu = ref<MenuInstance | null>(null);

const router = useRouter();

const dropdownItems =
  hasDropdown && dropdown
    ? dropdown.map((item: DropdownItem) => ({
        label: item.label,
        command: () => {
          router.push(item.path);
        },
      }))
    : [];

const toggleDropdown = (event: Event) => {
  if (menu.value && typeof menu.value.toggle === "function") {
    menu.value.toggle(event);
  }
};
</script>
