<script setup>
import { twMerge } from "tailwind-merge";
import { Collapse } from "vue-collapsed";
import { QuestionOutlined, DownOutlined } from "@ant-design/icons-vue";
const open = ref(false);
</script>
<template>
  <div
    class="flex gap-2 items-center px-4 py-2.5 transition-all opacity-70 hover:opacity-100 hover:bg-[#394B6A] hover:bg-opacity-50 cursor-pointer select-none"
    @click="open = !open"
  >
    <slot name="icon">
      <question-outlined />
    </slot>
    <div class="flex-1">
      <slot name="name" />
    </div>
    <DownOutlined
      :class="
        twMerge(`text-gray-300 transition-transform`, !open && `-rotate-90`)
      "
    />
  </div>
  <Collapse :when="open">
    <slot />
  </Collapse>
</template>
