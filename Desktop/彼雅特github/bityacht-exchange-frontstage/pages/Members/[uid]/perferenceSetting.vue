<template>
  <NuxtLayout name="members">
    <section class="relative bg-white 3xl:px-[300px] xl:px-[200px] lg:px-[100px] px-[30px] pb-10">
      <div class="z-30 relative border border-gray_800 rounded-xl xs:p-[55px] p-[25px]">
        <h1 class="flex items-start mb-5">
          <img src="/icon/hexagon.svg" alt="icon" class="mr-2 w-[18px]" />
          <div class="flex flex-col">
            <p class="text-black text-subTitle font-bold">
              {{ $t('memberCenter.receiptSetting') }}
            </p>
            <p class="text-black text-normal mt-2">{{ $t('memberCenter.mobileBarcodeBind') }}</p>
            <p class="md:block hidden text-grayBorder text-subText mt-2">{{ $t('memberCenter.loveEarth') }}</p>
            <p v-if="haveMobileBarcode && !goBind" class="text-black text-normal mt-5">{{ $t('memberCenter.alreadyMobileBind') }} {{ formState.mobileBarcode }}</p>
          </div>
        </h1>
        <Transition>
          <template v-if="goBind">
            <a-form :model="formState" @finish="submit" :rules="rules" class="ml-6">
              <a-form-item name="mobileBarcode">
                <FormLabel label="memberCenter.barcodeSetting"></FormLabel>
                <div class="flex md:flex-row flex-col">
                  <FormInput v-model="formState.mobileBarcode" :placeholder="$t('signUp.pleaseEnter') + $t('memberCenter.mobileBarcode')" class="3xl:w-[300px] 2xl:w-[250px] xl:w-[220px] lg:w-[180px] md_lg:w-[15px] w-full" />
                  <Button variant="waterBlue" class="md:w-auto w-[80px] mt-2 :md:ml-4 mx-auto"> {{ $t('memberCenter.submit') }}</Button>
                </div>
              </a-form-item>
            </a-form>
          </template>
        </Transition>
        <div class="flex items-center flex-wrap">
          <img src="/icon/light.svg" class="w-[30px]" />
          <p class="text-yellow text-subText ml-2">{{ $t('memberCenter.howtoWatchBarcode1') }}</p>
          <p class="text-grayBorder text-small md:mt-0 mt-2">{{ $t('memberCenter.howtoWatchBarcode2') }}</p>
        </div>
        <p class="cursor-pointer text-subText text-waterBlue underline absolute top-[60px] right-10" @click="goBind = !goBind">{{ goBind ? $t('memberCenter.cancel') : $t('memberCenter.bind') }}</p>
      </div>
    </section>
  </NuxtLayout>
</template>
<script setup>
import useUserStore from '@/stores/user';

const userStore = useUserStore();
const { t } = useI18n();
const formState = ref({
  mobileBarcode: '',
});
const rules = {
  mobileBarcode: {
    required: true,
    validator: async (rule, value) => {
      const regex = /^\/[0-9A-Z\+\-\.]{7}$/;
      if (value !== '') {
        if (!regex.test(value)) {
          return Promise.reject(t('rules.barCodeError'));
        }
        return Promise.resolve(t('rules.actionSuccess'));
      }
    },
    trigger: 'blur',
  },
};
const goBind = ref(false);
const haveMobileBarcode = ref(false);
const submit = async () => {
  const result = await userStore.mobileBarcode(formState.value, t);
  if (result.status.value === 'success') {
    message.success(t('signUp.success'));
    getUserInfo();
    goBind.value = false;
  }
};
const getUserInfo = async () => {
  const result = await userStore.getUserInfo(t, true);
  if (result.data.value.mobileBarcode !== '') {
    haveMobileBarcode.value = true;
    formState.value.mobileBarcode = result.data.value.mobileBarcode;
  }
  // console.log('result :>> ', result.data.value);
};
onMounted(() => {
  getUserInfo();
});
</script>
