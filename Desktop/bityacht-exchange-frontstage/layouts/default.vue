<template>
  <div class="bg-darkBlueBg overflow-hidden">
    <BaseHeader v-if="!appMode" class="3xl:px-[300px] xl:px-[150px] lg:px-[100px] px-[50px]" />
    <slot></slot>
    <BaseFooter />
  </div>
</template>
<script setup>
const appMode = ref(false);
const mode = localStorage.getItem('mode');

onMounted(() => {
  if (mode === 'app') {
    appMode.value = true;
  }
});
</script>
