<template>
  <div class="bg-white py-4 md:relative fixed w-full z-[50] shadow-lg shadow-gray_700/20">
    <nav class="md:flex hidden justify-between items-center">
      <div class="flex items-center w-[50%]">
        <NuxtLink :to="localePath('/')">
          <img src="/img/logo.svg" alt="logo" class="lg:mr-12 mr-6 cursor-pointer" />
        </NuxtLink>
        <LayoutsNavItem :href="`/MyAssets`" class="lg:mr-6 mr-3 break-keep xl:text-normal text-subText"> {{ $t('index.myAssets') }} </LayoutsNavItem>
        <LayoutsNavItem :href="`/Trade`" class="lg:mr-6 mr-3 break-keep xl:text-normal text-subText"> {{ $t('index.business') }} </LayoutsNavItem>
        <LayoutsNavItem :href="`/Members/${uid}/level`" class="break-keep xl:text-normal text-subText">
          {{ $t('index.userCenter') }}
        </LayoutsNavItem>
      </div>
      <div class="flex w-[50%] items-center justify-end break-keep">
        <!-- <Button variant="waterBlueBorder" class="xl:mr-10 mr-2 flex items-center"
          ><DownloadOutlined /> <span class="xl:block hidden">{{ $t('index.downloadAPP') }}</span></Button
        > -->
        <template v-if="isLogin">
          <Button variant="outline" class="md_lg:mr-6 mr-4 xl:w-[100px] w-[80px]" @click="dialogOpen('transfer')">{{ $t('index.transferMoney') }}</Button>
          <a-avatar class="bg-waterBlue scale-150 md_lg:mr-8 mr-2 cursor-pointer" @click="dialogOpen('user')">
            <template #icon><UserOutlined /></template>
          </a-avatar>
        </template>
        <template v-else>
          <NuxtLink :to="localePath('/login')">
            <Button variant="white" class="mr-3">{{ $t('index.logIn') }}</Button>
          </NuxtLink>
          <NuxtLink :to="localePath('/signUp')">
            <Button variant="waterBlue" class="">{{ $t('index.signUp') }}</Button>
          </NuxtLink>
        </template>
        <!-- <div class="flex md_lg:ml-6 ml-3">
          <div class="xl:w-[40px] w-[30px] xl:h-[40px] h-[30px] border border-black rounded-full flex items-center justify-center mx-2 cursor-pointer" :class="{ 'bg-waterBlue text-white border-none': isLocale === 'zh-TW' }" @click="setLanguage('zh-TW')">繁</div>
          <span class="block border border-gray_400"></span>
          <div class="xl:w-[40px] w-[30px] xl:h-[40px] h-[30px] border border-black rounded-full flex items-center justify-center mx-2 cursor-pointer" :class="{ 'bg-waterBlue text-white border-none': isLocale === 'en-US' }" @click="setLanguage('en-US')">EN</div>
        </div> -->
      </div>
    </nav>
    <Transition>
      <div v-if="userDialog" ref="userDialogRef" class="bg-white min-w-[300px] fixed top-[100px] right-[100px] p-8 rounded-2xl shadow shadow-blue-500/40">
        <div class="relative">
          <CloseOutlined class="cursor-pointer absolute -top-5 -right-4" @click="userDialog = !userDialog" />
          <div class="flex items-center">
            <a-avatar :size="40" class="bg-waterBlue mr-4 cursor-pointer">
              <template #icon><UserOutlined /></template>
            </a-avatar>
            <span class="text-subTitle font-bold">{{ userData.data.account }}</span>
          </div>
          <!-- <LayoutsNavItem :href="`/Members/${uid}/level`" class="flex items-center mt-6" @click="userDialog = !userDialog">
            <img src="/icon/hexagon.svg" alt="icon" class="mr-4" />
            <span class="text-normal font-bold cursor-pointer">{{ $t('index.userCenter') }}</span>
          </LayoutsNavItem> -->
          <!-- <LayoutsNavItem :href="`/Gift`" class="flex items-center mt-6" @click="userDialog = !userDialog">
            <img src="/icon/hexagon.svg" alt="icon" class="mr-4" />
            <span class="text-normal font-bold cursor-pointer">{{ $t('index.goodGift') }}</span>
          </LayoutsNavItem> -->
          <div class="flex justify-end mt-6">
            <Button variant="waterBlue" class="shadow-xl shadow-gray-700/40" @click="logOut">{{ $t('index.logOut') }}</Button>
          </div>
        </div>
      </div>
    </Transition>
    <Transition>
      <div v-if="transferDialog" ref="transferDialogRef" class="bg-white min-w-[350px] fixed top-[100px] right-[100px] p-8 rounded-2xl shadow shadow-blue-500/40 border border-waterBlue">
        <div class="relative">
          <CloseOutlined class="cursor-pointer absolute -top-5 -right-4" @click="transferDialog = !transferDialog" />
          <!-- <p class="font-bold text-center text-[20px] mb-4">{{ $t('index.legel') }}</p>
          <div class="flex justify-between">
            <a-button class="w-[140px] h-[50px] border border-waterBlue rounded-3xl text-waterBlue text-subTitle p-2"><download-outlined class="text-[26px]" />{{ $t('index.inputMoney') }}</a-button>
            <a-button class="w-[140px] h-[50px] border border-waterBlue rounded-3xl text-waterBlue text-subTitle p-2"><upload-outlined class="text-[26px]" />{{ $t('index.outputMoney') }}</a-button>
          </div>
          <div class="bg-[#e3e3e3] text-small py-3 px-5 mt-5">{{ $t('index.legelText') }}</div>
          <a-divider class="border-waterBlue"></a-divider> -->
          <p class="font-bold text-center text-[20px] mb-4">{{ $t('index.digitalAssets') }}</p>
          <div class="flex justify-between">
            <NuxtLink :to="localePath('/MyAssets/cryptocurrencyTakeOver')" @click="transferDialog = !transferDialog">
              <a-button class="w-[140px] h-[50px] flex items-center justify-center border border-waterBlue rounded-3xl text-waterBlue text-subTitle p-2"><img src="/icon/inputCount.svg" class="mr-2" />{{ $t('index.inputCount') }}</a-button>
            </NuxtLink>
            <NuxtLink :to="localePath('/MyAssets/cryptocurrencyGet')" @click="transferDialog = !transferDialog">
              <a-button class="w-[140px] h-[50px] flex items-center justify-center border border-waterBlue rounded-3xl text-waterBlue text-subTitle p-2"><img src="/icon/outputCount.svg" class="mr-2" />{{ $t('index.outputCount') }}</a-button>
            </NuxtLink>
          </div>
          <div class="bg-[#e3e3e3] text-small py-3 px-5 mt-5">
            <span> {{ $t('index.pleaseDredge') }}</span>
            <span class="text-red"> {{ $t('index.IDVerify') }}</span>
            <span> ，{{ $t('index.dredgeCountAction') }}</span>
          </div>
        </div>
      </div>
    </Transition>

    <!-- mobile area -->
    <nav class="md:hidden flex justify-between items-center">
      <div>
        <NuxtLink :to="localePath('/')">
          <img src="/img/logo.svg" alt="logo" class="cursor-pointer" />
        </NuxtLink>
      </div>
      <div>
        <MenuOutlined v-if="!mobileDialog" class="cursor-pointer scale-125 pb-1 text-gray_800" @click="toggleDialog" />
        <CloseOutlined v-else class="cursor-pointer scale-125 pb-1 pr-2 text-gray_800" @click="toggleDialog" />
      </div>
    </nav>
    <a-drawer v-model:open="mobileDialog" placement="top" class="!h-screen">
      <div class="relative xs:px-8 xxs:px-4 px-2 mt-10">
        <img src="/img/whiteHexagon.svg" class="absolute right-0 top-0 w-[400px]" />
        <div class="relative z-10">
          <template v-if="isLogin">
            <div class="flex items-center">
              <a-avatar class="bg-waterBlue scale-125 mr-8">
                <template #icon><UserOutlined /></template>
              </a-avatar>
              <span class="text-subTitle font-bold">{{ userData.data.account }}</span>
            </div>
            <div class="mt-6">
              <Button variant="fullWaterBlueBorder" class="flex items-center justify-center text-normal" @click="transferMobileArea = !transferMobileArea"><DollarCircleOutlined class="scale-125 mr-2" /> {{ $t('index.transferMoney') }}</Button>
            </div>
            <Transition>
              <div v-if="transferMobileArea" class="relative border border-waterBlue p-5 mx- rounded-2xl mt-5">
                <!-- <p class="font-bold text-center text-normal mb-4">{{ $t('index.legel') }}</p> -->
                <!-- <div class="flex justify-between">
                  <a-button class="flex items-center justify-center w-[120px] h-[35px] border border-waterBlue rounded-3xl text-waterBlue text-subText p-2"><download-outlined class="text-subText" />{{ $t('index.inputMoney') }}</a-button>
                  <a-button class="flex items-center justify-center w-[120px] h-[35px] border border-waterBlue rounded-3xl text-waterBlue text-subText p-2"><upload-outlined class="text-subText" />{{ $t('index.outputMoney') }}</a-button>
                </div> -->
                <!-- <div class="bg-[#e3e3e3] text-small py-3 px-5 mt-5">{{ $t('index.legelText') }}</div> -->
                <!-- <a-divider class="border-waterBlue"></a-divider> -->
                <p class="font-bold text-center text-normal mb-4">{{ $t('index.digitalAssets') }}</p>
                <div class="flex justify-between">
                  <NuxtLink :to="localePath('/MyAssets/cryptocurrencyTakeOver')" @click="mobileDialog = !mobileDialog">
                    <a-button class="w-[120px] h-[35px] flex items-center justify-center border border-waterBlue rounded-3xl text-waterBlue text-subText p-2"><img src="/icon/inputCount.svg" class="mr-2" />{{ $t('index.inputCount') }}</a-button>
                  </NuxtLink>
                  <NuxtLink :to="localePath('/MyAssets/cryptocurrencyGet')" @click="mobileDialog = !mobileDialog">
                    <a-button class="w-[120px] h-[35px] flex items-center justify-center border border-waterBlue rounded-3xl text-waterBlue text-subText p-2"><img src="/icon/outputCount.svg" class="mr-2" />{{ $t('index.outputCount') }}</a-button>
                  </NuxtLink>
                </div>
                <div class="bg-[#e3e3e3] text-small py-3 px-5 mt-5">
                  <span> {{ $t('index.pleaseDredge') }}</span>
                  <span class="text-red"> {{ $t('index.IDVerify') }}</span>
                  <span> ，{{ $t('index.dredgeCountAction') }}</span>
                </div>
              </div>
            </Transition>
          </template>
          <LayoutsNavItem :href="`/Members/${uid}/level`" class="flex items-center mt-6">
            <img src="/icon/hexagon.svg" alt="icon" class="mr-4" />
            <span class="text-subTitle font-normal cursor-pointer" @click="toggleDialog">{{ $t('index.userCenter') }}</span>
          </LayoutsNavItem>
          <!-- <LayoutsNavItem class="flex items-center mt-6">
            <img src="/icon/hexagon.svg" alt="icon" class="mr-4" />
            <span class="text-subTitle font-normal cursor-pointer">{{ $t('index.goodGift') }}</span>
          </LayoutsNavItem> -->
          <a-divider class="border-gray_400"></a-divider>
          <LayoutsNavItem :href="`/MyAssets`" class="flex items-center mt-6">
            <img src="/icon/hexagon.svg" alt="icon" class="mr-4" />
            <span class="text-subTitle font-normal cursor-pointer" @click="toggleDialog">{{ $t('index.myAssets') }}</span>
          </LayoutsNavItem>
          <LayoutsNavItem :href="`/Trade`" class="flex items-center mt-6">
            <img src="/icon/hexagon.svg" alt="icon" class="mr-4" />
            <span class="text-subTitle font-normal cursor-pointer" @click="toggleDialog">{{ $t('index.business') }}</span>
          </LayoutsNavItem>
          <!-- <LayoutsNavItem :href="`/dayByday`" class="flex items-center mt-6">
            <img src="/icon/hexagon.svg" alt="icon" class="mr-4" @click="mobileDialog = !mobileDialog" />
            <span class="text-subTitle font-normal cursor-pointer">{{ $t('index.dayByday') }}</span>
          </LayoutsNavItem> -->
          <!-- <div class="flex justify-center mt-6 z-10">
            <div class="w-full border border-black rounded-xl xxs:mx-[50px] mx-[40px] p-4 cursor-pointer text-[20px] flex justify-center" :class="{ 'bg-waterBlue text-white border-none': isLocale === 'zh-TW' }" @click="setLanguage('zh-TW')">繁</div>
            <span class="block border border-gray_400"></span>
            <div class="w-full border border-black rounded-xl mx-[50px] p-4 cursor-pointer text-[20px] flex justify-center" :class="{ 'bg-waterBlue text-white border-none': isLocale === 'en-US' }" @click="setLanguage('en-US')">EN</div>
          </div> -->
          <div class="relative w-full bg-waterBlue rounded-2xl px-6 mt-[45px] z-10">
            <p class="text-h5 text-white pt-10">{{ $t('index.dowloadBitYacht') }}</p>
            <p class="text-subTitle text-white font-light mt-1">{{ $t('index.bestParner') }}</p>
            <div class="flex mt-5 pb-6">
              <img src="/assets/img/appStore.png" alt="appStore" class="mr-3 cursor-pointer" />
              <img src="/assets/img/googlePlay.png" alt="googlePlay" class="cursor-pointer" />
            </div>
          </div>
          <template v-if="isLogin">
            <div class="flex items-center justify-center my-[65px]">
              <Button variant="darkBlue" class="w-[123px] text-subTitle" @click="logOut">{{ $t('index.logOut') }}</Button>
            </div>
          </template>
          <template v-else
            ><div class="flex items-center justify-center my-[65px]">
              <Button variant="waterBlue" class="w-[123px] mr-7 text-subTitle" @click="linkTo('/signUp')">{{ $t('index.signUp') }}</Button>
              <Button variant="darkBlue" class="w-[123px] text-subTitle" @click="linkTo('/login')">{{ $t('index.logIn') }}</Button>
            </div></template
          >
        </div>
      </div>
    </a-drawer>
  </div>
