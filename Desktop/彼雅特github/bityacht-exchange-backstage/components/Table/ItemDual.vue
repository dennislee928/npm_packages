<script setup>
import { twMerge } from "tailwind-merge";

const { class: className } = defineProps({
  class: {
    type: String,
    required: false,
    default: "",
  },
});
const itemClass = computed(() => {
  return twMerge(
    "border-r border-[#DDE1E5] grid grid-cols-[repeat(2,1fr)] bg-white",
    className
  );
});
</script>
<template>
  <div :class="itemClass">
    <slot />
  </div>
</template>
