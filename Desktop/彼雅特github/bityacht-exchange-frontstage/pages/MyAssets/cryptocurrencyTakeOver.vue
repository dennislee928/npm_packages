<template>
  <section class="relative 3xl:px-[300px] xl:px-[200px] lg:px-[100px] px-[30px] bg-white md:pb-[140px] pb-[70px]">
    <h1 class="relative flex items-center justify-center xl:text-h2 md:text-h4 text-[26px] text-gray_800 text-center pt-[88px] z-30"><left-circle-outlined class="mr-4 text-h5" @click="back" />{{ $t('myAssets.cryptocurrency') }}{{ $t('myAssets.takeOver') }}</h1>
    <div class="relative z-30 md:w-[580px] w-full mx-auto flex flex-col justify-center md:border border-0 border-gray_800 rounded-xl md:mt-[70px] mt-8 py-6 md:px-[60px] px-2">
      <template v-if="step === 1">
        <img src="/assets/img/getSuccess.png" class="w-[450px] mx-auto" />
      </template>
      <template v-if="step === 2 && !isDeploying">
        <!-- <div :class="`bg-[url('${background64}.png')]`">test</div> -->
        <img :src="`data:image/png;base64,${background64}`" />
      </template>
      <div class="mt-6 flex flex-col">
        <Transition>
          <div v-if="step === 2" class="flex flex-col mb-6">
            <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('myAssets.inputCoinAddress') }}</p>
            <div class="flex justify-center items-center bg-grayBg text-black text-center p-3 rounded-lg break-all">
              {{ !isDeploying ? coinAddress : $t(`${coinAddress}`) }}<span v-if="!isDeploying" class="cursor-pointer ml-3 text-gray_500 flex" @click="copyText(coinAddress, t)"><CopyOutlined class="text-normal" /></span>
            </div>
          </div>
        </Transition>
        <div class="flex flex-col mb-6">
          <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('myAssets.choiceCoinType') }}</p>
          <!-- <FormAntdSelect v-model="formState.selectCoinType" :placeholder="$t('AuthVerify.pleaseSelect') + $t('myAssets.coinType')">
            <a-select-option v-for="item of coinTypeList" :key="item" :value="item">{{ item }}</a-select-option>
          </FormAntdSelect> -->
          <a-select v-model:value="formState.currencyType" :placeholder="$t('AuthVerify.pleaseSelect') + $t('myAssets.coinType')" @change="selectClear" :disabled="$route.query.coinType ? true : false">
            <a-select-option v-for="(item, index) of coinTypeList" :key="item" :value="item">{{ item }}</a-select-option>
          </a-select>
        </div>
        <div class="flex flex-col mb-6">
          <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('myAssets.choiceMainInternet') }}</p>
          <a-select v-model:value="formState.mainnet" :placeholder="$t('AuthVerify.pleaseSelect') + $t('myAssets.mainInternet')">
            <a-select-option v-for="item of mainInternetList" :key="item" :value="item">{{ item }}</a-select-option>
          </a-select>
          <!-- <FormAntdSelect v-model="formState.selectMainInternet" :placeholder="$t('AuthVerify.pleaseSelect') + $t('myAssets.mainInternet')" :allowClear="true">
            <a-select-option v-for="item of mainInternetList" :key="item" :value="item">{{ item }}</a-select-option>
          </FormAntdSelect> -->
        </div>
        <Transition>
          <div>
            <div v-if="step === 1" class="flex flex-col mb-6">
              <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('myAssets.inputCoinAddress') }}</p>
              <div class="bg-gray_300 rounded-lg flex justify-center items-center py-[50px] mt-2">
                <p class="text-grayText text-subTitle flex md:flex-row flex-col items-center justify-center">
                  {{ $t('myAssets.addressNotEnable') }} <a-button :disabled="!formState.currencyType || !formState.mainnet" :class="`w-[125px] my-3 md:ml-5 ml-0 rounded-3xl text-normal bg-waterBlue text-white hover:!text-white`" @click="produceAddress">{{ $t('myAssets.enableAddress') }}</a-button>
                </p>
              </div>
            </div>
            <div v-if="step === 2" class="bg-gray_800 rounded-lg mt-2 p-6 text-white">
              <div class="flex items-center mb-5">
                <img src="/icon/whiteLight.svg" class="w-8" />
                <p class="ml-2">{{ $t('myAssets.importNotice') }}</p>
              </div>
              <div class="flex items-center mb-1">
                <caret-right-outlined />
                <p class="ml-2">{{ $t('myAssets.importNoticeText1') }}{{ formState.currencyType }}</p>
              </div>
              <div class="flex items-center mb-1">
                <caret-right-outlined />
                <p class="ml-2">{{ $t('myAssets.importNoticeText2') }}{{ formState.currencyType }}({{ formState.mainnet }})</p>
              </div>
              <div class="flex items-center mb-1">
                <caret-right-outlined />
                <p class="ml-2">{{ formState.currencyType }}({{ formState.mainnet }}){{ $t('myAssets.importNoticeText3') }}</p>
              </div>
            </div>
          </div>
        </Transition>
      </div>
    </div>
  </section>
  <img src="/assets/img/Group16.png" class="absolute right-0 top-0 md:w-[600px] w-[400px] z-10" />
  <img src="/assets/img/Group37.png" class="absolute left-0 top-0 md:w-[800px] w-[200px] z-10" />
