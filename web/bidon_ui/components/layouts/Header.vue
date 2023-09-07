<template>
  <header class="z-5 py-4 bg-white shadow-md">
    <div class="container mx-auto px-6 lg:px-8">
      <div class="flex items-center justify-between">
        <div class="flex items-center">
          <div class="relative mx-4 lg:mx-0">
            <span class="left-0 p-3 flex items-center">
              <a class="text-blue-500 font-semibold" target="_blank" href="/"
                >BidON Management Console</a
              >
            </span>
          </div>
        </div>
        <Button
          v-if="user"
          label="Log out"
          icon="pi pi-sign-out"
          text
          @click="logout"
        />
      </div>
    </div>
  </header>
</template>

<script setup>
import { useAuthStore } from "@/stores/AuthStore.js";

const { user, logout } = useAuthStore();
</script>
