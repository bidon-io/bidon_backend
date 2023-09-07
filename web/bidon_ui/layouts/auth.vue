<template>
  <div class="flex h-screen bg-gray-100">
    <div class="flex flex-col flex-1">
      <Header></Header>
      <Toast />
      <ConfirmDialog />
      <div
        class="flex-1 flex flex-col justify-between overflow-x-hidden overflow-y-auto bg-gray-200"
      >
        <main class="container mx-auto px-8 py-10">
          <slot></slot>
        </main>
      </div>
    </div>
  </div>
</template>
