<script setup>
import { twMerge } from "tailwind-merge";
const { price, unit, showAllNumbers } = defineProps({
  price: {
    required: true,
  },
  unit: {
    type: String,
    required: false,
  },
  alignRight: {
    type: Boolean,
    default: true,
  },
  showAllNumbers: {
    type: Boolean,
    default: false,
  },
});
const parsedPrice = computed(() =>
  showAllNumbers
    ? parseFloat(price).toLocaleString(undefined, {})
    : parseFloat(price).toLocaleString(undefined, {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
      })
);
</script>
<template>
  <a-tooltip :mouseEnterDelay="0" :mouseLeaveDelay="0" color="white">
    <template #title>
      <span class="text-gray-950">
        {{ price }}
      </span>
    </template>
    <div
      :class="
        twMerge(`tracking-tighter`, alignRight ? `text-right` : `text-left`)
      "
    >
      <span>{{ parsedPrice }}</span>
      <span class="text-xs opacity-60 ml-[2px] mr-[0.05em]">{{ unit }}</span>
    </div>
  </a-tooltip>
</template>
