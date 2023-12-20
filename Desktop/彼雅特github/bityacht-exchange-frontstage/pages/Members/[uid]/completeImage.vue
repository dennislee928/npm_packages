<template>
  <section class="relative 3xl:px-[300px] xl:px-[200px] lg:px-[100px] px-[30px] bg-white pb-[100px]">
    <div class="relative z-20">
      <h1 class="md:text-h2 text-h4 text-gray_800 text-center pt-[88px] md:pb-[60px] pb-[30px]">{{ $t('AuthVerify.completeImage') }}</h1>
      <a-divider class="md:hidden block border border-gray_300"></a-divider>
      <a-form :model="localFormState" @finish="submit" :rules="rules">
        <div class="flex flex-col">
          <a-form-item name="idImage" class="w-full">
            <FormLabel label="AuthVerify.idImage" haveIcon required></FormLabel>
            <FormUpload type="idImage" @upload="handleUpload">
              <FormInput v-model="localFormState.idImageName" :placeholder="$t('AuthVerify.pleaseUpload') + $t('AuthVerify.idImage')" />
              <FolderOpenOutlined class="text-gray_800 absolute top-1/2 right-2 -translate-y-1/2" />
            </FormUpload>
          </a-form-item>
          <a-form-item name="idBackImage" class="w-full">
            <FormLabel label="AuthVerify.idBackImage" haveIcon required></FormLabel>
            <FormUpload type="idBackImage" @upload="handleUpload">
              <FormInput v-model="localFormState.idBackImageName" :placeholder="$t('AuthVerify.pleaseUpload') + $t('AuthVerify.idBackImage')" />
              <FolderOpenOutlined class="text-gray_800 absolute top-1/2 right-2 -translate-y-1/2" />
            </FormUpload>
          </a-form-item>
        </div>
        <div class="flex flex-col">
          <a-form-item name="passportImage" class="w-full">
            <FormLabel label="AuthVerify.passportImage" haveIcon required></FormLabel>
            <FormUpload type="passportImage" @upload="handleUpload">
              <FormInput v-model="localFormState.passportImageName" :placeholder="$t('AuthVerify.pleaseUpload') + $t('AuthVerify.passportImage')" />
              <FolderOpenOutlined class="text-gray_800 absolute top-1/2 right-2 -translate-y-1/2" />
            </FormUpload>
          </a-form-item>
          <a-form-item name="idAndFaceImage" class="w-full">
            <FormLabel label="AuthVerify.idAndFaceImage" haveIcon required></FormLabel>
            <FormUpload type="idAndFaceImage" @upload="handleUpload">
              <FormInput v-model="localFormState.idAndFaceImageName" :placeholder="$t('AuthVerify.pleaseUpload') + $t('AuthVerify.idAndFaceImage')" />
              <FolderOpenOutlined class="text-gray_800 absolute top-1/2 right-2 -translate-y-1/2" />
            </FormUpload>
          </a-form-item>
        </div>
        <a-divider class="md:hidden block border border-gray_300"></a-divider>
        <div class="bg-[#f1f3f6] p-8 rounded-lg mt-8">
          <p class="text-normal font-bold">{{ $t('AuthVerify.notice') }}</p>
          <p class="text-normal mt-2">{{ $t('AuthVerify.noticeText') }}</p>
        </div>
        <div class="md:ml-10 ml-5 my-5">
          <a-checkbox v-model:checked="checkSubmit">{{ $t('AuthVerify.agree') }}</a-checkbox>
        </div>
        <a-divider></a-divider>
        <div class="flex justify-center">
          <Button variant="waterBlue" html-type="submit" class="w-[100px] shadow-xl shadow-gray_700/20 mr-5" :disabled="!checkSubmit" :class="`${checkSubmit ? `` : `cursor-not-allowed`}`" :loading="isLoading">{{ $t('AuthVerify.submit') }}</Button>
          <Button variant="gray" class="w-[100px] shadow-xl shadow-gray_700/20" @click="back">{{ $t('AuthVerify.back') }}</Button>
        </div>
        <div class="flex justify-center mt-8 md:flex-row flex-col items-center">
          <p class="text-normal text-waterBlue">{{ $t('AuthVerify.noticeText2') }}</p>
          <p class="text-normal text-waterBlue">{{ $t('AuthVerify.noticeText3') }}</p>
        </div>
      </a-form>
    </div>
    <img src="/assets/img/Group16.png" class="absolute right-0 top-0 md:w-[600px] w-[400px] z-10" />
    <img src="/assets/img/Group37.png" class="absolute left-0 top-0 md:w-[800px] w-[200px] z-10" />
  </section>
