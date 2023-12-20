<script setup>
const route = useRoute();
const { uid } = route.params;
</script>
<template>
  <Transactions :id="uid" />
</template>
