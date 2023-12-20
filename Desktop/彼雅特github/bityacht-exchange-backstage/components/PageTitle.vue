<template>
  <div class="flex justify-between items-center m-6 mt-8">
    <h1 class="font-bold text-[#394B6A] text-3xl">
      <slot />
    </h1>
    <div>
      <slot name="action" />
    </div>
  </div>
</template>
