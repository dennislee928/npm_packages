<script setup>
import { QuestionOutlined } from "@ant-design/icons-vue";
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
        `flex gap-2 items-center px-4 py-2.5`,
        isActive
          ? `bg-[#394B6A] font-bold`
          : `opacity-70 hover:opacity-100 hover:bg-[#394B6A] hover:bg-opacity-50`
      )
    "
  >
    <slot name="icon">
      <question-outlined />
    </slot>
    <div>
      <slot />
    </div>
  </NuxtLink>
</template>
