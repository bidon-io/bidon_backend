<template>
  <div class="flex h-screen bg-gray-100">
    <div class="bg-white w-64 md:w-1/4 lg:w-1/5 h-full overflow-y-auto">
      <div class="flex items-center justify-center mt-8">
        <div class="flex items-center">
          <!-- <img class="h-6 w-6" src="favicon.ico" alt="logo" /> -->
          <span class="text-blue-500 text-2xl mx-2 font-semibold">BidON</span>
        </div>
      </div>

      <nav class="mt-10">
        <nuxt-link
          v-for="item in menuItems"
          :key="item.name"
          :to="item.link"
          class="flex items-center mt-4 py-2 px-6 text-gray-600 hover:bg-gray-700 hover:bg-opacity-25 hover:text-gray-100"
        >
          <span class="mx-3">{{ item.name }}</span>
        </nuxt-link>
      </nav>
    </div>

    <div class="flex flex-col flex-1">
      <Header></Header>
      <main class="flex-1 overflow-x-hidden overflow-y-auto bg-gray-200">
        <div class="container mx-auto px-8 py-10">
          <slot></slot>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup>
const menuItems = [
  { name: "Apps", link: "resources/apps/" },
  { name: "AuctionConfigurations", link: "/auction_configurations/" },
];
</script>
