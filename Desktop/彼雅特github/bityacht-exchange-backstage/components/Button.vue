<script setup>
import { twMerge } from "tailwind-merge";
const {
  variant,
  size,
  class: className,
  disabled,
  loading,
} = defineProps({
  variant: {
    type: String,
    default: "primary",
  },
  size: {
    type: String,
    default: "default",
  },
  disabled: {
    type: Boolean,
    default: false,
  },
  loading: {
    type: Boolean,
    default: false,
  },
  class: {
    type: String,
    default: "",
  },
});
const variantClass = computed(() => {
  switch (variant) {
    case "secondary":
      return "bg-[#0F62AE]";
    case "export":
      return "bg-[#25BBEE]";
    case "danger":
      return "bg-[#F5222D]";
    case "outline":
      return "bg-transparent border border-[#DDE1E5] text-[#394B6A] hover:text-opacity-80 active:text-opacity-90 hover:bg-[#DDE1E5] hover:bg-opacity-20 active:bg-opacity-40";
  }
});
const sizeClass = computed(() => {
  switch (size) {
    case "mini":
      return "px-2 py-0.5 text-sm";
  }
});
const buttonClass = computed(() =>
  twMerge(
    `btn px-4 py-2.5 rounded-md text-white bg-[#394B6A] hover:bg-opacity-80 active:bg-opacity-70 relative transition-all`,
    sizeClass.value,
    variantClass.value,
    className,
    disabled && `bg-[#DDE1E5] cursor-not-allowed pointer-events-none`,
    loading && `cursor-not-allowed pointer-events-none text-opacity-0`
  )
);
</script>
<template>
  <button :class="buttonClass">
    <div :class="twMerge(loading && 'opacity-0')">
      <slot />
    </div>
    <div
      :class="
        twMerge(
          `h-5 w-5 rounded-full animate-spin transition-all`,
          'border-2 border-gray-300 border-r-white',
          'absolute inset-0 m-auto',
          loading ? 'opacity-100' : 'opacity-0 hidden'
        )
      "
    />
  </button>
</template>
