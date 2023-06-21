<template>
  <div class="flex h-screen bg-gray-100">
    <div class="bg-white w-64 md:w-1/4 lg:w-1/5 h-full overflow-y-auto">
      <Logo />
      <SidebarNavigation />
    </div>

    <div class="flex flex-col flex-1">
      <Header></Header>
      <Toast />
      <ConfirmDialog />
      <main class="flex-1 overflow-x-hidden overflow-y-auto bg-gray-200">
        <div class="container mx-auto px-8 py-10">
          <slot></slot>
        </div>
      </main>
    </div>
  </div>
</template>
