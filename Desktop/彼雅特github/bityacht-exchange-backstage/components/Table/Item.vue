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
    `grid gap-2 2xl:gap-4 grid-cols-[var(--cols)] border border-[#DDE1E5] -mt-[1px] bg-white py-2.5 px-4 hover:bg-[#F9FAFB] items-center [&>*]:truncate [.table-header+&]:border-t-0 table-item`,
    className
  );
});
</script>
<template>
  <div :class="itemClass">
    <slot />
  </div>
</template>
