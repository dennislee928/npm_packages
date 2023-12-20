<template>
  <section class="relative 3xl:px-[300px] xl:px-[200px] lg:px-[100px] px-[30px] bg-white md:pb-[140px] pb-[70px]">
    <h1 class="relative flex items-center justify-center z-20 xl:text-h2 md:text-h4 text-[26px] text-gray_800 text-center pt-[88px] md:pb-[60px] pb-[20px]"><left-circle-outlined v-if="step === 1" class="mr-4 text-h5" @click="back" />{{ $t('myAssets.cryptocurrency') }}{{ $t('myAssets.withdraw') }}</h1>
    <div class="relative z-30 flex items-center justify-between xl:mx-10 lg:mx-5 sm:mx-2 mx-0">
      <div class="flex items-center md:flex-row flex-col justify-center">
        <h1 :class="`hexagon flex items-center justify-center md:w-[75px] w-[50px] md:h-[65px] h-[40px] xl:text-h2 md:text-h5 text-[26px] text-white md:mb-0 mb-2 ${step === 1 ? 'bg-waterBlue' : 'bg-gray_500'}`">1</h1>
        <p :class="`xl:text-h5 md:text-[24px] xs:text-normal text-subText  ml-2 ${step === 1 ? 'text-waterBlue' : 'text-gray_500'}`">{{ $t('myAssets.writeAddress') }}</p>
      </div>
      <span class="relative md:w-[90px] w-[35px] border-b border-gray_500 after:content-[''] after:block after:w-3 after:rotate-45 after:border-b after:border-gray_500 after:absolute after:-top-1 after:-right-0.5"></span>
      <div class="flex items-center md:flex-row flex-col justify-center">
        <h1 :class="`hexagon flex items-center justify-center md:w-[75px] w-[50px] md:h-[65px] h-[40px] xl:text-h2 md:text-h5 text-[26px] text-white md:mb-0 mb-2 ${step === 2 ? 'bg-waterBlue' : 'bg-gray_500'}`">2</h1>
        <p :class="`xl:text-h5 md:text-[24px] xs:text-normal text-subText  ml-2 ${step === 2 ? 'text-waterBlue' : 'text-gray_500'}`">{{ $t('myAssets.enterCount') }}</p>
      </div>
      <span class="relative md:w-[90px] w-[35px] border-b border-gray_500 after:content-[''] after:block after:w-3 after:rotate-45 after:border-b after:border-gray_500 after:absolute after:-top-1 after:-right-0.5"></span>
      <div class="flex items-center md:flex-row flex-col justify-center">
        <h1 :class="`hexagon flex items-center justify-center md:w-[75px] w-[50px] md:h-[65px] h-[40px] xl:text-h2 md:text-h5 text-[26px] text-white md:mb-0 mb-2 ${step === 3 ? 'bg-waterBlue' : 'bg-gray_500'}`">3</h1>
        <p :class="`xl:text-h5 md:text-[24px] xs:text-normal text-subText  ml-2 ${step === 3 ? 'text-waterBlue' : 'text-gray_500'}`">{{ $t('myAssets.resultConfirm') }}</p>
      </div>
    </div>
    <div v-if="step === 1" class="relative z-30 md:w-[580px] w-full mx-auto md:border border-0 border-gray_800 rounded-xl md:mt-[70px] mt-8 py-6 md:px-[60px] px-2">
      <div class="mt-4 flex flex-col">
        <div class="flex flex-col mb-6">
          <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('myAssets.choiceCoinType') }}</p>
          <a-select v-model:value="formState.currencyType" :placeholder="$t('AuthVerify.pleaseSelect') + $t('myAssets.coinType')" @change="selectClear" :disabled="$route.query.coinType ? true : false">
            <a-select-option v-for="item of coinTypeList" :key="item" :value="item">{{ item }}</a-select-option>
          </a-select>
        </div>
        <div class="flex flex-col mb-6">
          <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('myAssets.choiceMainInternet') }}</p>
          <a-select v-model:value="formState.mainnet" :placeholder="$t('AuthVerify.pleaseSelect') + $t('myAssets.mainInternet')">
            <a-select-option v-for="item of mainInternetList" :key="item" :value="item.name">{{ item.name }}</a-select-option>
          </a-select>
        </div>
        <a-form :model="formData">
          <a-form-item name="address">
            <div class="flex flex-col mb-1">
              <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('myAssets.getAddress') }}</p>
              <FormInput v-model="formData.address" :placeholder="$t('signUp.pleaseEnter') + $t('myAssets.getAddress')" class="relative" readonly><img v-if="whiteListStatus === 'notNull'" src="/icon/notebook.svg" class="w-[20px] absolute top-[50%] -translate-y-1/2 right-2 mt-1 cursor-pointer" @click="dialogWhiteList = !dialogWhiteList" /></FormInput>
            </div>
          </a-form-item>
        </a-form>
        <Transition>
          <div v-if="whiteListStatus" class="flex bg-shallow_gray rounded-xl p-2 items-center bg-opacity-50">
            <exclamation-circle-outlined class="text-waterBlueOld mr-1 text-[20px]" />
            <span class="text-subText">{{ $t('myAssets.addressNotice') }}</span>
            <span class="ml-4 text-waterBlueOld text-subText cursor-pointer" @click="whiteListStatus === 'null' ? (dialogNullWhiteList = !dialogNullWhiteList) : (dialogWhiteList = !dialogWhiteList)">{{ whiteListStatus === 'null' ? $t('myAssets.goAddAddress') : $t('myAssets.goSetAddress') }}</span>
          </div>
        </Transition>
      </div>
      <div class="flex justify-center mt-10">
        <a-button :disabled="stepOneDisabled" class="w-[82px] bg-waterBlue text-white rounded-2xl hover:!text-gray_300" @click="step++">{{ $t('myAssets.continue') }}</a-button>
      </div>
    </div>
    <div v-if="step === 2" class="relative z-30 flex md:flex-row flex-col justify-between xl:mx-10 lg:mx-5 sm:mx-2 mx-0 mt-[60px]">
      <div class="md:w-[31%] w-full md:mb-0 mb-5">
        <div class="md:flex hidden bg-gray_800 rounded-xl justify-center items-center text-[20px] text-white py-5">{{ $t('myAssets.withdraw') }}{{ $t('myAssets.information') }}</div>
        <div class="bg-gray_800 rounded-xl flex flex-col text-[20px] text-white py-5 2xl:px-10 lg:px-5 md:px-3 px-10 mt-4">
          <div class="flex items-start">
            <img src="/icon/play.svg" class="mt-1 mr-3" />
            <div class="flex flex-col">
              <span class="text-normal text-white mb-2">{{ $t('myAssets.coinType') }}</span>
              <span class="text-normal text-white">{{ formState.currencyType }}</span>
            </div>
          </div>
          <div class="flex items-start mt-5">
            <img src="/icon/play.svg" class="mt-1 mr-3" />
            <div class="flex flex-col">
              <span class="text-normal text-white mb-2">{{ $t('myAssets.mainInternet') }}</span>
              <span class="text-normal text-white">{{ formState.mainnet }}</span>
            </div>
          </div>
          <div class="flex items-start mt-5 md:mb-[80px] mb-0">
            <img src="/icon/play.svg" class="mt-1 mr-3" />
            <div class="flex flex-col flex-wrap">
              <span class="text-normal text-white mb-2">{{ $t('myAssets.getAddress') }}</span>
              <span class="md_lg:text-normal md:text-small text-normal text-white break-all">{{ formData.address }}</span>
            </div>
          </div>
        </div>
      </div>
      <div class="md:w-[67.5%] w-full py-[40px] md:px-[77px] sm:px-10 px-6 border border-gray_800 rounded-xl">
        <div class="flex justify-between items-center">
          <div class="flex items-center">
            <img src="/icon/grayMoney.svg" />
            <span class="text-normal font-bold ml-2">{{ $t('myAssets.canUseCount') }}</span>
          </div>
          <div class="text-normal">{{ canUseRemaingingAmount }} {{ formState.currencyType }}</div>
        </div>
        <div class="flex justify-between items-center mt-4">
          <div class="flex items-center xxs:w-auto w-[150px]">
            <img src="/icon/Vector.svg" />
            <span class="text-normal font-bold ml-2">{{ $t('myAssets.remainingCountOfDay') }}</span>
          </div>
          <div class="text-normal">{{ amountOfDay }} U</div>
        </div>
        <div class="flex flex-col mt-8">
          <span class="flex items-center text-normal font-bold mb-2"><img src="/icon/hexagon.svg" class="mr-2" />{{ $t('myAssets.withdraw') }}{{ $t('myAssets.count') }}</span>
          <FormInput v-model="formData.amount" type="number" :placeholder="$t('myAssets.pleaseEnterAmount')"></FormInput>
          <div :class="`flex mt-1 ${error ? 'justify-between' : 'justify-end'}`">
            <div v-if="error" class="text-red text-subText">
              {{ $t(`rules.${errorText}`) }}
              <span v-if="minError">{{ minNotional }}</span>
              <span v-if="maxError">{{ maxNotional }}</span>
            </div>
            <!-- <div class="text-small text-gray_500">= TWD</div> -->
          </div>
        </div>
        <a-divider></a-divider>
        <div class="flex md:justify-end justify-between items-center flex-wrap">
          <span class="mr-10 text-normal font-bold">{{ $t('myAssets.withdraw') }}{{ $t('myAssets.actual') }}{{ $t('myAssets.count') }}</span>
          <span class="text-[20px] font-bold">{{ actualAmount }} {{ formState.currencyType }}</span>
        </div>
        <div class="flex justify-end items-center mt-2">
          <span class="mr-2 text-small text-gray_500">{{ $t('myAssets.redemptionFee') }} : {{ redemptionFee }} {{ formState.currencyType }} </span>
          <a-tooltip>
            <template #title
              ><div class="flex flex-col justify-center items-center p-2">
                <span class="text-small text-white">{{ $t('myAssets.importNoticeText4') }} </span>
                <span class="text-small text-white mt-3">{{ $t('myAssets.importNoticeText5') }} </span>
              </div></template
            >
            <exclamation-circle-outlined class="text-normal text-gray_500" />
          </a-tooltip>
        </div>
        <div class="md:hidden flex flex-col justify-end items-center mt-5">
          <span class="mr-2 text-small text-gray_500">{{ $t('myAssets.importNoticeText4') }} </span>
          <span class="mr-2 text-small text-gray_500">{{ $t('myAssets.importNoticeText5') }} </span>
        </div>
      </div>
    </div>
    <div v-if="step === 2" class="relative z-30 flex md:justify-end justify-center xl:mx-10 lg:mx-5 sm:mx-2 mx-0">
      <a-button :class="`w-[75px] my-3 rounded-3xl text-small mr-3 bg-gray_500 text-white hover:!text-white`" @click="step--">{{ $t('myAssets.back') }}</a-button>
      <a-button :disabled="error" :loading="isLoading" :class="`w-[75px] my-3 rounded-3xl text-small bg-waterBlue text-white hover:!text-white`" @click="withdraw">{{ $t('myAssets.continue') }}</a-button>
    </div>
  </section>
  <img src="/assets/img/Group16.png" class="absolute right-0 top-0 md:w-[600px] w-[400px] z-10" />
  <img src="/assets/img/Group37.png" class="absolute left-0 top-0 md:w-[800px] w-[200px] z-10" />
  <template v-if="dialog2fa">
    <Dialog v-model="dialog2fa">
      <div class="relative flex flex-col items-center">
        <p class="text-subTitle font-bold mt-[70px] mb-5">{{ $t('signUp.enterPassCode') }}</p>
        <p class="text-normal text-grayText">{{ $t('signUp.passCodeAlreadySend') }}</p>
        <p class="text-normal text-grayText">{{ formatAccount(showEmail) }}</p>
        <p class="mt-5 text-normal font-bold">{{ $t('signUp.emailPassCode') }}</p>
        <FormEmailValid :codes="codesEmail" />
        <p class="mt-5 text-normal font-bold">{{ $t('signUp.smsPassCode') }}</p>
        <FormEmailValid :codes="codesSms" />
        <template v-if="googleAuthenticator">
          <p class="mt-5 text-normal font-bold">{{ $t('signUp.gaPassCode') }}</p>
          <FormEmailValid :codes="codesGa" />
        </template>
        <Button variant="waterBlue" :disabled="twoFADisabled" :loading="isLoading" class="w-[200px] mt-10 mb-4 text-[25px] rounded-xl" :class="{ 'bg-gray_300 cursor-not-allowed text-gray_400 hover:bg-opacity-100 active:bg-opacity-100': twoFADisabled }" @click="withdraw2fa">{{ $t('signUp.continue') }}</Button>
        <close-circle-outlined class="absolute top-0 -right-4 text-[20px]" @click="close" />
      </div>
    </Dialog>
  </template>
  <Dialog v-model="dialogNullWhiteList">
    <div class="flex flex-col justify-center items-center py-5 md:w-[550px] w-auto">
      <p class="font-bold text-subTitle mb-5">{{ $t('myAssets.choiceWithdrawAddress') }}</p>
      <img src="/assets/img/watch.png" />
      <p class="text-shallow_gray my-5 text-subText w-[300px] text-center">{{ $t('myAssets.choiceWithdrawAddressText') }}</p>
      <Button variant="waterBlue" class="w-full" @click="goAdd">{{ $t('myAssets.addWithdrawAddress') }}</Button>
    </div>
  </Dialog>
  <Dialog v-model="dialogAddWhiteList">
    <Transition>
      <div class="flex flex-col justify-center items-center py-5 md:w-[550px] w-[300px]">
        <p class="font-bold text-subTitle mb-5">{{ $t('myAssets.addWithdrawAddress') }}</p>
        <a-form :model="whiteListData" class="w-full">
          <a-form-item name="address">
            <div class="relative flex flex-col mb-1">
              <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('myAssets.getAddress') }}</p>
              <FormInput v-model="whiteListData.address" :placeholder="$t('signUp.pleaseEnter') + $t('myAssets.getAddress')"></FormInput>
            </div>
          </a-form-item>
          <a-form-item name="remark">
            <div class="flex flex-col mb-1">
              <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('myAssets.remark') }}</p>
              <FormInput v-model="whiteListData.extra.memo" :placeholder="$t('signUp.pleaseEnter') + $t('myAssets.remark')"></FormInput>
            </div>
          </a-form-item>
        </a-form>
        <Button variant="waterBlue" class="w-full" @click="addWhiteList">{{ $t('myAssets.addWithdrawAddress') }}</Button>
        <close-circle-outlined class="absolute top-3 right-3 text-[20px]" @click="dialogAddWhiteList = !dialogAddWhiteList" />
      </div>
    </Transition>
  </Dialog>
  <Dialog v-model="dialogWhiteList">
    <div class="flex flex-col justify-center items-center py-5 md:w-[550px] w-[320px]">
      <p class="font-bold text-subTitle mb-5">{{ $t('myAssets.choiceWithdrawAddress') }}</p>
      <div v-for="(item, index) of sliceWhiteList" class="flex justify-center items-center mb-5">
        <div class="md:w-8 w-6 md:h-8 h-6 bg-gray_400 rounded-full flex justify-center items-center mr-2 text-white md:text-subText text-small" :class="{ 'bg-waterBlueOld': selectIndex === index }">{{ index + 1 }}</div>
        <a-tooltip placement="topRight">
          <template #title>
            <p v-if="item.extra.memo !== ''" class="md:text-subText text-small">{{ $t('myAssets.remark') }}: {{ item.extra.memo }}</p>
            <p class="md:text-subText text-small">{{ item.address }}</p>
          </template>
          <div class="border border-gray-400 p-1.5 md:text-subText text-small cursor-pointer md:w-[350px] w-[290px] whitespace-nowrap text-ellipsis overflow-hidden" :class="{ 'border-waterBlueOld text-waterBlueOld': selectIndex === index }" @click="select(item, index)">{{ item.address }}</div>
        </a-tooltip>
        <close-circle-filled class="text-red md:text-[30px] text-[24px] ml-2 cursor-pointer" @click="deleteCheck(item.id)" />
      </div>
      <a-pagination v-model:current="tableConfig.current" v-model:pageSize="tableConfig.pageSize" simple :total="tableConfig.total" hideOnSinglePage class="my-2" @change="changePage" />
      <div class="flex mt-4">
        <Button variant="waterBlue" class="w-full mr-4" @click="confirmSelectAddress">{{ $t('myAssets.choiceWithdrawAddress') }}</Button>
        <Button variant="darkBlue" class="w-full" @click="addWithdrawAddress">{{ $t('myAssets.addNewWithdrawAddress') }}</Button>
      </div>

      <close-circle-outlined class="absolute top-3 right-3 text-[20px]" @click="dialogWhiteList = !dialogWhiteList" />
    </div>
  </Dialog>
  <Dialog v-model="deleteCheckDialog">
    <div class="flex flex-col items-center">
      <p class="lg:text-[24px] text-[18px] mt-5">{{ $t('memberCenter.deleteConfirm') }}</p>
      <div>
        <Button variant="gray" class="w-[120px] mt-10 mb-4 mr-4 lg:text-[16px] text-[14px] rounded-xl" @click="deleteCheckDialog = !deleteCheckDialog">{{ $t('memberCenter.cancel') }}</Button>
        <Button variant="waterBlue" class="w-[120px] mt-10 mb-4 lg:text-[16px] text-[14px] rounded-xl" :loading="isLoading" @click="deleteWhiteList">{{ $t('signUp.continue') }}</Button>
      </div>
    </div>
  </Dialog>
