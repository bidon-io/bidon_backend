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
        <div v-if="currentUser" class="flex items-center">
          <div>
            <span>{{ currentUser.email }}</span>
            <span class="text-red-500 font-bold pl-1">
              {{ currentUser.isAdmin ? "[Admin]" : "" }}
            </span>
          </div>

          <Button label="Log out" icon="pi pi-sign-out" text @click="logout" />
        </div>
      </div>
    </div>
  </header>
</template>

<script setup>
const { currentUser, logout } = useAuthStore();
</script>
