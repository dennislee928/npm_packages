<script setup>
const localePath = useLocalePath();
const { locale } = useI18n();
const { href } = defineProps({
  href: {
    type: String,
    required: true,
  },
});
const isActive = ref(locale.value === 'en-US' ? useRoute().path === '/' + locale.value + href : useRoute().path === href);
watch(
  () => useRoute().path,
  (newValue) => {
    isActive.value = newValue === href;
  }
);
const isMobile = ref(null);
onMounted(() => {
  const screen = window.innerWidth;
  if (screen < 768) {
    isMobile.value = true;
  }
});
</script>

<template>
  <NuxtLink :to="localePath(`${href}`)">
    <Button :variant="`${isActive ? 'waterBlue' : 'white'}`">
      <slot />
    </Button>
  </NuxtLink>
</template>
