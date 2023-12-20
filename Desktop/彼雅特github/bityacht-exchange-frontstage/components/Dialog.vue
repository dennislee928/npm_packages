<script setup>
import { CloseOutlined } from '@ant-design/icons-vue';
const { modelValue, title } = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  title: {
    type: String,
    default: '標題',
  },
});
</script>
<template>
  <Teleport to="body">
    <div
      class="fixed inset-0 m-auto h-full w-full bg-black bg-opacity-50 flex items-center justify-center cursor-pointer z-[100] p-4"
      :style="{
        pointerEvents: modelValue ? 'auto' : 'none',
        opacity: modelValue ? 1 : 0,
        transition: 'opacity 0.1s ease-in-out',
      }">
      <div
        class="max-w-[512px] cursor-auto shadow-xl max-h-[calc(100vh-4rem)] flex flex-col"
        @click.stop
        :style="{
          transform: modelValue ? 'scale(1)' : 'scale(0.8)',
          transition: 'transform 0.15s ease-in-out',
        }">
        <!-- <div class="cursor-pointer flex items-center justify-center w-8 h-8 rounded-full hover:bg-white hover:bg-opacity-10 active:bg-opacity-20" @click="$emit('update:modelValue', false)">
          <CloseOutlined class="leading-4" />
        </div> -->

        <div class="py-4 px-8 bg-white flex justify-center overflow-y-auto rounded-xl border border-waterBlue">
          <slot />
        </div>
      </div>
    </div>
  </Teleport>
</template>
