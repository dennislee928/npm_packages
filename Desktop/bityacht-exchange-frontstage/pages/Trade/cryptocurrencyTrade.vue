<template>
  <section class="relative 3xl:px-[300px] xl:px-[200px] lg:px-[100px] px-[30px] bg-white md:pb-[140px] pb-[70px]">
    <h1 v-if="step === 1" class="flex items-center justify-center xl:text-h2 md:text-h4 text-[30px] font-bold text-gray_800 text-center pt-[88px]"><left-circle-outlined class="mr-4 text-h5" @click="back" />{{ $t('trade.tradeCryptocurrency2') }}</h1>
    <h1 v-if="step === 2" class="xl:text-h2 md:text-h4 text-[30px] font-bold text-gray_800 text-center pt-[88px]">{{ $t('trade.tradeConfirm') }}</h1>
    <h1 v-if="step === 3 && resultData.status === 1" class="xl:text-h2 md:text-h4 text-[30px] font-bold text-gray_800 text-center pt-[88px]">{{ $t('trade.tradeOrderSend') }}</h1>
    <h1 v-else-if="step === 3 && resultData.status === 2" class="xl:text-h2 md:text-h4 text-[30px] font-bold text-red text-center pt-[88px]">{{ $t('trade.tradeCancel') }}</h1>
    <!-- <img v-if="step === 3" src="/assets/img/frame.png" class="w-[125px] mx-auto my-10" /> -->
    <div class="flex mx-auto md:w-[500px] w-[300px] justify-between mt-5">
      <a-spin :spinning="spinning">
        <p v-if="step === 1" class="flex justify-center border border-gray_800 md:p-5 p-2 text-gray_800 md:w-[200px] w-[145px] text-center md:text-normal text-subText font-semibold">{{ $t('trade.buyPrice') }} : {{ buyPrice }}</p>
      </a-spin>
      <a-spin :spinning="spinning">
        <p v-if="step === 1" class="flex justify-center border border-gray_800 md:p-5 p-2 text-gray_800 md:w-[200px] w-[145px] text-center md:text-normal text-subText font-semibold">{{ $t('trade.sellPrice') }} : {{ sellPrice }}</p>
      </a-spin>
      <p v-if="step === 2" class="flex justify-center mx-auto border border-lightRed p-2 rounded-lg text-lightRed md:w-[220px] w-[145px] text-center md:text-normal text-subText font-semibold"><img src="/icon/funnel.svg" class="mr-2" />{{ $t('trade.please') }} {{ lockTimeNumber }} {{ $t('trade.sec') }}{{ $t('trade.finishTrade') }}</p>
    </div>
    <div class="relative z-30 md:w-[580px] w-full mx-auto flex flex-col justify-center border border-gray_800 rounded-xl md:mt-[40px] mt-8 py-6 md:px-[60px] px-4">
      <div v-if="step === 1" class="absolute -top-5 right-2 text-small text-gray_500">{{ $t('trade.renewPriceText') }}</div>
      <template v-if="step < 3">
        <img src="/assets/img/getSuccess.png" class="md:hidden block w-[450px] mx-auto mb-5" />
        <div class="flex items-center justify-center">
          <div class="flex">
            <img :src="getAssetsFile(`${tradeCouple.side === '2' ? tradeCouple.baseSymbol : tradeCouple.quoteSymbol}.svg`)" class="w-[35px]" />
            <div class="flex flex-col ml-2">
              <span class="font-bold text-subTitle">{{ tradeCouple.side === '2' ? tradeCouple.baseSymbol : tradeCouple.quoteSymbol }}</span>
              <span class="text-subText text-gray_500">{{ tradeCouple.side === '2' ? tradeCouple.baseName : tradeCouple.quoteName }}</span>
            </div>
          </div>
          <caret-right-outlined class="text-waterBlue text-[18px] mx-8 my-5" />
          <!-- <caret-right-outlined v-else class="text-waterBlue text-[18px] mx-8 my-5"  /> -->
          <div class="flex">
            <img :src="getAssetsFile(`${tradeCouple.side === '2' ? tradeCouple.quoteSymbol : tradeCouple.baseSymbol}.svg`)" class="w-[35px]" />
            <div class="flex flex-col ml-2">
              <span class="font-bold text-subTitle">{{ tradeCouple.side === '2' ? tradeCouple.quoteSymbol : tradeCouple.baseSymbol }}</span>
              <span class="text-subText text-gray_500">{{ tradeCouple.side === '2' ? tradeCouple.quoteName : tradeCouple.baseName }}</span>
            </div>
          </div>
        </div>
        <div class="flex justify-center items-center mt-3 mb-1 md:text-subText text-small">{{ $t('trade.warningText') }}</div>
        <a-form :model="formState" @finish="submit" :rules="rules" ref="formRef">
          <div class="flex flex-col mt-6">
            <a-form-item name="payAmount">
              <div class="flex items-center">
                <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('trade.pay') }}</p>
                <p v-if="step === 1" class="xs:text-subText text-small text-gray_400 ml-8 mb-2">{{ $t('trade.canUseRemaingingAmount') }} : {{ canUseRemaingingAmount }} {{ formState.side === 1 ? tradeCouple.quoteSymbol : tradeCouple.baseSymbol }}</p>
              </div>
              <div class="relative">
                <FormInput v-model="formState.payAmount" type="number" :placeholder="$t('signUp.pleaseEnter') + $t('trade.pay')" :disabled="step === 2"><a-button v-if="step === 1" class="absolute top-[50%] -translate-y-1/2 md:right-9 right-2 rounded-2xl py-1 px-3 mt-1" @click="maxCanUse">100%</a-button></FormInput>
              </div>
            </a-form-item>
          </div>
          <div class="flex flex-col mt-3">
            <a-form-item name="earnAmount">
              <div class="flex items-center">
                <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('trade.get') }}</p>
                <a-spin :spinning="spinning">
                  <p v-if="step <= 2" class="xs:text-subText text-small text-gray_400 ml-8 mb-2 font-bold">
                    1 {{ tradeCouple.baseSymbol }} ≈ {{ exchange }}
                    {{ tradeCouple.quoteSymbol }}
                  </p>
                </a-spin>
              </div>
              <FormInput v-model="fakeEarnAmount" disabled />
            </a-form-item>
          </div>
          <div v-if="step === 1" class="flex flex-col items-end">
            <span class="text-small text-gray_500">{{ $t('trade.redemptionFee') }}</span>
            <span class="text-small text-gray_500"> = {{ formState.handlingCharge }} {{ tradeCouple.quoteSymbol }}</span>
            <!-- <exclamation-circle-outlined class="text-normal text-gray_500" /> -->
          </div>
          <div class="flex justify-center mt-6">
            <a-button v-if="step > 1" class="w-[82px] bg-gray_500 text-white rounded-2xl hover:!text-white mr-5" @click="backClear">{{ $t('trade.cancel') }}</a-button>
            <a-button v-if="step === 1" :disabled="!rulesValid" class="w-[82px] bg-waterBlue text-white rounded-2xl hover:!text-white" @click="lockPrice">{{ $t('trade.continue') }}</a-button>
            <a-button v-if="step === 2" class="w-[82px] bg-waterBlue text-white rounded-2xl hover:!text-white" :loading="loading" @click="submit">{{ $t('trade.submit') }}</a-button>
          </div>
        </a-form>
      </template>
      <template v-if="step === 3">
        <div class="flex flex-col items-center">
          <div class="flex justify-between xxs:w-[80%] w-full">
            <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('trade.order') }}</p>
            <p class="xs:text-subText text-small ml-10">{{ resultData.transactionsID }}</p>
          </div>
          <div class="flex justify-between xxs:w-[80%] w-full mt-3">
            <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('trade.tradePair') }}</p>
            <p class="xs:text-subText text-small ml-10">{{ resultData.baseSymbol }} / {{ resultData.quoteSymbol }}</p>
          </div>
          <div class="flex justify-between xxs:w-[80%] w-full mt-3">
            <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('trade.tradeSide') }}</p>
            <p class="xs:text-subText text-small ml-10">{{ resultData.side === 1 ? $t('trade.buy') : $t('trade.sale') }}</p>
          </div>
          <div class="flex justify-between xxs:w-[80%] w-full mt-3">
            <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('trade.pay') }}</p>
            <p class="xs:text-subText text-small ml-10">{{ resultData.amount }} {{ resultData.quoteSymbol }}</p>
          </div>
          <div class="flex justify-between xxs:w-[80%] w-full mt-3">
            <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('trade.get') }}</p>
            <p class="xs:text-subText text-small ml-10">{{ resultData.quantity }} {{ resultData.baseSymbol }}</p>
          </div>
          <div class="flex justify-between xxs:w-[80%] w-full mt-3">
            <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('trade.handlingCharge') }}</p>
            <p class="xs:text-subText text-small ml-10">{{ resultData.handlingCharge }} TWD</p>
          </div>
          <div class="flex justify-between xxs:w-[80%] w-full mt-3">
            <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('trade.tradeTime') }}</p>
            <p class="xs:text-subText text-small ml-10">{{ resultData.createdAt }}</p>
          </div>
        </div>
        <div class="flex justify-center mt-6">
          <LayoutsNavItem :href="`/Trade`" class="break-keep">
            <a-button class="w-[120px] bg-waterBlue text-white rounded-2xl hover:!text-white">{{ $t('trade.backTrade') }}</a-button>
          </LayoutsNavItem>
        </div>
      </template>
    </div>
  </section>
