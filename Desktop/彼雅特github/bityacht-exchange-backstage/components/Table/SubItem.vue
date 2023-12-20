<script setup>
import { twMerge } from "tailwind-merge";

const { class: className, alignTop } = defineProps({
  class: {
    type: String,
    required: false,
    default: "",
  },
  alignTop: {
    type: Boolean,
    required: false,
    default: false,
  },
});
const itemClass = computed(() => {
  return twMerge(
    "grid grid-cols-[250px_1fr]",
    "[&>*]:py-3 [&>*]:px-4 [&>*]:border-l [&>*]:border-b [&>*]:border-[#DDE1E5]",
    "[&>*]:flex [&>*]:gap-2",
    "[&>*:first-child]:bg-[#F1F3F6]",
    alignTop ? "[&>*]:items-start" : "[&>*]:items-center",
    className
  );
});
</script>
<template>
  <div :class="itemClass">
    <slot />
  </div>
</template>
