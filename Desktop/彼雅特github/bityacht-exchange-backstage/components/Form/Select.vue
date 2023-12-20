<script setup>
const { modelValue, label } = defineProps(["modelValue", "label"]);
</script>
<template>
  <FormBase :label="label">
    <select
      :value="modelValue"
      class="bg-white border border-[#DDE1E5] focus:border-[#00b96b] focus:ring-2 ring-[#00b96b] ring-opacity-10 rounded-[8px] px-4 py-2 w-full outline-0 transition-colors h-[44px]"
      @input="$emit('update:modelValue', $event.target.value)"
      :id="label"
    >
      <slot />
    </select>
  </FormBase>
</template>