</template>

<script setup>
import useUserStore from '@/stores/user';
import { formatPhone } from '@/config/config';

const localePath = useLocalePath();
const back = async () => {
  await navigateTo(localePath(`/Members/${userData.value.id}/level`));
};
const { t } = useI18n();
useHead({
  title: t('title.memberCenter_IDVerify'),
  meta: [{ name: 'description', content: '' }],
});
const userStore = useUserStore();
const isLoading = computed(() => userStore.loadingButton);
const country = ref('');
watch(
  () => country.value,
  (val) => {
    if (val === 'local') {
      localFormState.countriesCode = 'TWN';
    }
  }
);
const localFormState = reactive({
  phone: '',
  idImage: '', //file
  idImageName: '',
  idBackImage: '', //file
  idBackImageName: '',
  passportImage: '', //file
  passportImageName: '',
  idAndFaceImage: '', //file
  idAndFaceImageName: '',
});
const addressError = ref(false);
const rules = {
  idImage: {
    require: true,
    validator: async (rule, value) => {
      if (value === '') {
        return Promise.reject(t('rules.dontEmpty'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'submit',
  },
  idBackImage: {
    require: true,
    validator: async (rule, value) => {
      if (value === '') {
        return Promise.reject(t('rules.dontEmpty'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'submit',
  },
  passportImage: {
    require: true,
    validator: async (rule, value) => {
      if (value === '') {
        return Promise.reject(t('rules.dontEmpty'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'submit',
  },
  idAndFaceImage: {
    require: true,
    validator: async (rule, value) => {
      if (value === '') {
        return Promise.reject(t('rules.dontEmpty'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'submit',
  },
};
// const validateRules = (name, status, errorMsgs) => {
//   if (!status) {
//     message.error(t(`AuthVerify.${name}`) + errorMsgs[0]);
//   }
// };

const checkSubmit = ref(false);
const dialogOpen = ref(false);
const codes = ref(['', '', '', '', '', '']);
const continueDisabled = ref(true);
watch(
  () => [codes.value],
  ([newVal]) => {
    const noEmpty = newVal.every((item) => item !== '');
    if (noEmpty) continueDisabled.value = false;
    else continueDisabled.value = true;
  },
  { deep: true }
);
const submit = async () => {
  userStore.loadingButton = true;
  const idvResult = await userStore.patchImage(localFormState, t);
  userStore.loadingButton = false;
  if (idvResult.status.value === 'success') {
    message.success(t('AuthVerify.IDVerifySent'));
    await navigateTo({
      path: localePath('/status'),
      query: {
        title: 'AuthVerify.IDVerifySent',
        imgSource: 'dataSuccess',
      },
    });
  }
};

const handleUpload = (file, type, fileName) => {
  if (type === 'idImage') {
    localFormState.idImage = file;
    localFormState.idImageName = fileName;
  } else if (type === 'idBackImage') {
    localFormState.idBackImage = file;
    localFormState.idBackImageName = fileName;
  } else if (type === 'passportImage') {
    localFormState.passportImage = file;
    localFormState.passportImageName = fileName;
  } else if (type === 'idAndFaceImage') {
    localFormState.idAndFaceImage = file;
    localFormState.idAndFaceImageName = fileName;
  }
};
const userData = ref();
onMounted(async () => {
  userData.value = JSON.parse(localStorage.getItem('userInfo'));
  const result = await userStore.getIdvOption(t);
});
</script>
<style>
.code::-webkit-outer-spin-button,
.code::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}
</style>