</template>

<script setup>
import useAssetsStore from '@/stores/assets';
import { copyText } from '@/config/config';

const localePath = useLocalePath();
const { t } = useI18n();
useHead({
  title: t('title.myAssets_cryptocurrencyTakeOver'),
  meta: [{ name: 'description', content: '' }],
});
const back = async () => {
  await navigateTo(localePath('/MyAssets'));
};
const assetsStore = useAssetsStore();
const step = ref(1);
const formState = ref({
  currencyType: undefined,
  mainnet: undefined,
});
const formData = ref({
  currencyType: undefined,
  mainnet: undefined,
});
const orignalList = ref({
  mainnets: {
    BTC: {
      BTC: 'Bitcoin',
    },
    ETH: {
      ETH: 'Ethereum',
    },
    USDC: {
      ETH: 'Ethereum (ERC20)',
    },
    USDT: {
      ETH: 'Ethereum (ERC20)',
      TRX: 'Tron (TRC20)',
    },
  },
});
const coinTypeList = ref([]);
const mainInternetList = ref([]);

const getOptions = async () => {
  // const result = await assetsStore.getSpotOptions();
  // orignalList.value = result.data.value.mainnets;
  for (const coinType in orignalList.value.mainnets) {
    coinTypeList.value.push(coinType);
  }
};
const selectClear = () => {
  formState.value.mainnet = undefined;
};
watch(
  () => formState.value.currencyType,
  (newVal) => {
    mainInternetList.value = orignalList.value.mainnets[newVal];
  }
);

const coinAddress = ref('');
const background64 = ref('');
const isDeploying = ref(false);
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
    } else if (newVal === 'Ethereum (ERC20)') {
      formData.value.mainnet = 3;
    } else {
      formData.value.mainnet = 4;
    }
    localStorage.setItem(`${formState.value.currencyType}-mainnet`, JSON.stringify({ mainnet: formData.value.mainnet, currencyType: formState.value.currencyType, mainnetName: newVal }));
    getAddress();
  }
);
const produceAddress = async () => {
  const result = await assetsStore.produceAddress(formData.value, t);
  if (result.status.value === 'success') {
    message.success(t('myAssets.enableAddress') + t('myAssets.success'));
    step.value = 2;
    getAddress();
  }
};
const getAddress = async () => {
  const result = await assetsStore.getAddress(formData.value.currencyType, formData.value.mainnet);
  if (result.status.value === 'success') {
    step.value = 2;
    if (result.data.value.code === 5038) {
      isDeploying.value = true;
      coinAddress.value = 'myAssets.addressDeploying';
    } else {
      isDeploying.value = false;
      coinAddress.value = result.data.value.address;
      background64.value = result.data.value.qrCodeData;
    }
  } else {
    step.value = 1;
    message.warn(t('myAssets.notYetProduceAddress'));
  }
};
onMounted(async () => {
  getOptions();
  const route = useRoute();
  formState.value.currencyType = route.query.coinType;
  const haveMainnet = JSON.parse(localStorage.getItem(`${formState.value.currencyType}-mainnet`));
  if (haveMainnet) {
    if (formState.value.currencyType === haveMainnet.currencyType && haveMainnet.currencyType === 'BTC') {
      formData.value.currencyType = 1;
    } else if (formState.value.currencyType === haveMainnet.currencyType && haveMainnet.currencyType === 'ETH') {
      formData.value.currencyType = 2;
    } else if (formState.value.currencyType === haveMainnet.currencyType && haveMainnet.currencyType === 'USDC') {
      formData.value.currencyType = 3;
    } else if (formState.value.currencyType === haveMainnet.currencyType && haveMainnet.currencyType === 'USDT') {
      haveMainnet.currencyType = 4;
    }
    formData.value.mainnet = haveMainnet.mainnet;
    formState.value.mainnet = haveMainnet.mainnetName;
  }
});
</script>
<style>
.ant-select-selector {
  height: 40px !important;
  padding: 5px 10px !important;
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
