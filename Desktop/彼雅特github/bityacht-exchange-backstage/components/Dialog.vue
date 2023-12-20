<script setup>
import { CloseOutlined } from "@ant-design/icons-vue";
import { Collapse } from "vue-collapsed";
import { twMerge } from "tailwind-merge";
const { modelValue, title, loading, variant } = defineProps({
  modelValue: {
    default: false,
  },
  title: {
    type: String,
    default: "標題",
  },
  loading: {
    type: Boolean,
    default: false,
  },
  maxWidth: {
    type: String,
    default: "512px",
  },
  variant: {
    type: String,
    default: "default",
  },
});
</script>
<template>
  <Teleport to="body">
    <div
      :class="
        twMerge(
          `fixed inset-0 m-auto h-full w-full bg-black bg-opacity-75 flex cursor-pointer z-10`,
          variant === 'sidebar'
            ? `items-center justify-end`
            : `items-center justify-center p-4`
        )
      "
      @click="$emit('update:modelValue', false)"
      :style="{
        pointerEvents: modelValue ? 'auto' : 'none',
        opacity: modelValue ? 1 : 0,
        transition: 'opacity 0.25s cubic-bezier(0.33, 1, 0.68, 1)',
      }"
    >
      <div
        :class="
          twMerge(
            `w-full cursor-auto shadow-xl flex flex-col bg-white`,
            variant === 'sidebar'
              ? `h-full`
              : `max-h-[calc(100vh-4rem)] rounded`
          )
        "
        @click.stop
        :style="{
          transform:
            variant === 'sidebar'
              ? modelValue
                ? 'translateX(0)'
                : 'translateX(25%)'
              : modelValue
              ? 'scale(1)'
              : 'scale(0.95)',
          opacity: modelValue ? 1 : 0,
          transition: 'all 0.25s cubic-bezier(0.33, 1, 0.68, 1)',
          maxWidth: maxWidth || '512px',
        }"
      >
        <a-spin
          :spinning="loading"
          :wrapperClassName="
            twMerge(
              `[&>.ant-spin-container]:flex [&>.ant-spin-container]:flex-col`,
              variant === 'sidebar'
                ? `[&>.ant-spin-container]:h-[100svh] [&>.ant-spin-container]:w-full `
                : `[&>.ant-spin-container]:max-h-[calc(100vh-4rem)]`
            )
          "
        >
          <div
            :class="
              twMerge(
                `flex justify-between items-center pl-8 pr-2 py-3 text-white`,
                variant === 'sidebar'
                  ? `bg-[#19253A] py-1.5`
                  : `bg-[#8E9BA5] rounded-t`
              )
            "
            v-if="variant !== 'image'"
          >
            <div class="text-xl font-bold">{{ title }}</div>
            <div
              class="cursor-pointer flex items-center justify-center w-8 h-8 rounded-full hover:bg-white hover:bg-opacity-10 active:bg-opacity-20"
              @click="$emit('update:modelValue', false)"
            >
              <CloseOutlined class="leading-4" />
            </div>
          </div>
          <div
            :class="
              twMerge(
                `py-4 px-8 bg-white flex-1 overflow-x-hidden overflow-y-auto h-full`,
                variant === 'image' ? 'rounded-t' : ''
              )
            "
          >
            <Collapse
              :when="!!modelValue"
              :style="{
                transition:
                  !modelValue && 'all 0.25s cubic-bezier(0.33, 1, 0.68, 1)',
              }"
              v-if="variant !== 'sidebar'"
            >
              <slot />
            </Collapse>
            <slot v-else />
          </div>
          <div
            :class="
              twMerge(
                `py-2 px-8 bg-white rounded-b flex items-center justify-end gap-2 border-t border-[#DDE1E5]`,
                `[&>.btn]:py-2 [&>.btn]:px-5`
              )
            "
          >
            <slot name="actions" />
          </div>
        </a-spin>
      </div>
    </div>
  </Teleport>
</template>
