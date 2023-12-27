<template>
  <NuxtLayout name="members">
    <section class="relative bg-white 3xl:px-[300px] xl:px-[200px] lg:px-[100px] px-[30px]">
      <div class="relative z-30 flex justify-between md_lg:flex-row flex-col pb-[75px]">
        <div class="relative md_lg:w-[49.5%] w-full flex flex-col border border-gray_800 rounded-xl lg:px-[45px] xs:px-[35px] px-[18px] py-10">
          <h1 class="flex items-start">
            <img src="/icon/hexagon.svg" alt="icon" class="mr-2 w-[24px]" />
            <div class="flex flex-col">
              <p class="text-black md:text-subTitle text-normal font-bold">
                {{ $t('memberCenter.loginPassword') }}
              </p>
              <p class="text-black md:text-normal text-subText mt-1">
                {{ $t('memberCenter.forLoginPassword') }}
              </p>
            </div>
          </h1>
          <div class="items-start md_lg:hidden flex mt-2">
            <img src="/icon/light.svg" class="w-[30px]" />
            <div class="flex flex-col mt-2 ml-3">
              <p class="text-small text-yellow">{{ $t('memberCenter.passwordContent1') }}</p>
              <p class="text-small text-grayBorder">{{ $t('memberCenter.passwordContent2') }}</p>
            </div>
          </div>
          <Transition>
            <template v-if="!isMobile">
              <a-form :model="formState" :rules="rules" @finish="submit">
                <div class="flex flex-col justify-center items-center bg-waterBlue bg-opacity-80 rounded-xl py-3 px-10 my-6">
                  <p class="text-white text-small mb-2">{{ $t('memberCenter.notice1') }}</p>
                  <p class="text-white text-small">{{ $t('memberCenter.notice2') }}</p>
                </div>
                <div class="flex justify-between items-center">
                  <div class="md_lg:w-[auto] w-full">
                    <a-form-item name="oldPassword">
                      <FormLabel label="memberCenter.nowPassword" haveIcon></FormLabel>
                      <FormInput v-model="formState.oldPassword" :password="true" :placeholder="$t('memberCenter.nowPassword')" class="3xl:w-[250px] 2xl:w-[210px] xl:w-[180px] lg:w-[150px] md_lg:w-[145px] w-full" />
                    </a-form-item>
                    <NuxtLink :to="localePath('/forgotPassword')">
                      <p class="text-small text-pink2 -mt-5">{{ $t('memberCenter.forgotPassword') }}</p>
                    </NuxtLink>
                  </div>
                  <div class="items-start md_lg:flex hidden">
                    <img src="/icon/light.svg" class="w-[30px]" />
                    <div class="flex flex-col mt-2">
                      <p class="text-small text-yellow">{{ $t('memberCenter.passwordContent1') }}</p>
                      <p class="text-small text-grayBorder">{{ $t('memberCenter.passwordContent2') }}</p>
                    </div>
                  </div>
                </div>
                <div class="flex justify-between md_lg:flex-nowrap flex-wrap">
                  <div class="md_lg:w-auto w-full md_lg:mt-0 mt-5">
                    <a-form-item name="newPassword">
                      <FormLabel label="memberCenter.newPassword" haveIcon></FormLabel>
                      <FormInput v-model="formState.newPassword" :password="true" :placeholder="$t('memberCenter.enter') + $t('memberCenter.nowPassword')" class="3xl:w-[250px] 2xl:w-[210px] xl:w-[180px] lg:w-[150px] md_lg:w-[145px] w-full" />
                    </a-form-item>
                    <p class="text-small text-pink2 -mt-5">{{ $t('memberCenter.passwordText') }}</p>
                  </div>
                  <div class="md_lg:w-auto w-full md_lg:mt-0 mt-5">
                    <a-form-item name="confirmPassword">
                      <FormLabel label="memberCenter.confirmPassword" haveIcon></FormLabel>
                      <FormInput v-model="formState.confirmPassword" :password="true" :placeholder="$t('memberCenter.again') + $t('memberCenter.enter') + $t('memberCenter.newPassword')" class="3xl:w-[250px] 2xl:w-[210px] xl:w-[180px] lg:w-[150px] md_lg:w-[145px] w-full" />
                    </a-form-item>
                  </div>
                </div>
                <div class="flex justify-center md_lg:mt-8 mt-5">
                  <Button variant="waterBlue" class="w-[85px]">{{ $t('memberCenter.submit') }}</Button>
                </div>
              </a-form>
            </template>
          </Transition>

          <div class="md:hidden block cursor-pointer text-subText text-waterBlue underline absolute top-[20px] sm:right-[50px] right-[25px]" @click="isMobile = !isMobile">{{ !isMobile ? $t('memberCenter.cancel') : $t('memberCenter.edit') }}</div>
        </div>
        <div class="relative md_lg:w-[49.5%] w-full flex flex-col border border-gray_800 rounded-xl lg:px-[45px] xs:px-[35px] px-[18px] py-10 md_lg:mt-0 mt-10">
          <h1 class="flex items-start">
            <img src="/icon/hexagon.svg" alt="icon" class="mr-2 w-[24px]" />
            <div class="flex flex-col">
              <p class="text-black md:text-subTitle text-normal font-bold">
                {{ $t('memberCenter.doubleValid') }}
                <span v-if="googleAuthenticator" class="border border-green text-green py-1 px-2 rounded-xl sm:text-subText text-small ml-2">{{ $t('memberCenter.alreadyStart') }}</span>
              </p>
              <p class="text-black md:text-normal text-subText mt-1">
                {{ $t('memberCenter.forLoginbalala') }}
              </p>
            </div>
          </h1>
          <div class="flex md_lg:flex-row flex-col items-center md_lg:mt-[50px] mt-5">
            <!-- <img src="/assets/img/doubleValid.png" class="3xl:w-[320px] xl:w-[200px] md_lg:w-[150px] w-[220px] md_lg:block hidden" /> -->
            <img v-if="background64" :src="`data:image/png;base64,${background64}`" class="md_lg:w-[200px] sm:w-[300px] w-auto md_lg:h-[200px] sm:h-[300px] h-auto" />
            <div class="flex items-start">
              <div class="flex flex-col mt-2 ml-3 w-full">
                <p class="text-subText">{{ $t('memberCenter.googleAuth') }}</p>
                <p class="text-small text-pink2 mt-4">{{ $t('memberCenter.noticeText1') }}</p>
                <p class="text-small text-pink2 mt-2">{{ $t('memberCenter.noticeText2') }}</p>
                <p class="text-small text-pink2 mt-2">{{ $t('memberCenter.noticeText3') }}</p>
                <p class="text-small text-pink2 mt-2">{{ $t('memberCenter.noticeText4') }}<a class="text-pink2 hover:text-waterBlue transition-all duration-300" href="mailto:support@bityacht.io"> support@bityacht.io </a>{{ $t('memberCenter.noticeText5') }}</p>
              </div>
            </div>
          </div>
          <Transition>
            <template v-if="!isMobile2">
              <div class="flex flex-col items-center mt-2">
                <div v-if="secret" class="w-full bg-gray_500 rounded-xl bg-opacity-30 flex flex-col items-center py-4">
                  <p class="text-subText mb-2">{{ $t('memberCenter.key') }}</p>
                  <p class="text-waterBlueOld font-bold flex items-center sm:text-normal text-small">{{ secret }} <CopyOutlined class="text-gray_500 text-normal ml-1" @click="copyText(secret, t)" /></p>
                </div>
                <template v-if="!googleAuthenticator">
                  <Button v-if="step === 1" variant="outline" :class="`w-[110px] mx-auto 3xl:mt-10 mt-6 `" @click="post2fa('open')">{{ $t('memberCenter.start') }}</Button>
                  <Button v-if="step === 2" variant="outline" :class="`w-[110px] mx-auto 3xl:mt-10 mt-6 `" @click="postGoogle()">{{ $t('memberCenter.start') }}</Button>
                </template>
                <template v-else>
                  <Button variant="outline" :class="`w-[110px] mx-auto 3xl:mt-10 mt-6 `" @click="post2fa('close')">{{ $t('memberCenter.close') }}</Button>
                </template>
              </div>
            </template>
          </Transition>
          <div class="md:hidden block cursor-pointer text-subText text-waterBlue underline absolute top-[20px] sm:right-[50px] right-[25px]" @click="isMobile2 = !isMobile2">{{ !isMobile2 ? $t('memberCenter.cancel') : $t('memberCenter.edit') }}</div>
        </div>
      </div>
    </section>
    <Dialog v-model="dialogOpen">
      <div class="relative flex flex-col items-center pb-5">
        <p class="text-h4 font-bold mt-[50px] mb-5">{{ $t('memberCenter.editPasswordSuccess') }}</p>
        <img src="/assets/img/changePassword.png" />
        <div class="flex flex-col justify-center items-center bg-waterBlue opacity-60 rounded-xl py-3 px-10 my-6">
          <p class="text-white text-small mb-2">{{ $t('memberCenter.notice1') }}</p>
          <p class="text-white text-small">{{ $t('memberCenter.notice2') }}</p>
        </div>
        <Button variant="waterBlue" class="w-[150px] mt-5 mb-5 text-[25px] rounded-xl" @click="close">{{ $t('memberCenter.close') }}</Button>
      </div>
    </Dialog>
    <Dialog v-model="dialog2fa">
      <div class="relative flex flex-col items-center">
        <p class="text-subTitle font-bold mt-[70px] mb-5">{{ $t('signUp.enterEmailPassCode') }}</p>
        <p class="text-normal text-grayText">{{ $t('signUp.passCodeAlreadySend') }}</p>
        <p class="text-normal text-grayText">{{ formatAccount(showEmail) }}</p>
        <FormEmailValid :codes="codes2fa" />
        <Button variant="waterBlue" :disabled="continueDisabled" :loading="isLoading" class="w-[200px] mt-10 mb-4 text-[25px] rounded-xl" :class="{ 'bg-gray_300 cursor-not-allowed text-gray_400 hover:bg-opacity-100 active:bg-opacity-100': continueDisabled }" @click="get2fa">{{ $t('signUp.continue') }}</Button>
        <close-circle-outlined class="absolute top-0 -right-4 text-[20px]" @click="dialog2fa = !dialog2fa" />
      </div>
    </Dialog>
    <Dialog v-model="dialogGoogleAuth">
      <div class="relative flex flex-col items-center">
        <p class="text-subTitle font-bold mt-[70px] mb-5">{{ $t('signUp.enterGoogleAuth') }}</p>
        <FormEmailValid :codes="codesGa" />
        <Button variant="waterBlue" :disabled="continueDisabled" :loading="isLoading" class="w-[200px] mt-10 mb-4 text-[25px] rounded-xl" :class="{ 'bg-gray_300 cursor-not-allowed text-gray_400 hover:bg-opacity-100 active:bg-opacity-100': continueDisabled }" @click="patch2fa">{{ $t('signUp.continue') }}</Button>
        <close-circle-outlined class="absolute top-0 -right-4 text-[20px]" @click="dialogGoogleAuth = !dialogGoogleAuth" />
      </div>
    </Dialog>
  </NuxtLayout>
