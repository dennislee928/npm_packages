<script setup>
import { twMerge } from "tailwind-merge";
const { href } = defineProps({
  href: {
    type: String,
    required: true,
  },
});
const isActive = computed(() => useRoute().path === href);
</script>
<template>
  <router-link
    :to="href"
    :class="
      twMerge(
        `border-b-[3px] border-[#394B6A] p-1 text-[#394B6A] transition-all`,
        isActive
          ? `font-bold`
          : `opacity-70 border-opacity-70 hover:opacity-100 hover:border-opacity-90`
      )
    "
  >
    <slot />
  </router-link>
</template>
