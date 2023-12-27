<template>
  <div class="relative bg-white flex flex-col items-center">
    <h1 class="md:text-h2 text-h3 text-gray_800 pt-[120px] md:mb-[55px] mb-[35px]">{{ $t('signUp.start') }}</h1>
    <div class="bg-white border border-grayBorder rounded-md p-[35px] md:w-[450px] xxs:w-[350px] w-[300px] md:mb-[175px] mb-[75px] z-10">
      <a-form :model="formState" @finish="submit" :rules="rules">
        <a-form-item name="account">
          <FormLabel label="signUp.email" haveIcon></FormLabel>
          <FormInput v-model="formState.account" :placeholder="$t('signUp.pleaseEnter') + $t('signUp.email')" />
        </a-form-item>
        <a-form-item name="password">
          <FormLabel label="signUp.password" haveIcon></FormLabel>
          <FormInput v-model="formState.password" :password="true" :placeholder="$t('signUp.pleaseEnter') + $t('signUp.password')" @change="checkValue" />
          <template v-if="formState.password !== ''">
            <div class="mt-2 ml-2">
              <p class="flex items-center" :class="`${haveUpperLower ? `text-waterBlue` : `text-gray_400`}`"><span class="inline-block w-[12px] h-[12px] rounded-full mr-2" :class="`${haveUpperLower ? `bg-waterBlue` : `bg-gray_400`}`"></span>包含至少一大寫及小寫英文字母</p>
              <p class="flex items-center" :class="`${enoughLength ? `text-waterBlue` : `text-gray_400`}`"><span class="inline-block w-[12px] h-[12px] rounded-full mr-2" :class="`${enoughLength ? `bg-waterBlue` : `bg-gray_400`}`"></span>長度須為 8 - 16 位數</p>
            </div>
          </template>
        </a-form-item>
        <FormLabel label="signUp.inviteCode" haveIcon class="md:mt-6 mt-3"></FormLabel>
        <FormInput v-model="formState.inviteCode" :placeholder="$t('signUp.pleaseEnter') + $t('signUp.inviteCode')" :disabled="inviteCodeDisabled" />
        <div class="flex flex-col items-center">
          <p class="xs:text-subText text-small mt-8 md:w-auto w-full">{{ $t('signUp.notice1') }}</p>
          <p class="xs:text-subText text-small mb-5 md:w-auto w-full">{{ $t('signUp.notice2') }}</p>
          <Button variant="waterBlue" class="w-[100px] shadow-xl shadow-gray_700/20" :loading="isLoading">{{ $t('signUp.signUp') }}</Button>
          <p class="text-gray_500 text-small mt-6">
            {{ $t('signUp.alreadyAccount') }} <NuxtLink :to="localePath('/login')" class="text-small text-waterBlue underline">{{ $t('signUp.loginRightNow') }}</NuxtLink>
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
      <p class="text-normal text-grayText">{{ $t('signUp.passCodeAlreadySend') }}</p>
      <p class="text-normal text-grayText">{{ formatAccount(showEmail) }}</p>
      <Button variant="white" :disabled="buttonDisabled" class="mt-[50px]" :class="{ 'bg-gray_300 cursor-not-allowed text-gray_400 hover:bg-opacity-100 active:bg-opacity-100': buttonDisabled }" @click="reSend"
        >{{ $t('signUp.reSend') }} <span v-if="buttonDisabled">{{ countDown }} s</span></Button
      >
      <FormEmailValid :codes="codes" />
      <Button v-if="continueDisabled" variant="disabled" :disabled="continueDisabled" class="w-[200px] mt-10 mb-4 text-[25px] rounded-xl" @click="verifyContinue">{{ $t('signUp.continue') }}</Button>
      <Button v-else variant="waterBlue" class="w-[200px] mt-10 mb-4 text-[25px] rounded-xl" :loading="isLoading" @click="verifyContinue">{{ $t('signUp.continue') }}</Button>
    </div>
  </Dialog>
</template>

<script setup>
import useUserStore from '@/stores/user';
import formatAccount from '@/config/config';

const localePath = useLocalePath();
const userStore = useUserStore();
const { t } = useI18n();
const dialogOpen = ref(false);
const codes = ref(['', '', '', '', '', '']);
const isLoading = computed(() => userStore.loadingButton);

watch(
  () => [codes.value],
  ([newVal]) => {
    const noEmpty = newVal.every((item) => item !== '');
    if (noEmpty) continueDisabled.value = false;
    else continueDisabled.value = true;
  },
  { deep: true }
);

const continueDisabled = ref(true);
const showEmail = ref('xxxxx@email.com');

const rules = {
  account: {
    require: true,
    validator: async (rule, value) => {
      const emailPattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
      if (!emailPattern.test(value)) {
        return Promise.reject(t('rules.pleaseEnterEmail'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'blur',
  },
  password: {
    required: true,
    validator: async (rule, value) => {
      const regex = /^(?=.*[a-z])(?=.*[A-Z]).{8,16}$/;
      if (!regex.test(value)) {
        return Promise.reject(t('rules.passwordError'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'blur',
  },
};
const formState = ref({
  account: '',
  password: '',
  inviteCode: '',
});
const submit = async () => {
  userStore.loadingButton = true;
  const result = await userStore.register(formState.value, t);
  userStore.loadingButton = false;
  if (result?.status.value === 'success') {
    showEmail.value = formState.value.account;
    localStorage.setItem('id', result.data.value.id);
    dialogOpen.value = true;
    reSend('runSetTime');
    message.success(t('signUp.successStep1'));
  }
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
  if (result?.code === 4022 || result?.code === 4023) {
    cancelTimer();
    message.error(t('signUp.passCodeError'));
  } else if (result?.status?.value === 'success') {
    message.success(t('signUp.successStep2'));
    navigateTo(localePath('/login'));
  }
};
const buttonDisabled = ref(false);
const countDown = ref(60);
const reSend = async (type) => {
  if (type === 'runSetTime') {
  } else {
    const data = {
      id: Number(localStorage.getItem('id')),
    };
    const result = await userStore.resendEmailVerify(data, t);
    if (result.status.value === 'success') {
      message.success(t('signUp.successStep1'));
    }
  }
  buttonDisabled.value = true;
  countDownHandler();
};
const cancelTimer = () => {
  clearInterval(timer);
  countDown.value = 60;
  buttonDisabled.value = false;
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

const haveUpperLower = ref(false);
const enoughLength = ref(false);
const checkValue = () => {
  const uppercaseRegex = /[A-Z]+/;
  const lowercaseRegex = /[a-z]+/;
  const haveUppercase = uppercaseRegex.test(formState.value.password);
  const haveLowercase = lowercaseRegex.test(formState.value.password);
  if (haveUppercase && haveLowercase) {
    haveUpperLower.value = true;
  } else {
    haveUpperLower.value = false;
  }
  if (formState.value.password.length >= 8 && formState.value.password.length < 17) {
    enoughLength.value = true;
  } else {
    enoughLength.value = false;
  }
};
const inviteCodeDisabled = ref(false);
onMounted(() => {
  // console.log('useRoute() :>> ', useRoute());
  const route = useRoute();
  if (route.query.inviteCode) {
    formState.value.inviteCode = route.query.inviteCode;
    inviteCodeDisabled.value = true;
  }
});
</script>
<style>
.code::-webkit-outer-spin-button,
.code::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}
</style>
