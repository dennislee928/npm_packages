<script setup>
const localePath = useLocalePath();
const { locale } = useI18n();
const { href } = defineProps({
  href: {
    type: String,
    required: true,
  },
});
const isActive = ref(useRoute().path.startsWith(href.split('/')[1], locale.value === 'en-US' ? 7 : 1));
// console.log('isActive :>> ', isActive.value);
// console.log('href :>> ', href);
// console.log('useRoute() :>> ', useRoute().path);
// const isActive = ref(useRoute().path === href);
watch(
  () => useRoute().path,
  (newValue) => {
    isActive.value = newValue.startsWith(href.split('/')[1], locale.value === 'en-US' ? 7 : 1);
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
  <NuxtLink
    :to="localePath(`${href}`)"
    :class="`p-1 transition-all duration-300 cursor-pointer hover:text-waterBlue
    ${isActive ? `text-waterBlue font-bold` : `text-black`} ${isActive && isMobile ? `border-b-0` : ``}`">
    <slot />
  </NuxtLink>
</template>
