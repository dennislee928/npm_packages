<script setup>
import { twMerge } from "tailwind-merge";
const { href, exact } = defineProps({
  href: String,

  exact: {
    type: Boolean,
    default: false,
  },
});
const isActive = computed(() =>
  exact ? useRoute().path === href : useRoute().path.startsWith(href)
);
</script>
<template>
  <NuxtLink
    :to="href"
    :class="
      twMerge(
        `flex gap-2 items-center pl-[2.5rem] pr-4 py-2 transition-all`,
        isActive
          ? `bg-[#394B6A] font-bold`
          : `opacity-70 hover:opacity-100 hover:bg-[#394B6A] hover:bg-opacity-50`
      )
    "
  >
    <div>
      <slot />
    </div>
  </NuxtLink>
</template>
