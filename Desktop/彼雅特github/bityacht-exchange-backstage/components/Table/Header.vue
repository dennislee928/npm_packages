<script setup>
import { twMerge } from "tailwind-merge";
const { class: className } = defineProps({
  class: {
    type: String,
    default: "",
  },
});
const headerClass = computed(() => {
  return twMerge(
    "grid gap-2 2xl:gap-4 grid-cols-[var(--cols)] border border-[#19253A] bg-[#19253A] text-white font-bold py-2 px-4 table-header [.table-item+&]:-mt-[1px]",
    className
  );
});
</script>
<template>
  <div :class="headerClass">
    <slot />
  </div>
</template>
