<script setup>
import { twMerge } from 'tailwind-merge';

const { variant, loading } = defineProps({
  variant: {
    type: String,
    default: 'primary',
  },
  loading: {
    type: Boolean,
    default: false,
  },
});
const color = computed(() => {
  switch (variant) {
    case 'waterBlueBorder':
      return 'border border-waterBlue text-waterBlue';
    case 'waterBlue':
      return 'bg-waterBlue text-white';
    case 'darkBlue':
      return 'bg-darkBlue text-white';
    case 'white':
      return 'bg-white text-gray_800 border border-gray_800';
    case 'whiteBlue':
      return 'bg-white text-waterBlue border border-waterBlue';
    case 'gray':
      return 'bg-gray_500 text-white';
    case 'shallow_gray':
      return 'bg-shallow_gray text-white';
    case 'outline':
      return 'bg-transparent border border-gray_500 text-[#394B6A] hover:text-opacity-80 active:text-opacity-90';
    case 'fullWaterBlueBorder':
      return 'w-full border border-waterBlue text-waterBlue';
    case 'disabled':
      return 'bg-gray_500 text-white cursor-not-allowed';
  }
});
</script>
<template>
  <button :class="twMerge(`text-normal 2xl:px-7 px-5 md:py-1.5 py-1 rounded-3xl hover:bg-opacity-80 active:bg-opacity-90`, color)" :disabled="loading">
    <loading-outlined v-if="loading" class="mr-2" />
    <slot />
  </button>
</template>
