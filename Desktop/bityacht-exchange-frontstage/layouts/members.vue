<template>
  <section class="relative 3xl:px-[300px] xl:px-[200px] lg:px-[100px] px-[30px] bg-white">
    <h1 class="text-h2 text-gray_800 text-center pt-[88px] pb-[60px]">{{ $t('memberCenter.memberCenterTitle') }}</h1>
    <div class="relative z-30 border border-gray_800 rounded-xl md:p-10 p-5">
      <div class="flex md:flex-row flex-col justify-around">
        <div class="md:flex hidden lg:mb-0 xs:mb-10 mb-6">
          <div class="relative flex flex-col items-center justify-end">
            <a-avatar :size="70" class="bg-waterBlue">
              <template #icon><UserOutlined /></template>
            </a-avatar>
            <template v-if="dataReady">
              <WarningFilled v-if="userData.status === 3" class="text-red text-h4 absolute -top-1 -right-2" />
            </template>
            <Button variant="whiteBlue" v-if="dataReady" class="mt-3 break-keep !py-[1px]"> LV {{ userData.level }}</Button>
          </div>
          <div class="flex flex-col justify-center ml-3 mb-5">
            <p v-if="dataReady" class="text-black text-normal font-bold">{{ userData.account }}</p>
            <p v-if="dataReady" class="text-grayText text-small">UID :{{ userData.id }}</p>
            <!-- <Button variant="waterBlue" class="flex items-center justify-center mt-4 w-[170px]"><img src="/icon/gift.svg" class="mr-2" />{{ $t('memberCenter.shareInviteCode') }}</Button> -->
          </div>
        </div>
        <div class="md:hidden flex lg:mb-0 xs:mb-10 mb-6">
          <div class="relative">
            <a-avatar :size="70" class="bg-waterBlue">
              <template #icon><UserOutlined /></template>
            </a-avatar>
            <template v-if="dataReady">
              <WarningFilled v-if="userData.status === 3" class="text-red text-h4 absolute -top-1 -right-3" />
            </template>
          </div>
          <div class="flex flex-col ml-5 justify-around">
            <div class="flex">
              <Button variant="whiteBlue" v-if="dataReady" class="break-keep mr-2 md:text-normal text-small"> LV {{ userData.level }}</Button>
              <!-- <Button variant="waterBlue" class="flex items-center xs:text-normal text-small"><img src="/icon/gift.svg" class="mr-2" />{{ $t('memberCenter.shareInviteCode') }}</Button> -->
            </div>
            <p v-if="dataReady" class="text-black text-normal font-bold mt-1">{{ userData.account }}</p>
            <p v-if="dataReady" class="text-grayText text-small">UID :{{ userData.id }}</p>
          </div>
        </div>
        <div class="flex lg:flex-row md:flex-col flex-row items-center border border-gray_300 rounded-xl p-5 lg:justify-start md:justify-end justify-center md:mb-0 xs:mb-10 mb-6">
          <img src="/icon/bank01.svg" class="md:mr-1 xs:mr-8 mr-4" />
          <div v-if="dataReady" class="flex 3xl:flex-row flex-col items-center">
            <template v-if="userStatus.idVerificationStatus === 1">
              <div class="flex flex-col break-keep mx-3">
                <p class="text-black text-normal font-bold">{{ $t('memberCenter.IDVerify') }}{{ $t('memberCenter.waitingForAccept') }}</p>
                <p class="text-grayText text-small mt-2">{{ $t('memberCenter.needsTime') }}</p>
              </div>
            </template>
            <template v-else-if="userStatus.idVerificationStatus === 2">
              <div class="flex flex-col break-keep mx-3">
                <p class="text-black text-normal font-bold">
                  {{ $t('memberCenter.IDVerify') }}<span class="text-green">{{ $t('memberCenter.accept') }}</span>
                </p>
                <p v-if="userStatus.bankAccountStatus !== 2" class="text-grayText text-subText mt-2">{{ $t('memberCenter.goBankToBind') }}</p>
              </div>
            </template>
            <template v-else-if="userStatus.idVerificationStatus === 3">
              <div class="flex flex-col break-keep mx-3">
                <p class="text-black text-normal font-bold">
                  {{ $t('memberCenter.IDVerify') }}<span class="text-ref">{{ $t('memberCenter.reject') }}</span>
                </p>
              </div>
            </template>
            <template v-else-if="userStatus.idVerificationStatus === 4">
              <div class="flex flex-col break-keep mx-3">
                <p class="text-black text-normal font-bold">
                  {{ $t('memberCenter.IDVerify') }}<span class="text-ref">{{ $t('memberCenter.reject') }}</span>
                </p>
                <p class="text-grayText text-small">{{ $t('memberCenter.pleaseIDVerifyAgain') }}</p>
              </div>
            </template>
            <template v-else>
              <div class="flex flex-col break-keep mx-3">
                <p class="text-black text-normal font-bold">{{ $t('memberCenter.IDVerify') }}{{ $t('memberCenter.notAccept') }}</p>
                <p class="text-grayText text-small">{{ $t('memberCenter.pleaseIDVerifyFirst') }}</p>
              </div>
            </template>
            <NuxtLink v-if="userStatus.idVerificationStatus === 0 && userData.status !== 3" :to="localePath(`/Members/${uid}/verify`)">
              <Button variant="waterBlue" class="3xl:mt-0 mt-2">{{ $t('memberCenter.IDVerify') }}</Button>
            </NuxtLink>
            <NuxtLink v-if="userStatus.idVerificationStatus === 4 && userData.countriesCode !== 'TWN'" :to="localePath(`/Members/${uid}/completeImage`)">
              <Button variant="waterBlue" class="3xl:mt-0 mt-2">{{ $t('memberCenter.reIDVerify') }}</Button>
            </NuxtLink>
            <!-- <a v-if="userStatus.idVerificationStatus === 4 && userData.countriesCode === 'TWN'" :href="kryptoGoUrl" :target="kryptoGoUrl === '' ? '' : '_blank'">
              <Button variant="waterBlue" class="3xl:mt-0 mt-2">{{ $t('memberCenter.reIDVerify') }}</Button> -->

            <template v-if="userStatus.idVerificationStatus === 4 && userData.countriesCode === 'TWN'">
              <Button variant="waterBlue" class="3xl:mt-0 mt-2" @click="getKryptoGoUrl">{{ $t('memberCenter.reIDVerify') }}</Button>
            </template>
          </div>
        </div>
        <div class="flex lg:flex-row md:flex-col flex-row items-center border border-gray_300 rounded-xl p-5 lg:justify-start md:justify-end justify-center">
          <img src="/icon/bank.svg" class="md:mr-1 xs:mr-8 mr-4" />
          <template v-if="dataReady">
            <div v-if="userStatus.bankAccountStatus === 2" class="flex 3xl:flex-row flex-col">
              <div class="flex flex-col break-keep mx-3">
                <p class="text-black text-normal font-bold">
                  {{ $t('memberCenter.bankAccount') }}<span class="text-green">{{ $t('memberCenter.accept') }}</span>
                </p>
              </div>
            </div>
            <div v-else class="flex 3xl:flex-row flex-col justify-center items-center">
              <div class="flex flex-col break-keep mx-3">
                <p class="text-black text-normal font-bold">{{ $t('memberCenter.unLockLegalTender') }}</p>
                <template v-if="dataReady">
                  <p v-if="userData.status === 3" class="text-grayText text-small">{{ $t('memberCenter.accountFreeze') }}</p>
                  <p v-else class="text-grayText text-small">{{ $t('memberCenter.pleaseIDVerifyFirst') }}</p>
                </template>
              </div>
              <template v-if="dataReady">
                <Button v-if="userData.status === 3" variant="disabled" class="3xl:mt-0 mt-2">{{ $t('memberCenter.contractService') }}</Button>
                <NuxtLink v-else :to="localePath(`/Members/${uid}/bankAccount`)">
                  <Button variant="waterBlue" class="3xl:mt-0 mt-2">{{ $t("trade.comingSoon") }}</Button>
                </NuxtLink>
              </template>
            </div>
          </template>
        </div>
      </div>
    </div>
    <div class="relative flex justify-center md:flex-nowrap flex-wrap py-10 z-30">
      <LayoutsNavButton :href="`/Members/${uid}/level`" class="break-keep xs:mr-5 mr-2 mb-2 xs:text-normal text-small"> {{ $t('memberCenter.certificationLevel') }} </LayoutsNavButton>
      <LayoutsNavButton :href="`/Members/${uid}/safe`" class="break-keep xs:mr-5 mr-2 mb-2 xs:text-normal text-small"> {{ $t('memberCenter.accountSafe') }} </LayoutsNavButton>
      <LayoutsNavButton :href="`/Members/${uid}/bankAccount`" class="break-keep xs:mr-5 mr-2 mb-2 xs:text-normal text-small"> {{ $t('memberCenter.bankAccount') }} </LayoutsNavButton>
      <LayoutsNavButton :href="`/Members/${uid}/perferenceSetting`" class="break-keep xs:mr-5 mr-2 mb-2 xs:text-normal text-small"> {{ $t('memberCenter.preferenceSettings') }} </LayoutsNavButton>
      <LayoutsNavButton :href="`/Members/${uid}/loginRecord`" class="break-keep xs:mr-5 mr-2 mb-2 xs:text-normal text-small"> {{ $t('memberCenter.loginRecord') }} </LayoutsNavButton>
      <LayoutsNavButton :href="`/Members/${uid}/inviteFriend`" class="break-keep xs:mr-5 mr-2 mb-2 xs:text-normal text-small"> {{ $t('memberCenter.inviteFriends') }} </LayoutsNavButton>
    </div>
    <img src="/assets/img/Group16.png" class="absolute right-0 top-0 md:w-[600px] w-[400px] z-10" />
    <img src="/assets/img/Group37.png" class="absolute left-0 top-0 md:w-[800px] w-[200px] z-10" />
  </section>
  <slot />
