<template>
  <div class="relative bg-white flex flex-col items-center">
    <h1 class="md:text-h2 text-h3 text-gray_800 pt-[120px] md:mb-[55px] mb-[35px]">{{ $t('signUp.loginRightNow') }}</h1>
    <div class="bg-white border border-grayBorder rounded-md p-[35px] md:w-[450px] xxs:w-[350px] w-[300px] md:mb-[175px] mb-[75px] z-10">
      <a-form :model="formState" @finish="submit" :rules="rules">
        <a-form-item name="account">
          <FormLabel label="signUp.memberAccount" haveIcon></FormLabel>
          <FormInput v-model="formState.account" :placeholder="$t('signUp.pleaseEnter') + $t('signUp.email')" />
        </a-form-item>
        <a-form-item name="password">
          <FormLabel label="signUp.memberPassword" haveIcon></FormLabel>
          <FormInput v-model="formState.password" :password="true" :placeholder="$t('signUp.pleaseEnter') + $t('signUp.password')" />
        </a-form-item>
        <div class="text-small text-gray_500 text-right underline -mt-4 hover:text-waterBlue">
          <NuxtLink :to="localePath('/forgotPassword')" class="cursor-pointer">{{ $t('signUp.forgotPassword') }}?</NuxtLink>
        </div>
        <div class="flex flex-col items-center">
          <Button variant="waterBlue" class="w-[120px] shadow-xl shadow-gray_700/20" :loading="isLoading">{{ $t('signUp.login') }}</Button>
          <p class="text-gray_500 text-small mt-6">
            {{ $t('signUp.firstComing') }} <NuxtLink :to="localePath('/signUp')" class="text-small text-waterBlue underline">{{ $t('signUp.signUpRightNow') }}</NuxtLink>
          </p>
        </div>
      </a-form>
    </div>
    <img src="/assets/img/Group16.png" class="absolute right-0 md:w-[640px] w-[400px]" />
    <img src="/assets/img/Group37.png" class="absolute left-0 top-0 md:w-[500px] w-[100px]" />
  </div>
  <Dialog v-model="dialogOpen">
    <div class="flex flex-col items-center">
      <p class="text-subTitle font-bold mt-[70px] mb-5">{{ $t('signUp.enterEmailPassCode') }}</p>
      <!-- <p class="text-normal text-grayText">{{ $t('signUp.passCodeAlreadySend') }}</p> -->
      <p class="text-normal text-grayText">{{ formatAccount(showEmail) }}</p>
      <Button variant="white" :disabled="buttonDisabled" class="mt-[50px]" :class="{ 'bg-gray_300 cursor-not-allowed text-gray_400 hover:bg-opacity-100 active:bg-opacity-100': buttonDisabled }" @click="reSend"
        >{{ $t('signUp.reSend') }} <span v-if="buttonDisabled">{{ countDown }} s</span></Button
      >
      <FormEmailValid :codes="codes" />
      <Button variant="waterBlue" :disabled="continueDisabled" class="w-[200px] mt-10 mb-4 text-[25px] rounded-xl" :class="{ 'bg-gray_300 cursor-not-allowed text-gray_400 hover:bg-opacity-100 active:bg-opacity-100': continueDisabled }" :loading="isLoading" @click="verifyContinue">{{ $t('signUp.continue') }}</Button>
    </div>
  </Dialog>
  <Dialog v-model="dialog2fa">
    <div class="flex flex-col items-center">
      <p class="text-subTitle font-bold mt-[70px] mb-5">{{ $t('signUp.enterEmailPassCode') }}</p>
      <!-- <p class="text-normal text-grayText">{{ $t('signUp.passCodeAlreadySend') }}</p> -->
      <p class="text-normal text-grayText">{{ formatAccount(showEmail) }}</p>
      <FormEmailValid :codes="codes" />
      <Button variant="waterBlue" :disabled="continueDisabled" :loading="isLoading" class="w-[200px] mt-10 mb-4 text-[25px] rounded-xl" :class="{ 'bg-gray_300 cursor-not-allowed text-gray_400 hover:bg-opacity-100 active:bg-opacity-100': continueDisabled }" @click="login2fa">{{ $t('signUp.continue') }}</Button>
    </div>
  </Dialog>
