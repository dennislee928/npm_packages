<script setup>
import { twMerge } from "tailwind-merge";
const { modelValue, label, placeholder, type, required, disabled } =
  defineProps([
    "modelValue",
    "label",
    "placeholder",
    "type",
    "required",
    "accept",
    "disabled",
  ]);
</script>
<template>
  <FormBase :label="label">
    <input
      :value="modelValue"
      :class="
        twMerge(
          `bg-white border border-[#DDE1E5] focus:border-[#00b96b] focus:ring-2 ring-[#00b96b] ring-opacity-10 rounded-[8px] px-4 py-2 h-[44px] w-full outline-0 transition-colors`,
          `file:mr-2 file:py-1 file:px-3 file:rounded-[4px] file:border-0 file:text-sm file:bg-[#00b96b] file:bg-opacity-10 file:text-[#00b96b] file:cursor-pointer file:hover:bg-opacity-20 file:active:bg-opacity-30 file:-ml-2`,
          disabled &&
            `bg-[#ebebeb] text-[#cec1c1] cursor-not-allowed pointer-events-none`
        )
      "
      @input="$emit('update:modelValue', $event.target.value)"
      :placeholder="placeholder"
      :id="label"
      :type="type"
      :required="required"
      :accept="accept"
      :disabled="disabled"
    />
  </FormBase>
</template>