</template>
<script setup>
import useUserStore from '@/stores/user';

const localePath = useLocalePath();
const { t } = useI18n();
useHead({
  title: t('title.memberCenter'),
  meta: [{ name: 'description', content: '' }],
});
const userStore = useUserStore();
const route = useRoute();
const { uid } = route.params;
const dataReady = ref(false);
const userData = ref();
const userStatus = ref();
const getUserInfo = async () => {
  const result = await userStore.getUserInfo(t);
  userStatus.value = result.data.value;
  // console.log('userStatus.value :>> ', userStatus.value);
};
const kryptoGoUrl = ref('');
const getKryptoGoUrl = async () => {
  const result = await userStore.kryptoGoUrl(t);
  if (result.status.value === 'success') {
    window.open(result.data.value.idVerificationUrl, '_blank');
  }
};
onMounted(async () => {
  await getUserInfo();
  const userInfo = JSON.parse(localStorage.getItem('userInfo'));
  userData.value = userInfo;
  // console.log('userData.value :>> ', userData.value);
  dataReady.value = true;
  if (userData.value.type === 1) {
    if ((userStatus.value.idVerificationStatus === 0 || userStatus.value.idVerificationStatus === 1 || userStatus.value.idVerificationStatus === 3) && userData.value.level > 0) {
      await userStore.token(true);
      const userInfo = JSON.parse(localStorage.getItem('userInfo'));
      userData.value = userInfo;
    } else if (userStatus.value.idVerificationStatus === 2 && userData.value.level === 0) {
      await userStore.token(true);
      const userInfo = JSON.parse(localStorage.getItem('userInfo'));
      userData.value = userInfo;
    } else if (userStatus.value.idVerificationStatus === 4 && userData.value.countriesCode === '') {
      await userStore.token(true);
      const userInfo = JSON.parse(localStorage.getItem('userInfo'));
      userData.value = userInfo;
    }
  }
  // if (userStatus.value.idVerificationStatus === 4) {
  //   if (userData.value.countriesCode === 'TWN') {
  //     const result = await userStore.kryptoGoUrl(t);
  //     kryptoGoUrl.value = result.data.value.idVerificationUrl;
  //   } else {
  //   }
  // }
});
</script>