</template>

<script setup>
import useAssetsStore from '@/stores/assets';
import { getAssetsFile, formatValueByDigits } from '@/config/config';
import { Decimal } from 'decimal.js';

const localePath = useLocalePath();
const { t } = useI18n();
useHead({
  title: t('title.trade_tradeCryptocurrency'),
  meta: [{ name: 'description', content: '' }],
});
const assetsStore = useAssetsStore();
const loading = computed(() => assetsStore.loadingButton);
const step = ref(1);
const canUseRemaingingAmount = ref(0);
const route = useRoute();
const tradeCouple = ref(route.query);
const formState = ref({
  payAmount: '',
  earnAmount: 0,
  handlingCharge: 0,
  price: 0,
  side: Number(tradeCouple.value.side),
  baseSymbol: tradeCouple.value.baseSymbol,
  quoteSymbol: tradeCouple.value.quoteSymbol,
});
// const tradeChange = async (side) => {
//   tradeCouple.value.side = side;
//   const runtimeConfig = useRuntimeConfig();
//   let baseUrl = runtimeConfig.app.baseURL ? runtimeConfig.app.baseURL : ' ';
//   baseUrl = baseUrl.slice(0, -1);
//   const href = `${baseUrl}${route.path}?baseName=${tradeCouple.value.baseName}&baseSymbol=${tradeCouple.value.baseSymbol}&quoteName=${tradeCouple.value.quoteName}&quoteSymbol=${tradeCouple.value.quoteSymbol}&side=${tradeCouple.value.side}&handlingChargeRate=${tradeCouple.value.handlingChargeRate}`;
//   location.href = href;
// };
const saletransactionAmount = computed(() => Decimal(formState.value.payAmount).mul(spotInfo.value.sellPrice).toDP(spotInfo.value.basePrecision, Decimal.ROUND_FLOOR));

