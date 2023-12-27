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
  haveIcon: {
    type: Boolean,
    default: false,
  },
  notice: {
    type: String,
    default: '',
  },
});
</script>
<template>
  <label class="w-full flex font-semibold">
    <img src="/icon/hexagon.svg" alt="icon" class="mr-2" v-if="haveIcon" />
    <span class="text-red mr-1" v-if="required">*</span>
    {{ $t(label) }}
    <span v-if="notice" class="text-red ml-2">{{ $t(notice) }}</span>
  </label>
</template>
