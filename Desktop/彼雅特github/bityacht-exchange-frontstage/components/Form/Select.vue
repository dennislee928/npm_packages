<script setup>
const { modelValue, label, disabled } = defineProps(['modelValue', 'label', 'disabled']);
</script>
<template>
  <FormBase :label="label">
    <select :value="modelValue" :disabled="disabled" class="bg-white border border-[#DDE1E5] focus:border-waterBlue focus:ring-2 ring-waterBlue ring-opacity-10 rounded-lg pl-3 pr-5 py-2.5 w-full outline-0 transition-colors" @input="$emit('update:modelValue', $event.target.value)" :id="label">
      <slot />
    </select>
  </FormBase>
</template>
