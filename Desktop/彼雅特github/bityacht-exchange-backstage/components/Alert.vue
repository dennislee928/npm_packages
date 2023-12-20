<script setup>
import { UpOutlined, DownOutlined } from "@ant-design/icons-vue";
const showDatail = ref(false);
</script>
<template>
  <div class="bg-white py-3 px-6 border-2 border-[#E1EFF8] rounded-xl">
    <div
      class="flex justify-between items-center cursor-pointer"
      @click="showDatail = !showDatail"
      :aria-expanded="showDatail"
      aria-controls="collapse"
    >
      <div class="flex gap-2 font-bold items-center">
        <slot />
      </div>
      <UpOutlined v-if="showDatail" />
      <DownOutlined v-else />
    </div>
    <div class="mt-2" v-if="showDatail">
      <slot name="content" />
    </div>
  </div>
</template>
