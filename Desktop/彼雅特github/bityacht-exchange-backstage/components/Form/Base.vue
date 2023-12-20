<script setup>
const { label } = defineProps(["label"]);
</script>
<template>
  <div>
    <label :for="label" class="text-sm color-[#6B6C6C] pl-2" v-if="label">
      {{ label }}
    </label>
    <slot />
  </div>
</template>
