<script setup>
const { label } = defineProps(['label']);
</script>
<template>
  <div>
    <label :for="label" class="text-sm color-[#6B6C6C]">{{ label }}</label>
    <slot />
  </div>
</template>
