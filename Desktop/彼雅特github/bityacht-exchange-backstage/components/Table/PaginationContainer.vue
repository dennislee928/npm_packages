<script setup>
import {
  DatabaseOutlined,
  LeftOutlined,
  RightOutlined,
} from "@ant-design/icons-vue";
const { cols, page, perPage, totalPages } = defineProps({
  cols: {
    type: String,
    default: "repeat(auto-fit, minmax(4rem, 1fr))",
  },
  page: {
    default: 1,
  },
  perPage: {
    default: 25,
  },
  totalPages: {
    default: 1,
  },
});
</script>
<template>
  <div :style="`--cols: ${cols}`" class="m-6">
    <div class="overflow-x-auto">
      <div class="min-w-[1000px]">
        <slot />
      </div>
    </div>
    <div
      v-if="!totalPages"
      class="border border-[#DDE1E5] border-t-0 bg-white py-16 px-4 items-center"
    >
      <div class="text-center text-[#C5C6C9] text-4xl mb-2">
        <database-outlined />
      </div>
      <div class="text-center text-[#C5C6C9]">查無資料</div>
    </div>
    <div class="flex gap-2 items-center mt-2" v-if="totalPages > 0">
      <button
        class="bg-[#DDE1E5] bg-opacity-0 hover:bg-opacity-50 active:bg-opacity-100 rounded-sm focus:outline-none w-8 h-8 text-center disabled:opacity-50 disabled:pointer-events-none"
        :disabled="page === 1"
        @click="$emit('update:page', page - 1)"
      >
        <LeftOutlined />
      </button>
      <select
        :value="page"
        @input="$emit('update:page', $event.target.value)"
        class="border border-[#DDE1E5] rounded-sm focus:outline-none h-8"
      >
        <option v-for="i in totalPages" :key="i" :value="i">
          {{ i }}
        </option>
      </select>
      <button
        class="bg-[#DDE1E5] bg-opacity-0 hover:bg-opacity-50 active:bg-opacity-100 rounded-sm focus:outline-none w-8 h-8 text-center disabled:opacity-50 disabled:pointer-events-none"
        :disabled="page === totalPages"
        @click="$emit('update:page', page + 1)"
      >
        <RightOutlined />
      </button>

      <select
        :value="perPage"
        @input="$emit('update:perPage', $event.target.value)"
        class="border border-[#DDE1E5] rounded-sm focus:outline-none h-8"
      >
        <option v-for="i in [25, 50, 100]" :key="i" :value="i">
          {{ i }} 筆 / 頁
        </option>
      </select>
    </div>
  </div>
</template>
