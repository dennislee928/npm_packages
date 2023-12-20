<template>
  <div class="relative bg-white flex flex-col items-center">
    <h1 class="md:text-h2 text-h3 text-pink pt-[120px] md:mb-[55px] mb-[35px]">{{ $t('signUp.accountSuspension') }}</h1>
    <div class="bg-white border border-pink rounded-md py-[35px] md:px-[60px] px-[45px] md:w-[450px] xxs:w-[350px] w-[300px] md:mb-[175px] mb-[75px] z-10 flex flex-col items-center">
      <img src="/assets/img/exclamationMark.png" class="md:w-[120px] w-[200px]" />
      <p class="text-pink my-8 text-center">{{ $t('signUp.suspensionNotice') }}</p>
      <NuxtLink :to="localePath('/')"
        ><Button variant="waterBlue" class="w-[100px] shadow-xl shadow-gray_700/20">{{ $t('signUp.back') }}</Button></NuxtLink
      >
      <p class="text-gray_500 text-small mt-6">
        <NuxtLink :to="localePath('/login')" class="text-small">{{ $t('signUp.contactCustomerService') }}</NuxtLink>
      </p>
    </div>
    <img src="/assets/img/Group16.png" class="absolute right-0 md:w-[640px] w-[400px]" />
    <img src="/assets/img/Group37.png" class="absolute left-0 top-0 md:w-[500px] w-[100px]" />
  </div>
</template>

<script setup>
const localePath = useLocalePath();
</script>