</template>
<script setup>
import useUserStore from '@/stores/user';
import formatAccount from '@/config/config';
import { copyText } from '@/config/config';
import { useLocalStorage } from '@vueuse/core';

const localePath = useLocalePath();
const isLoading = computed(() => userStore.loadingButton);
const userStore = useUserStore();
const { t } = useI18n();
const dialogOpen = ref(false);
const dialog2fa = ref(false);
const dialogGoogleAuth = ref(false);
const userData = ref(null);
const showEmail = ref('xxxxx@gmail.com');
const codes2fa = ref(['', '', '', '', '', '']);
const codesGa = ref(['', '', '', '', '', '']);
const continueDisabled = ref(true);
watch(
  () => [codes2fa.value],
  ([newVal]) => {
    const noEmpty = newVal.every((item) => item !== '');
    if (noEmpty) continueDisabled.value = false;
    else continueDisabled.value = true;
  },
  { deep: true }
);
watch(
  () => [codesGa.value],
  ([newVal]) => {
    const noEmpty = newVal.every((item) => item !== '');
    if (noEmpty) continueDisabled.value = false;
    else continueDisabled.value = true;
  },
  { deep: true }
);

const formState = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
});
const rules = {
  oldPassword: {
    required: true,
    validator: async (rule, value) => {
      if (value === '') {
        return Promise.reject(t('rules.dontEmpty'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'blur',
  },
  newPassword: {
    required: true,
    validator: async (rule, value) => {
      const regex = /^(?=.*[a-z])(?=.*[A-Z]).{8,12}$/;
      if (!regex.test(value)) {
        return Promise.reject(t('rules.passwordError'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'blur',
  },
};
const close = () => {
  dialogOpen.value = false;
  dialogGoogleAuth.value = false;
  setTimeout(() => {
    location.reload();
  }, 1000);
};

const submit = async () => {
  if (formState.value.newPassword !== formState.value.confirmPassword) {
    message.error(t('signUp.pleaseConfirmPassword'));
    return;
  }
  const result = await userStore.editPassword(formState.value, t);
  if (result.status.value === 'success') {
    message.success(t('signUp.passwordEditSuccess'));
    dialogOpen.value = true;
  }
};
const actionStatus = ref('');
const post2fa = async (action) => {
  await userStore.post2fa(t);
  dialog2fa.value = true;
  actionStatus.value = action;
};
const postGoogle = async () => {
  dialogGoogleAuth.value = true;
};
const background64 = ref(null);
const secret = ref(null);
const token = ref(null);
const step = ref(1);
const get2fa = async () => {
  userStore.loadingButton = true;
  const verificationCode = codes2fa.value.join('');
  const result = await userStore.get2fa(verificationCode, t);
  userStore.loadingButton = false;
  if (result.status.value === 'success') {
    dialog2fa.value = false;
    token.value = result.data.value.token;
    if (actionStatus.value === 'open') {
      step.value++;
      background64.value = result.data.value.qrCode;
      secret.value = result.data.value.secret;
    } else {
      dialogGoogleAuth.value = true;
    }
  }
};
const patch2fa = async () => {
  userStore.loadingButton = true;
  const verificationCode = codesGa.value.join('');
  const data = {
    verificationCode: verificationCode,
    token: token.value,
  };
  const result = await userStore.patch2fa(data, t);
  userStore.loadingButton = false;
  if (result.status.value === 'success') {
    message.success(t('signUp.success'));
    close();
  }
};
const googleAuthenticator = ref(false);
const getUserInfo = async () => {
  const result = await userStore.getUserInfo(t, true);
  googleAuthenticator.value = result.data.value.googleAuthenticator;
};
const isMobile = ref(false);
const isMobile2 = ref(false);
onMounted(async () => {
  await getUserInfo();
  userData.value = JSON.parse(localStorage.getItem('userInfo'));
  showEmail.value = userData.value.account;
  const screen = window.innerWidth;
  if (screen < 768) {
    isMobile.value = true;
    isMobile2.value = true;
  }
});
</script>
