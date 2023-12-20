<script setup>
const { cols } = defineProps({
  cols: {
    type: String,
    default: "250px 1fr",
  },
});
</script>
<template>
  <div :style="`--cols: ${cols}`" class="m-6 overflow-x-auto">
    <div class="min-w-[1000px]">
      <slot />
    </div>
  </div>
</template>
