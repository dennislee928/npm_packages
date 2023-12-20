<script setup>
import { twMerge } from 'tailwind-merge';
const {
  status,
  variant,
  class: className,
} = defineProps({
  status: {
    type: String,
    required: true,
  },
  variant: {
    type: String,
    default: 'default',
  },
  class: {
    type: String,
    default: '',
  },
});
const statusColor = {
  'in-progress': 'bg-[#DDE1E5] text-gray-400',
  completed: 'bg-green-500 text-green-100',
  cancelled: 'bg-red-500 text-red-100',
  frozen: 'bg-blue-200 text-blue-800',
  warn: 'bg-yellow-500 text-yellow-100',
  checking: 'bg-[#DDE1E5] text-gray-600',
};
const filledClassName = computed(() => twMerge(`py-0.5 px-1 rounded-md text-center text-white max-w-[4em] tracking-tighter`, statusColor[status], className));
const defaultClassName = computed(() => twMerge(`rounded-full w-2.5 h-2.5`, statusColor[status], className));
</script>
<template>
  <div class="flex gap-1.5 items-center" v-if="variant === `default`">
    <div :class="defaultClassName"></div>
    <span class="flex-1 truncate">
      <slot />
    </span>
  </div>
  <div :class="filledClassName" v-if="variant === `filled`">
    <slot />
  </div>
</template>
