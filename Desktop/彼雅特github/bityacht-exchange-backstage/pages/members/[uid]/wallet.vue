<script setup>
const route = useRoute();
const { uid } = route.params;
</script>
<template>
  <Spots :id="uid" />
</template>