</template>
<script setup>
import useAssetsStore from '@/stores/assets';
import useUserStore from '@/stores/user';
import formatAccount from '@/config/config';
import { Decimal } from 'decimal.js';

const localePath = useLocalePath();
const back = async () => {
  await navigateTo(localePath('/MyAssets'));
};
const { t } = useI18n();
useHead({
  title: t('title.myAssets_cryptocurrencyGet'),
  meta: [{ name: 'description', content: '' }],
});
const assetsStore = useAssetsStore();
const userStore = useUserStore();
const user = JSON.parse(localStorage.getItem('userInfo'));
const dialog2fa = ref(false);
const dialogNullWhiteList = ref(false);
const dialogAddWhiteList = ref(false);
const dialogWhiteList = ref(false);
const deleteCheckDialog = ref(false);
const deleteID = ref(null);
const deleteCheck = (id) => {
  deleteCheckDialog.value = true;
  deleteID.value = id;
};
const whiteListData = ref({
  address: '',
  extra: {
    memo: '',
  },
  mainnet: '',
});
const whiteList = ref(null);
const addWhiteList = async () => {
  whiteListData.value.mainnet = formData.value.mainnet;
  const result = await assetsStore.postWhiteList(whiteListData.value, t);
  if (result.status.value === 'success') {
    message.success(t('signUp.success'));
    await getWhiteList();
    dialogAddWhiteList.value = false;
    dialogWhiteList.value = true;
  }
};
const whiteListStatus = ref('');
const tableConfig = ref({
  position: ['bottomCenter'],
  current: 1,
  pageSize: 5,
  total: 0,
});
const getWhiteList = async () => {
  const result = await assetsStore.getWhiteList(formData.value.mainnet, t);
  if (result.status.value === 'success') {
    if (result.data.value.data.length === 0) {
      whiteListStatus.value = 'null';
      dialogWhiteList.value = false;
      formData.value.address = '';
    } else {
      whiteListStatus.value = 'notNull';
      whiteList.value = result.data.value.data;
      tableConfig.value.total = result.data.value.data.length;
      if (result.data.value.data.length < 6) {
        tableConfig.value.current = 1;
      }
    }
  }
};
const sliceWhiteList = computed(() => {
  if (!whiteList.value) return;
  const start = (tableConfig.value.current - 1) * tableConfig.value.pageSize;
  const end = start + tableConfig.value.pageSize;
  return whiteList.value.slice(start, end);
});
const changePage = () => {
  selectIndex.value = null;
};
const goAdd = () => {
  dialogNullWhiteList.value = false;
  dialogAddWhiteList.value = true;
  whiteListData.value.address = '';
  whiteListData.value.extra.memo = '';
};
const selectIndex = ref(null);
const selectAddress = ref('');
const selectAddressID = ref(null);
const select = (item, index) => {
  selectIndex.value = index;
  selectAddress.value = item.address;
  selectAddressID.value = item.id;
};
const confirmSelectAddress = () => {
  formData.value.address = selectAddress.value;
  formData.value.whitelistID = selectAddressID.value;
  dialogWhiteList.value = false;
};
const addWithdrawAddress = () => {
  whiteListData.value.address = '';
  whiteListData.value.extra.memo = '';
  dialogWhiteList.value = false;
  dialogAddWhiteList.value = true;
};
const deleteWhiteList = async () => {
  assetsStore.loadingButton = true;
  const result = await assetsStore.deleteWhiteList(deleteID.value, t);
  assetsStore.loadingButton = false;
  if (result.status.value === 'success') {
    message.success(t('signUp.success'));
    await getWhiteList();
    deleteCheckDialog.value = false;
    selectIndex.value = '';
  }
};
const codesEmail = ref(['', '', '', '', '', '']);
const codesSms = ref(['', '', '', '', '', '']);
const codesGa = ref(['', '', '', '', '', '']);
const showEmail = ref(user.account);
const continueDisabledEmail = ref(true);
const continueDisabledSms = ref(true);
const continueDisabledGa = ref(true);
const isLoading = computed(() => assetsStore.loadingButton);
watch(
  () => [codesEmail.value],
  ([newVal]) => {
    const noEmpty = newVal.every((item) => item !== '');
    if (noEmpty) continueDisabledEmail.value = false;
    else continueDisabledEmail.value = true;
  },
  { deep: true }
);
watch(
  () => [codesSms.value],
  ([newVal]) => {
    const noEmpty = newVal.every((item) => item !== '');
    if (noEmpty) continueDisabledSms.value = false;
    else continueDisabledSms.value = true;
  },
  { deep: true }
);
watch(
  () => [codesGa.value],
  ([newVal]) => {
    const noEmpty = newVal.every((item) => item !== '');
    if (noEmpty) continueDisabledGa.value = false;
    else continueDisabledGa.value = true;
  },
  { deep: true }
);
const twoFADisabled = computed(() => {
  if (googleAuthenticator.value) {
    if (!continueDisabledSms.value && !continueDisabledSms.value && !continueDisabledGa.value) return false;
    else return true;
  } else {
    if (!continueDisabledSms.value && !continueDisabledSms.value) return false;
    else return true;
  }
});
const step = ref(1);
const formState = ref({
  currencyType: undefined,
  mainnet: undefined,
});
const formData = ref({
  currencyType: undefined,
  mainnet: undefined,
  address: '',
  amount: '',
  whitelistID: '',
});
const actualAmount = ref(0);
const redemptionFee = ref(0);
const minNotional = ref(0);
const maxNotional = ref(0);
const amountOfDay = ref(0);
watch(
  () => formState.value.mainnet,
  async (newVal) => {
    if (newVal === undefined) return;
    if (formState.value.currencyType === 'BTC') {
      formData.value.currencyType = 1;
    } else if (formState.value.currencyType === 'ETH') {
      formData.value.currencyType = 2;
    } else if (formState.value.currencyType === 'USDC') {
      formData.value.currencyType = 3;
    } else {
      formData.value.currencyType = 4;
    }
    if (newVal === 'Bitcoin') {
      formData.value.mainnet = 1;
    } else if (newVal === 'Ethereum') {
      formData.value.mainnet = 2;
    } else if (newVal === 'Ethereum(ERC20)') {
      formData.value.mainnet = 3;
    } else {
      formData.value.mainnet = 4;
    }
    await getWhiteList();
    const result = await assetsStore.withdrawOption(formData.value.currencyType, formData.value.mainnet, t);
    if (result.status.value === 'success') {
      redemptionFee.value = parseFloat(result.data.value.withdrawFee);
      minNotional.value = parseFloat(result.data.value.withdrawMin);
      maxNotional.value = parseFloat(result.data.value.withdrawMax);
      amountOfDay.value = result.data.value.levelLimit.perDay.max - parseFloat(result.data.value.accWithdrawInDay);
    }
  }
);
const orignalList = ref({
  mainnets: {
    BTC: {
      BTC: {
        name: 'Bitcoin',
        value: 1,
      },
    },
    ETH: {
      ETH: {
        name: 'Ethereum',
        value: 2,
      },
    },
    USDC: {
      ETH: {
        name: 'Ethereum(ERC20)',
        value: 3,
      },
    },
    USDT: {
      ETH: {
        name: 'Ethereum(ERC20)',
        value: 3,
      },
      TRX: {
        name: 'Tron(TRC20)',
        value: 4,
      },
    },
  },
});
// const checkAddress = ref(false);
// const rules = {
//   address: {
//     require: true,
//     validator: async (rule, value) => {
//       const regexBTC = /\b[13][a-km-zA-HJ-NP-Z1-9]{25,34}\b/;
//       if (formState.value.currencyType === 'BTC') {