</template>

<script setup>
import useUserStore from '@/stores/user';
import formatAccount from '@/config/config';

const localePath = useLocalePath();
const userStore = useUserStore();
const route = useRoute();
const { t } = useI18n();
const isLoading = computed(() => userStore.loadingButton);
const dialogOpen = ref(false);
const dialog2fa = ref(false);
const codes = ref(['', '', '', '', '', '']);
const showEmail = ref('xxxxx@email.com');
watch(
  () => [codes.value],
  ([newVal]) => {
    const noEmpty = newVal.every((item) => item !== '');
    if (noEmpty) continueDisabled.value = false;
    else continueDisabled.value = true;
  },
  { deep: true }
);
const formState = ref({
  account: '',
  password: '',
});

const rules = {
  account: {
    require: true,
    validator: async (rule, value) => {
      if (value === '') {
        return Promise.reject(t('rules.dontEmpty'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'blur',
  },
  password: {
    required: true,
    validator: async (rule, value) => {
      if (value === '') {
        return Promise.reject(t('rules.dontEmpty'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'blur',
  },
};

const continueDisabled = ref(true);
const buttonDisabled = ref(false);
const countDown = ref(60);
const reSend = async () => {
  const data = {
    id: Number(localStorage.getItem('id')),
  };
  const result = await userStore.resendEmailVerify(data, t);
  if (result.status.value === 'success') {
    message.success(t('signUp.successStep1'));
  }
  buttonDisabled.value = true;
  countDownHandler();
};
let timer;
const countDownHandler = () => {
  timer = setInterval(() => {
    countDown.value--;
    if (countDown.value === 0) {
      clearInterval(timer);
      countDown.value = 60;
      buttonDisabled.value = false;
    }
  }, 1000);
};
const verifyContinue = async () => {
  userStore.loadingButton = true;
  const data = {
    id: Number(localStorage.getItem('id')),
    verificationCode: '',
  };
  data.verificationCode = codes.value.join('');
  const result = await userStore.emailVerify(data, t);
  userStore.loadingButton = false;
  if (result.code === 4022 || result.code === 4023) {
    cancelTimer();
    message.error(t('signUp.passCodeError'));
  } else if (result.status.value === 'success') {
    setTimeout(() => {
      location.reload();
    }, 500);
    message.success(t('signUp.successStep2'));
  }
};
const cancelTimer = () => {
  clearInterval(timer);
  countDown.value = 60;
  buttonDisabled.value = false;
};
const submit = async () => {
  userStore.loadingButton = true;
  const result = await userStore.login(formState.value, t);
  userStore.loadingButton = false;
  showEmail.value = formState.value.account;
  if (result?.data.value.accountNotVerified) {
    localStorage.setItem('id', result.data.value.id);
    dialogOpen.value = true;
    message.warn(t('signUp.notYetValid'));
  } else if (result?.data.value.login2FAType) {
    message.warn(t('signUp.warning2fa'));
    localStorage.setItem('id', result.data.value.id);
    localStorage.setItem('onePassKey', result.data.value.onePassKey);
    dialog2fa.value = true;
  } else if (result?.status.value === 'success') {
    localStorage.setItem('id', result.data.value.id);
    message.success(t('signUp.login') + t('signUp.success'));
    navigateTo(localePath('/MyAssets'));
  }
};
const login2fa = async () => {
  userStore.loadingButton = true;
  const data = {
    id: Number(localStorage.getItem('id')),
    onePassKey: localStorage.getItem('onePassKey'),
    verificationCode: '',
  };
  data.verificationCode = codes.value.join('');
  const result = await userStore.login2fa(data, t);
  userStore.loadingButton = false;
  if (result?.status.value === 'success') {
    message.success(t('signUp.login') + t('signUp.success'));
    if(route.query.isVerify){
      navigateTo(localePath(`/Members/${data.id}/verify`));
    }else{
      navigateTo(localePath('/MyAssets'));
    }
  }
};
</script>
