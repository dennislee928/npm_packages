<template>
  <div class="relative border border-gray_800 rounded-xl bg-white min-h-[275px] p-6 mb-10">
    <div :class="`absolute md:-top-4 -top-3 left-5 bg-cover md:text-h1 text-h4 text-white font-medium  flex justify-center items-center ${isBlue ? `bg-[url('assets/img/labelBlue.png')] xs:w-[140px] xs:h-[135px] w-[110px] h-[100px]` : `bg-[url('assets/img/labelYellow.png')] xs:w-[120px] xs:h-[128px] w-[90px] h-[98px]`}`"><span class="mr-1">LV</span>{{ level }}</div>
    <div class="absolute top-[28px] right-[28px]">
      <div class="flex flex-col">
        <template v-if="realLevel < level">
          <img src="/icon/lock.svg" />
          <p class="text-small text-gray_500 mt-1">{{ $t('memberCenter.notAccept') }}</p>
        </template>
        <template v-else>
          <img src="/icon/unlock.svg" />
          <p class="text-small text-gray_500 mt-1">{{ $t('memberCenter.accept') }}</p>
        </template>
      </div>
    </div>
    <div class="xs:mt-[130px] mt-[100px]">
      <slot />
    </div>
  </div>
</template>
<script setup>
const { level, isBlue, realLevel } = defineProps({
  level: {
    type: Number || String,
  },
  isBlue: {
    type: Boolean,
    default: false,
  },
  realLevel: {
    type: Number || String,
  },
});
</script>