</template>

<script setup>
import { DownloadOutlined, MenuOutlined, CloseOutlined, UserOutlined, DollarCircleOutlined } from '@ant-design/icons-vue';
import { parseJwt, clearAllCookie } from '@/config/config';
import useUserStore from '@/stores/user';
import { onClickOutside } from '@vueuse/core';

const userDialogRef = ref(null);
const transferDialogRef = ref(null);
onClickOutside(userDialogRef, () => (userDialog.value = false));
onClickOutside(transferDialogRef, () => (transferDialog.value = false));
const userStore = useUserStore();
const isLogin = useCookie('isLogin');
const token = useCookie('token');
const uid = ref('');
const userData = reactive({ data: {} });
watch(
  () => isLogin.value,
  () => {
    const userInfo = JSON.parse(localStorage.getItem('userInfo'));
    uid.value = userInfo.id;
    nextTick(() => {
      userData.data = parseJwt(token.value);
    });
  }
);
const localePath = useLocalePath();
const { t, setLocale, locale } = useI18n();
const isLocale = ref(locale.value);
const setLanguage = (lang) => {
  isLocale.value = lang;
  setLocale(lang);
};
nextTick(() => {
  if (token.value) {
    userData.data = parseJwt(token.value);
  }
});

const logOut = async () => {
  const result = await userStore.logOut(t);
  if (result.status.value === 'success') {
    clearAllCookie();
    navigateTo(localePath('/'));
    setTimeout(() => {
      location.reload();
    }, 500);
    message.success(t('signUp.logOut') + t('signUp.success'));
  }
};
const linkTo = async (link) => {
  mobileDialog.value = false;
  await navigateTo(localePath(link));
};
const mobileDialog = ref(false);
const toggleDialog = () => {
  mobileDialog.value = !mobileDialog.value;
};
const userDialog = ref(false);
const transferDialog = ref(false);
const transferMobileArea = ref(false);
const dialogOpen = (type) => {
  if (type === 'user') {
    transferDialog.value = false;
    userDialog.value = !userDialog.value;
  } else {
    transferDialog.value = !transferDialog.value;
    userDialog.value = false;
  }
};
onMounted(() => {
  if (isLogin.value !== 0) {
    const userInfo = JSON.parse(localStorage.getItem('userInfo'));
    uid.value = userInfo.id;
  }
});
</script>
<style>
.v-enter-active,
.v-leave-active {
  transition: opacity 0.3s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>