watch(
  () => formState.value.payAmount,
  (newVal) => {
    getHandlingCharge();
  }
);
const getHandlingCharge = () => {
  let transactionAmount = 0;
    let precision = 0;
    if (formState.value.payAmount === '') {
      transactionAmount = Decimal(0);
    } else if (formState.value.side === 1) {
      // 成交額 = pay / buyPrice || pay * sellPrice
      transactionAmount = Decimal(formState.value.payAmount).div(spotInfo.value.buyPrice).toDP(spotInfo.value.basePrecision, Decimal.ROUND_FLOOR);
      precision = spotInfo.value.basePrecision;
    } else {
      transactionAmount = saletransactionAmount.value;
      precision = spotInfo.value.quotePrecision;
    }
    // 手續費 = 成交額 * (手續費率/100)
    formState.value.handlingCharge = transactionAmount.mul(spotInfo.value.handlingChargeRate).toFixed(precision, Decimal.ROUND_CEIL);
    // 獲得數量 = 成交額 - 手續費
    formState.value.earnAmount = transactionAmount.sub(formState.value.handlingCharge).toDP(precision, Decimal.ROUND_FLOOR);
}
const fakeEarnAmount = computed(() => formState.value.earnAmount.toString());
const rulesValid = ref(false);
const rules = {
  payAmount: {
    require: true,
    validator: async (rule, value) => {
      rulesValid.value = false;
      if (value === '') {
        rulesValid.value = false;
        return Promise.reject(t('rules.dontEmpty'));
      }
      let notinal = formState.value.side === 1 ? value : saletransactionAmount.value;
      let precision = formState.value.side === 1 ? spotInfo.value.quotePrecision : spotInfo.value.basePrecision;
      if (isCopyStatus.value) {
        isCopyStatus.value = false;
        if (Number(notinal) < minNotional.value) {
          return Promise.reject(t('rules.payAmountMinLimit') + `${minNotional.value}` + '' + spotInfo.value.quoteSymbol);
        } else if (Number(notinal) > maxNotional.value) {
          return Promise.reject(t('rules.payAmountMaxLimit') + ' ' + maxNotional.value + ' ' + spotInfo.value.quoteSymbol);
        } else {
          return (rulesValid.value = true);
        }
      }
      if (value.split('.').length == 2 && value.split('.')[1].length > precision) {
        return Promise.reject(t('rules.pleaseNotOver8Digits'));
      } else if (Number(value) > canUseRemaingingAmount.value) {
        rulesValid.value = false;
        return Promise.reject(t('rules.payAmountError'));
      } else if (Number(notinal) < minNotional.value) {
        return Promise.reject(t('rules.payAmountMinLimit') + ' ' + minNotional.value + ' ' + spotInfo.value.quoteSymbol);
      } else if (Number(notinal) > maxNotional.value) {
        return Promise.reject(t('rules.payAmountMaxLimit') + ' ' + maxNotional.value + ' ' + spotInfo.value.quoteSymbol);
      }
      rulesValid.value = true;
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'change',
  },
};
const backClear = () => {
  isCopyStatus.value = false;
  rulesValid.value = false;
  step.value--;
  lockTimeNumber.value = 30;
  formState.value.payAmount = '';
  clearInterval(lockTimer);
  timer = setInterval(() => {
    getTrend();
  }, 5000);
};
const lockTimeNumber = ref(30);
let lockTimer = ref(null);
const lockPrice = () => {
  step.value++;
  clearTimeout(timer);
  lockTimer = setInterval(() => {
    lockTimeNumber.value--;
    if (lockTimeNumber.value === 0) {
      backClear();
      clearTimeout(lockTimer);
    }
  }, 1000);
};
const back = async () => {
  await navigateTo(localePath('/Trade'));
};
const isCopyStatus = ref(false);
const formRef = ref(null);
const maxCanUse = () => {
  formState.value.payAmount = Number(canUseRemaingingAmount.value);
  isCopyStatus.value = true;
  formRef.value.validate(['payAmount']);
};
const userInfo = ref(undefined);
const getAssets = async () => {
  const result = await assetsStore.getAssets(t);
  userInfo.value = result.data.value.assetInfos;
  userInfo.value.forEach((item) => {
    if (formState.value.side === 1) {
      if (item.currenciesSymbol === tradeCouple.value.quoteSymbol) {
        canUseRemaingingAmount.value = item.freeAmount;
      }
    } else {
      if (item.currenciesSymbol === tradeCouple.value.baseSymbol) {
        canUseRemaingingAmount.value = item.freeAmount;
      }
    }
  });
};
const spotInfo = ref(undefined);
const exchange = ref(0);
const minNotional = ref(0);
const maxNotional = ref(0);
const sellPrice = ref(0);
const buyPrice = ref(0);
const spinning = ref(false);
let timeOut = ref(null);

const getTrend = async () => {
  spinning.value = true;
  const result = await assetsStore.getTrend(t);
  timeOut = setTimeout(() => {
    spinning.value = false;
  }, 500);
  result.data.value.forEach((item) => {
    if (item.baseSymbol === tradeCouple.value.baseSymbol && item.quoteSymbol === tradeCouple.value.quoteSymbol) {
      sellPrice.value = formatValueByDigits(item.sellPrice, 2);
      buyPrice.value = formatValueByDigits(item.buyPrice, 2);
      spotInfo.value = item;
      spotInfo.value.buyPrice = Decimal(spotInfo.value.buyPrice);
      spotInfo.value.sellPrice = Decimal(spotInfo.value.sellPrice);
      spotInfo.value.handlingChargeRate = Decimal(spotInfo.value.handlingChargeRate).div(100);
    }
  });
  minNotional.value = Number(spotInfo.value.minNotional);
  maxNotional.value = Number(spotInfo.value.maxNotional);

  let transactionAmount = 0;
  let precision = 0;
  if (formState.value.payAmount === '') {
    transactionAmount = Decimal(0);
  } else if (formState.value.side === 1) {
    // 成交額 = pay / buyPrice || pay * sellPrice
    transactionAmount = Decimal(formState.value.payAmount).div(spotInfo.value.buyPrice).toDP(spotInfo.value.basePrecision, Decimal.ROUND_FLOOR);
    precision = spotInfo.value.basePrecision;
  } else {
    transactionAmount = saletransactionAmount.value;
    precision = spotInfo.value.quotePrecision;
  }
  formState.value.earnAmount = transactionAmount.sub(formState.value.handlingCharge).toDP(precision, Decimal.ROUND_FLOOR);
  if (formState.value.side === 1) {
    // exchange.value = Decimal(1).div(spotInfo.value.buyPrice).toDP(spotInfo.value.basePrecision, Decimal.ROUND_CEIL);
    exchange.value = formatValueByDigits(spotInfo.value.buyPrice, 2);
    formState.value.price = spotInfo.value.buyPrice;
  } else {
    exchange.value = formatValueByDigits(spotInfo.value.sellPrice, 2);
    // exchange.value = spotInfo.value.sellPrice.toDP(spotInfo.value.quotePrecision, Decimal.ROUND_CEIL);
    formState.value.price = spotInfo.value.sellPrice;
  }
  getHandlingCharge();
};

const resultData = ref(undefined);
const submit = async () => {
  clearInterval(lockTimer);
  assetsStore.loadingButton = true;
  formState.value.earnAmount = formState.value.earnAmount.toString();
  formState.value.handlingCharge = formState.value.handlingCharge.toString();
  const result = await assetsStore.transactions(formState.value, t);
  assetsStore.loadingButton = false;
  if (result.status.value === 'success') {
    resultData.value = result.data.value;
    step.value++;
    if (result.data.value.status === 1) {
      message.success(t('trade.tradeSuccess'));
    } else {
      message.error(t('trade.tradeCancel'));
    }
  }
};

let timer = ref(null);
onMounted(async () => {
  await getAssets();
  await getTrend();
  timer = setInterval(() => {
    getTrend();
  }, 5000);
});
onUnmounted(() => {
  clearTimeout(timer);
  clearTimeout(timeOut);
});
</script>
<style>
input[type='number']::-webkit-outer-spin-button,
input[type='number']::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}
</style>
