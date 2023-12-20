<template>
  <div class="relative bg-white flex flex-col items-center">
    <h1 class="md:text-h2 text-h4 text-gray_800 pt-[120px] md:mb-[55px] mb-[35px]">{{ $t(`${title}`) }}</h1>
    <div class="bg-white border border-gray_800 rounded-md py-[45px] md:px-[60px] px-[45px] md:w-[650px] xxs:w-[450px] w-[300px] md:mb-[175px] mb-[75px] z-10 flex flex-col items-center">
      <img v-if="imgSource" :src="getAssetsFile(`${imgSource}.png`)" />
      <template v-if="type === 'withdraw'">
        <div class="flex flex-col mt-5 md:text-normal text-subText">
          <p class="flex my-1 text-gray_800">
            <span class="md:w-[100px] w-[90px]">{{ $t('myAssets.thisCount') }}</span> : {{ amount }} <span class="ml-2">{{ currencyType }}</span>
          </p>
          <p class="flex my-1 text-gray_800">
            <span class="md:w-[100px] w-[90px]">{{ $t('myAssets.canUseCount') }}</span> : {{ freeAmount }}<span class="ml-2">{{ currencyType }}</span>
          </p>
          <p class="flex my-1 text-gray_800">
            <span class="md:w-[100px] w-[90px]">{{ $t('myAssets.lockCount') }}</span> : {{ lockedAmount }} <span class="ml-2">{{ currencyType }}</span>
          </p>
        </div>
      </template>
      <a-button v-if="url" :href="url" class="w-full h-auto bg-waterBlue md:text-h5 text-[20px] text-white hover:!text-white py-6 rounded-xl">{{ $t('AuthVerify.kryptoHref') }}</a-button>
      <p v-if="url" class="mt-10">{{ $t('AuthVerify.kryptoText') }}</p>
      <NuxtLink :to="localePath(`/Members/${uid}/level`)"
        ><Button variant="waterBlue" class="w-[160px] shadow-xl shadow-gray_700/20 md:mt-[50px] mt-[30px]">{{ $t('signUp.backToMemberCenter') }}</Button></NuxtLink
      >
    </div>
    <img src="/assets/img/Group16.png" class="absolute right-0 md:w-[640px] w-[400px]" />
    <img src="/assets/img/Group37.png" class="absolute left-0 top-0 md:w-[500px] w-[100px]" />
  </div>
</template>

<script setup>
import useAssetsStore from '@/stores/assets';
import { getAssetsFile } from '@/config/config';

const { t } = useI18n();
const assetsStore = useAssetsStore();
const getAssets = async () => {
  const result = await assetsStore.getAssets(t);
  result.data.value.assetInfos.forEach((item) => {
    if (item.currenciesSymbol === currencyType.value) {
      freeAmount.value = item.freeAmount;
      lockedAmount.value = item.lockedAmount;
    }
  });
};
const localePath = useLocalePath();
const route = useRoute();
const isLogin = useCookie('isLogin');
const title = ref(route.query.title);
const imgSource = ref(route.query.imgSource);
const type = ref(route.query.type);
const url = ref(route.query.url);
const amount = ref(route.query.amount);
const freeAmount = ref();
const lockedAmount = ref();
const currencyType = ref(route.query.currencyType);
// const status = ref('');
const uid = ref('');

onMounted(async () => {
  if (isLogin.value !== 0) {
    const userInfo = JSON.parse(localStorage.getItem('userInfo'));
    uid.value = userInfo.id;
  }
  if (type.value === 'withdraw') {
    await getAssets();
  }
});
</script>