//         if (!regexBTC.test(value)) {

//           checkAddress.value = false;
//           return Promise.reject(t('rules.addressError'));
//         }
//         checkAddress.value = true;
//         return Promise.resolve(t('rules.actionSuccess'));
//       } else if (formState.value.currencyType === 'ETH') {
//       } else if (formState.value.currencyType === 'USDC') {
//       } else {
//       }
//     },
//     trigger: 'change',
//   },
// };
const stepOneDisabled = computed(() => {
  if (formData.value.currencyType === undefined || formData.value.mainnet === undefined || formData.value.address === '') return true;
  else return false;
});
const coinTypeList = ref([]);
const mainInternetList = ref([]);

const selectClear = () => {
  formState.value.mainnet = undefined;
};
watch(
  () => formState.value.currencyType,
  (newVal, oldVal) => {
    mainInternetList.value = orignalList.value.mainnets[newVal];
    if (oldVal !== newVal) {
      formData.value.mainnet = undefined;
    }
    if (userInfo.value) {
      userInfo.value.forEach((item) => {
        if (item.currenciesSymbol === newVal) {
          canUseRemaingingAmount.value = parseFloat(item.freeAmount);
        }
      });
    }
  }
);
const error = ref(false);
const minError = ref(false);
const maxError = ref(false);
const errorText = ref('');
// watch(
//   () => formData.value.amount,
//   (newVal) => {
//     if (newVal === '' || newVal === 0) {
//       actualAmount.value = 0;
//       return;
//     }
//     actualAmount.value = formData.value.amount - redemptionFee.value;
//   }
// );
watch(
  () => formData.value.amount,
  (newVal) => {
    if (newVal.split('.').length == 2 && newVal.split('.')[1].length > 8) {
      error.value = true;
      minError.value = false;
      maxError.value = false;
      errorText.value = 'pleaseNotOver8Digits';
      return false;
    } else if (newVal === '' || newVal === 0) {
      error.value = true;
      minError.value = false;
      maxError.value = false;
      actualAmount.value = 0;
      errorText.value = 'dontEmpty';
      return false;
    } else if (newVal <= 0) {
      error.value = true;
      minError.value = false;
      maxError.value = false;
      errorText.value = 'noZeroAndNegativeNumber';
      return false;
    } else if (minNotional.value >= 0 && newVal < minNotional.value) {
      error.value = true;
      minError.value = true;
      maxError.value = false;
      errorText.value = 'payAmountMinLimit';
      return false;
    } else if (maxNotional.value >= 0 && newVal > maxNotional.value) {
      error.value = true;
      minError.value = false;
      maxError.value = true;
      errorText.value = 'payAmountMaxLimit';
      return false;
    } else if (newVal > canUseRemaingingAmount.value) {
      error.value = true;
      minError.value = false;
      maxError.value = false;
      errorText.value = 'canUseAmountNotEnough';
      return false;
    }
    actualAmount.value = Decimal(formData.value.amount).sub(Decimal(redemptionFee.value));
    minError.value = false;
    error.value = false;
    return true;
  }
);
const userInfo = ref(undefined);
const canUseRemaingingAmount = ref(0);
const getAssets = async () => {
  const result = await assetsStore.getAssets(t);
  userInfo.value = result.data.value.assetInfos;
};
const getOptions = async () => {
  // const result = await assetsStore.getSpotOptions();
  // orignalList.value = result.data.value.mainnets;
  for (const coinType in orignalList.value.mainnets) {
    coinTypeList.value.push(coinType);
  }
};
const twoFAType = ref(0);
const withdraw = async () => {
  assetsStore.loadingButton = true;
  const result = await assetsStore.withdraw(formData.value, t, lastChangePasswordAt.value);
  assetsStore.loadingButton = false;
  if (result.status.value === 'success') {
    if (result.data.value.onePassKey) {
      message.warn(t('signUp.warning2fa'));
      localStorage.setItem('onePassKey', result.data.value.onePassKey);
      twoFAType.value = result.data.value.twoFAType;
      dialog2fa.value = true;
    } else {
      message.success(t('signUp.success'));
      await navigateTo({
        path: localePath('/status'),
        query: {
          title: 'myAssets.withdrawSuccess',
          imgSource: 'getSuccess',
          type: 'withdraw',
          amount: formData.value.amount,
          currencyType: formState.value.currencyType,
        },
      });
    }
  } else {
    await navigateTo({
      path: localePath('/status'),
      query: {
        title: 'myAssets.withdrawFail',
        imgSource: 'getFail',
      },
    });
  }
};
const close = async () => {
  dialog2fa.value = false;
  codesEmail.value = ['', '', '', '', '', ''];
  codesSms.value = ['', '', '', '', '', ''];
  codesGa.value = ['', '', '', '', '', ''];
};
const withdraw2fa = async () => {
  assetsStore.loadingButton = true;
  const data = {
    onePassKey: localStorage.getItem('onePassKey'),
    emailVerificationCode: '',
    smsVerificationCode: '',
  };

  if (twoFAType.value === 7) {
    data['gaVerificationCode'] = codesGa.value.join('');
  }
  data.emailVerificationCode = codesEmail.value.join('');
  data.smsVerificationCode = codesSms.value.join('');
  const result = await assetsStore.twofaWithdraw(data, t);
  assetsStore.loadingButton = false;
  if (result.status.value === 'success') {
    message.success(t('signUp.success'));
    await navigateTo({
      path: localePath('/status'),
      query: {
        title: 'myAssets.withdrawSuccess',
        imgSource: 'getSuccess',
        type: 'withdraw',
        amount: formData.value.amount,
        currencyType: formState.value.currencyType,
      },
    });
  } else {
    await navigateTo({
      path: localePath('/status'),
      query: {
        title: 'myAssets.withdrawFail',
        imgSource: 'getFail',
      },
    });
  }
};
const lastChangePasswordAt = ref('');
const googleAuthenticator = ref(null);
const getUserInfo = async () => {
  const result = await userStore.getUserInfo(t, true);
  lastChangePasswordAt.value = result.data.value.lastChangePasswordAt;
  googleAuthenticator.value = result.data.value.googleAuthenticator;
};
onMounted(async () => {
  await getUserInfo();
  await getOptions();
  await getAssets();
  const route = useRoute();
  formState.value.currencyType = route.query.coinType;
});
</script>
<style>
.hexagon {
  clip-path: polygon(75% 0%, 100% 50%, 75% 100%, 25% 100%, 0% 50%, 25% 0%);
}
.ant-select-selector {
  height: 40px !important;
  padding: 5px 10px !important;
}
input[type='number']::-webkit-outer-spin-button,
input[type='number']::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}
.v-enter-active,
.v-leave-active {
  transition: opacity 0.3s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>
