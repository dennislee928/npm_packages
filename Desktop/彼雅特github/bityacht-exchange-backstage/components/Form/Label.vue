<script setup>
const { required, label } = defineProps({
  label: {
    type: String,
    required: true,
  },
  required: {
    type: Boolean,
    default: false,
  },
});
</script>
<template>
  <label class="w-full text-right">
    <span class="text-red-500" v-if="required">*</span>
    {{ label }}
  </label>
</template>
