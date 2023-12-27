<template>
  <div class="relative bg-white flex flex-col items-center">
    <h1 class="md:text-h2 text-h3 text-gray_800 pt-[120px] md:mb-[55px] mb-[35px]">{{ $t('signUp.forgotPassword') }}</h1>
    <div class="bg-white border border-grayBorder rounded-md p-[35px] md:w-[450px] xxs:w-[350px] w-[300px] md:mb-[175px] mb-[75px] z-10">
      <template v-if="step === 1">
        <div class="bg-waterBlue px-9 py-6 text-center rounded-xl opacity-60 mb-10">
          <p class="text-white text-subText mb-2">{{ $t('signUp.notice3') }}</p>
          <p class="text-white text-subText">{{ $t('signUp.notice4') }}</p>
        </div>
        <a-form :model="formState" @finish="submit(1)" :rules="rules">
          <a-form-item name="account">
            <FormLabel label="signUp.memberAccount" haveIcon></FormLabel>
            <FormInput v-model="formState.account" :placeholder="$t('signUp.pleaseEnter') + $t('signUp.email')" />
          </a-form-item>
          <div class="flex flex-col items-center">
            <Button variant="waterBlue" class="w-[100px] shadow-xl shadow-gray_700/20">{{ $t('signUp.continue') }}</Button>
            <p class="text-gray_500 text-small mt-6">
              <NuxtLink :to="localePath('/login')" class="text-small underline">{{ $t('signUp.loginGoBack') }}</NuxtLink>
            </p>
          </div>
        </a-form>
      </template>
      <template v-if="step === 2">
        <a-form :model="passwordState" @finish="submit(2)" :rules="rules">
          <a-form-item name="password">
            <FormLabel label="signUp.password" haveIcon></FormLabel>
            <FormInput v-model="passwordState.password" :password="true" :placeholder="$t('signUp.pleaseEnter') + $t('signUp.password')" @change="checkValue" />
            <template v-if="passwordState.password !== ''">
              <div class="mt-2 ml-2">
                <p class="flex items-center" :class="`${haveUpperLower ? `text-waterBlue` : `text-gray_400`}`"><span class="inline-block w-[12px] h-[12px] rounded-full mr-2" :class="`${haveUpperLower ? `bg-waterBlue` : `bg-gray_400`}`"></span>{{ $t('signUp.passwordNotice1') }}</p>
                <p class="flex items-center" :class="`${enoughLength ? `text-waterBlue` : `text-gray_400`}`"><span class="inline-block w-[12px] h-[12px] rounded-full mr-2" :class="`${enoughLength ? `bg-waterBlue` : `bg-gray_400`}`"></span>{{ $t('signUp.passwordNotice2') }}</p>
              </div>
            </template>
          </a-form-item>
          <FormLabel label="signUp.confirmPassword" haveIcon class="md:mt-6 mt-3"></FormLabel>
          <FormInput v-model="passwordState.confirmPassword" :password="true" :placeholder="$t('signUp.pleaseEnter') + $t('signUp.confirmPassword')" />
          <div class="flex flex-col items-center mt-7">
            <Button variant="waterBlue" class="w-[100px] shadow-xl shadow-gray_700/20">{{ $t('signUp.continue') }}</Button>
            <p class="text-gray_500 text-small mt-6">
              <NuxtLink :to="localePath('/login')" class="text-small underline">{{ $t('signUp.loginGoBack') }}</NuxtLink>
            </p>
          </div>
        </a-form>
      </template>
    </div>
    <img src="/assets/img/Group16.png" class="absolute right-0 md:w-[640px] w-[400px]" />
    <img src="/assets/img/Group37.png" class="absolute left-0 top-0 md:w-[500px] w-[100px]" />
  </div>
  <Dialog v-model="dialogOpen">
    <div class="flex flex-col items-center">
      <p class="text-subTitle font-bold mt-[70px] mb-5">{{ $t('signUp.enterEmailPassCode') }}</p>
      <p class="text-normal text-grayText">{{ $t('signUp.passCodeAlreadySend') }}</p>
      <p class="text-normal text-grayText">{{ formatAccount(showEmail) }}</p>
      <FormEmailValid :codes="codes" />
      <Button variant="waterBlue" :disabled="continueDisabled" class="w-[200px] mt-10 mb-4 text-[25px] rounded-xl" :class="{ 'bg-gray_300 cursor-not-allowed text-gray_400 hover:bg-opacity-100 active:bg-opacity-100': continueDisabled }" @click="verifyContinue">{{ $t('signUp.continue') }}</Button>
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
const step = ref(1);

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
      const regex = /^(?=.*[a-z])(?=.*[A-Z]).{8,12}$/;
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
});
const passwordState = ref({
  password: '',
  confirmPassword: '',
  token: '',
});
const submit = async (step) => {
  if (step === 1) {
    const result = await userStore.forgotPassword(formState.value, t);
    if (result.status.value === 'success') {
      showEmail.value = formState.value.account;
      localStorage.setItem('id', result.data.value.id);
      dialogOpen.value = true;
      message.success(t('signUp.successStep1'));
    }
  } else {
    if (passwordState.value.password !== passwordState.value.confirmPassword) {
      message.error(t('signUp.pleaseConfirmPassword'));
      return;
    }
    const result = await userStore.resetPassword(passwordState.value, t);
    if (result.status.value === 'success') {
      message.success(t('signUp.passwordEditSuccess'));
      navigateTo(localePath('/login'));
    }
  }
};
const verifyContinue = async () => {
  const data = {
    account: formState.value.account,
    verificationCode: '',
  };
  data.verificationCode = codes.value.join('');
  const result = await userStore.verifyResetPassword(data);
  if (result.code === 4022 || result.code === 4023) {
    message.error(t('signUp.passCodeError'));
    dialogOpen.value = false;
    codes.value = ['', '', '', '', '', ''];
  } else if (result.status.value === 'success') {
    passwordState.value.token = result.data.value.token;
    message.success(t('signUp.success'));
    dialogOpen.value = false;
    step.value = 2;
  }
};
const haveUpperLower = ref(false);
const enoughLength = ref(false);
const checkValue = () => {
  const uppercaseRegex = /[A-Z]+/;
  const lowercaseRegex = /[a-z]+/;
  const haveUppercase = uppercaseRegex.test(passwordState.value.password);
  const haveLowercase = lowercaseRegex.test(passwordState.value.password);
  if (haveUppercase && haveLowercase) {
    haveUpperLower.value = true;
  } else {
    haveUpperLower.value = false;
  }
  if (passwordState.value.password.length >= 8 && passwordState.value.password.length < 17) {
    enoughLength.value = true;
  } else {
    enoughLength.value = false;
  }
};
</script>
<style>
.code::-webkit-outer-spin-button,
.code::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}
</style>
