<template>
  <section class="relative 3xl:px-[300px] xl:px-[200px] lg:px-[100px] px-[15px] bg-white pb-10">
    <h1 class="relative z-30 flex items-center justify-center md:text-h2 text-h4 text-gray_800 pt-[88px] md:pb-[60px] pb-[40px]">{{ $t('trade.trade') }}</h1>
    <div class="relative z-30 flex justify-center md:mb-[60px] mb-6">
      <a-button :class="`md:w-[145px] w-[40%] md:text-subTitle text-subText md:py-5 py-3 rounded-3xl flex items-center justify-center md:mr-5 mr-4 ${selectCoinType === 'legalTender' ? 'bg-waterBlue text-white hover:!text-white' : ''}`" @click="coinTypeSelect('legalTender')">{{ $t('trade.tradeCryptocurrency') }}</a-button>
      <a-button :class="`md:w-[145px] w-[40%] md:text-subTitle text-subText md:py-5 py-3  rounded-3xl flex items-center justify-center md:mr-0 sm:mr-4 mr-0 ${selectCoinType === 'cryptocurrency' ? 'bg-waterBlue text-white hover:!text-white' : ''}`" @click="coinTypeSelect('cryptocurrency')">{{ $t('trade.tradeLegelCount') }}</a-button>
    </div>
    <div class="relative z-30 border border-gray_800 rounded-xl 3xl:px-10 xl:px-8 md:px-5 xs:px-10 px-6 xl:py-10 py-6">
      <div class="flex justify-between">
        <p class="flex items-center font-bold text-subTitle flex-wrap">
          <img src="/icon/hexagon.svg" class="mr-2" />{{ $t('trade.myAssets') }} <span class="text-normal font-normal ml-1">({{ $t('trade.totalAssets') }}: {{ formatValueByDigits(totalCount,8) }})</span>
        </p>
        <NuxtLink :to="localePath('/Trade/tradeRecord')" class="md:block hidden">
          <div class="w-[140px] mx-auto border border-gray_800 rounded-2xl py-1 cursor-pointer">
            <div class="flex justify-center items-center">
              <img src="/icon/walletGray.svg" />
              <span class="text-gray_800 text-normal ml-3">{{ $t('trade.tradeRecord') }}</span>
            </div>
          </div>
        </NuxtLink>
      </div>
      <div class="flex md:flex-row flex-col mt-5 justify-around">
        <div v-for="(item, index) of data" :key="`${item.currenciesSymbol}+${index}`" class="md:w-auto w-full flex md:flex-col flex-row items-center justify-between md:mx-2 mx-0 md:mb-0 mb-5 md:border-b-0 border-b border-gray_300 md:pb-0 pb-5 last:border-0">
          <div class="flex items-center">
            <img :src="getAssetsFile(`${item.currenciesSymbol}.svg`)" class="xs:w-10 w-8 xs:h-10 h-8 mr-3" />
            <div class="flex flex-col items-start">
              <p class="font-bold xs:text-subTitle text-normal text-black">{{ item.currenciesSymbol }}</p>
              <p class="xs:text-subText text-small text-grayText">{{ item.currenciesName }}</p>
            </div>
          </div>
          <p class="md:block hidden border border-black  w-10 my-2"></p>
          <div class="xs:text-subTitle text-normal">
            {{ formatValueByDigits(item.freeAmount,8) }}
          </div>
        </div>
      </div>
      <div class="bg-gray_800 rounded-xl py-3 cursor-pointer mb-2 md:hidden block">
        <NuxtLink :to="localePath('/Trade/tradeRecord')">
          <div class="flex justify-center items-center">
            <img src="/icon/wallet.svg" class="xs:w-[30px] w-[25px]" />
            <span class="text-white xs:text-[22px] text-subTitle ml-3">{{ $t('trade.tradeRecord') }}</span>
          </div>
        </NuxtLink>
      </div>
    </div>
    <div class="relative z-30 mt-10 md:block hidden">
      <div class="flex items-center justify-center">
        <p class="text-gray_800 md:text-h4 text-[28px] w-[50%] text-center cursor-pointer" :class="`${isBuy ? `border-t border-l border-r border-gray_800 rounded-tr-lg rounded-tl-lg` : `border-b border-gray_800`}`" @click="marketAction('buy')">
          <span v-if="!isBuy" class="block w-[95%] h-full bg-waterBlue text-white px-5 mx-auto rounded-tr-lg rounded-tl-lg border-t border-l border-r border-gray_800">{{ $t('index.buyCurrency') }}</span
          >{{ isBuy ? $t('index.buyCurrency') : '' }}
        </p>
        <p class="text-gray_800 md:text-h4 text-[28px] w-[50%] text-center cursor-pointer" :class="`${isBuy ? `border-b border-gray_800` : `border-t border-l border-r border-gray_800 rounded-tl-lg rounded-tr-lg`}`" @click="marketAction('sale')">
          <span v-if="isBuy" class="block w-[95%] h-full bg-waterBlue text-white px-5 mx-auto rounded-tr-lg rounded-tl-lg border-t border-l border-r border-gray_800">{{ $t('index.saleCurrency') }}</span
          >{{ !isBuy ? $t('index.saleCurrency') : '' }}
        </p>
      </div>
      <div class="border-l border-r border-b border-gray_800 rounded-br-lg rounded-bl-lg">
        <Table :data="tableData.data" :columns="columns" :pagination="false" class="tradeTable">
          <template v-slot:symbol="{ record }">
            <div class="flex items-center justify-between">
              <div class="flex w-[35%] 3xl:justify-start justify-end">
                <img :src="getAssetsFile(`${!isBuy ? record.baseSymbol : record.quoteSymbol}.svg`)" class="w-[30px]" />
                <div class="3xl:flex hidden flex-col ml-2">
                  <span class="font-bold text-normal">{{ !isBuy ? record.baseSymbol : record.quoteSymbol }}</span>
                  <span class="text-small text-gray_500">{{ !isBuy ? record.baseName : record.quoteName }}</span>
                </div>
              </div>
              <caret-right-outlined class="text-gray_500" />
              <div class="flex w-[35%] justify-start">
                <img :src="getAssetsFile(`${!isBuy ? record.quoteSymbol : record.baseSymbol}.svg`)" class="w-[30px]" />
                <div class="3xl:flex hidden flex-col ml-2">
                  <span class="font-bold text-normal">{{ !isBuy ? record.quoteSymbol : record.baseSymbol }}</span>
                  <span class="text-small text-gray_500">{{ !isBuy ? record.quoteName : record.baseName }}</span>
                </div>
              </div>
            </div>
          </template>
          <template v-slot:buyPrice="{ record }">
            <div class="flex items-end justify-center flex-wrap">
              <span class="text-gray_500">1</span>
              <span class="text-gray_500 ml-1 text-small">{{ record.baseSymbol }} </span>
              <span class="text-gray_500 mx-1">≈ </span>
              <span class="text-gray_500 text-subText">{{ formatValue(record.buyPrice) }}</span>
              <span class="text-gray_500 ml-1 text-small">{{ record.quoteSymbol }}</span>
            </div>
          </template>
          <template v-slot:sellPrice="{ record }">
            <div class="flex items-end justify-center flex-wrap">
              <span class="text-gray_500">1</span>
              <span class="text-gray_500 ml-1 text-small">{{ record.baseSymbol }}</span>
              <span class="text-gray_500 mx-1">≈ </span>
              <span class="text-gray_500 text-subText">{{ formatValue(record.sellPrice) }}</span>
              <span class="text-gray_500 ml-1 text-small">{{ record.quoteSymbol }}</span>
            </div>
          </template>
          <template v-slot:updown="text">
            <span v-if="Number(text.text) > 0" class="xxs:text-subText text-small text-orange">{{ text.text }}%</span>
            <span v-else class="xxs:text-subText text-small text-green">{{ text.text }}%</span>
          </template>
          <template v-if="dataReady" v-slot:trend="{ record }">
            <Echarts :id="record.baseSymbol + record.quoteSymbol" :series="record.trend" :max="record.max" :min="record.min" location="trade" :show="false" symbol="none" :inverse="false" />
          </template>
          <template v-slot:operate="{ record }">
            <a-button class="bg-waterBlue text-white rounded-2xl hover:!text-white" @click="tradeRouter(record, isBuy ? 1 : 2)"> {{ isBuy ? $t('trade.buy') : $t('trade.sale') }}</a-button>
          </template>
        </Table>
        <div v-if="closeData.data.length !== 0" class="bg-gray_800 text-white mx-10 py-2.5 rounded-lg font-bold text-subText 2xl:pl-[100px] pl-[80px]">
          {{ $t('trade.temporaryClose') }}
        </div>
        <Table v-if="closeData.data.length !== 0" :data="closeData.data" :columns="closeColumns" :pagination="false" :show-header="false" class="tradeTable opacity-50">
          <template v-slot:symbol="{ record }">
            <div class="flex items-center justify-between">
              <div class="flex w-[35%] 3xl:justify-start justify-end">
                <img :src="getAssetsFile(`${record.baseSymbol}.svg`)" class="w-[30px]" />
                <div class="3xl:flex hidden flex-col ml-2">
                  <span class="font-bold text-normal">{{ record.baseSymbol }}</span>
                  <span class="text-small text-gray_500">{{ record.baseName }}</span>
                </div>
              </div>
              <caret-left-outlined class="text-gray_500" />
              <!-- <caret-right-outlined v-else class="text-gray_500" /> -->
              <div class="flex w-[35%] 3xl:justify-start">
                <img :src="getAssetsFile(`${record.quoteSymbol}.svg`)" class="w-[30px]" />
                <div class="3xl:flex hidden flex-col ml-2">
                  <span class="font-bold text-normal">{{ record.quoteSymbol }}</span>
                  <span class="text-small text-gray_500">{{ record.quoteName }}</span>
                </div>
              </div>
            </div>
          </template>
          <template v-slot:buyPrice="{ record }">
            <div class="flex items-center justify-center">
              <template v-if="Number(record.buyPrice) !== 0">
                <span>1</span>
                <span class="text-gray_500 ml-1 text-small">{{ record.quoteSymbol }}</span>
                <span> = {{ formatValue(1 / record.buyPrice) }}</span>
                <span class="text-gray_500 ml-1 text-small">{{ record.baseSymbol }}</span>
              </template>
              <template v-else>-</template>
            </div>
          </template>
          <template v-slot:sellPrice="{ record }">
            <div class="flex items-center">
              <span>1</span>
              <span class="text-gray_500 ml-1 text-small">{{ record.baseSymbol }}</span>
              <span> = {{ record.buyPrice }}</span>
              <span class="text-gray_500 ml-1 text-small">{{ record.quoteSymbol }}</span>
            </div>
          </template>
          <template v-slot:updown="text">
            <span v-if="Number(text.text) > 0" class="xxs:text-subText text-small text-orange">{{ text.text }}%</span>
            <span v-else class="xxs:text-subText text-small text-green">{{ text.text }}%</span>
          </template>
          <template v-if="dataReady" v-slot:trend="{ record }">
            <Echarts :id="record.baseSymbol + record.quoteSymbol" :series="record.trend" :max="record.max" :min="record.min" location="trade" :show="false" symbol="none" :inverse="false" />
          </template>
          <template v-slot:operate="{ record }">
            <a-button class="bg-waterBlue text-white rounded-2xl hover:!text-white" :disabled="true"> {{ isBuy ? $t('trade.buy') : $t('trade.sale') }}</a-button>
          </template>
        </Table>
      </div>
    </div>
    <div class="relative z-30 mt-10 md:hidden block">
      <div class="flex -mx-[30px]">
        <p :class="`w-[50%] text-center py-3 text-subTitle font-bold cursor-pointer ${isBuy ? 'bg-white text-gray_800 border-t border-r border-gray_800' : 'bg-waterBlue text-white border-b border-gray_800'}`" @click="marketAction('buy')">{{ $t('trade.buy') }}</p>
        <p :class="`w-[50%] text-center py-3 text-subTitle font-bold cursor-pointer ${isBuy ? 'bg-waterBlue text-white border-b border-gray_800' : 'bg-white text-gray_800 border-t border-r border-gray_800'}`" @click="marketAction('sale')">{{ $t('trade.sale') }}</p>
      </div>
      <div class="bg-gray_800 text-white text-center mt-10 rounded-xl py-2">{{ $t('trade.tradePair') }}</div>
      <div v-for="(item, index) of tableData.data">
        <div class="border border-gray_800 rounded-xl p-3 mt-10">
          <div class="flex items-center xs:justify-center justify-between">
            <div class="flex xs:mx-3 mx-0">
              <img :src="getAssetsFile(`${!isBuy ? item.baseSymbol : item.quoteSymbol}.svg`)" class="w-[30px]" />
              <div class="flex flex-col ml-2">
                <span class="font-bold text-normal">{{ !isBuy ? item.baseSymbol : item.quoteSymbol }}</span>
                <span class="text-small text-gray_500">{{ !isBuy ? item.baseName : item.quoteName }}</span>
              </div>
            </div>
            <div class="flex xs:mx-3 mx-0">
              <caret-right-outlined class="text-gray_500" />
              <caret-right-outlined class="text-gray_500" />
              <caret-right-outlined class="text-gray_500" />
            </div>
            <div class="flex xs:mx-3 mx-0">
              <img :src="getAssetsFile(`${!isBuy ? item.quoteSymbol : item.baseSymbol}.svg`)" class="w-[30px]" />
              <div class="flex flex-col ml-2">
                <span class="font-bold text-normal">{{ !isBuy ? item.quoteSymbol : item.baseSymbol }}</span>
                <span class="text-small text-gray_500">{{ !isBuy ? item.quoteName : item.baseName }}</span>
              </div>
            </div>
          </div>
          <div class="flex items-center justify-center mt-2 flex-wrap">
            <span class="font-bold text-subText mr-1">1</span>
            <span class="text-gray_500 text-small sm:mx-2 mx-1">{{ item.baseSymbol }}</span>
            <span class="font-bold text-subText mr-1">≈ {{ isBuy ? formatValueByDigits(item.buyPrice, 8) : formatValueByDigits(item.sellPrice, 8) }}</span>
            <span class="text-gray_500 text-small sm:mx-2 mx-1">{{ isBuy ? item.quoteSymbol : item.quoteSymbol }}</span>
            <img src="/assets/img/twoWay.svg" class="sm:mx-2 mx-1 xxs:w-auto w-[15px]" />
            <span v-if="Number(item.upsAndDowns) > 0" class="xxs:text-subText text-small text-orange sm:mx-2 mx-1">{{ item.upsAndDowns }}%</span>
            <span v-else class="xxs:text-subText text-small text-green sm:mx-2 mx-1">{{ item.upsAndDowns }}%</span>
          </div>
        </div>
        <div class="flex justify-center mt-5 xs:mb-[60px] mb-10">
          <a-button class="bg-waterBlue text-white rounded-2xl hover:!text-white" @click="mobileTradeRouter(index, isBuy ? 1 : 2)"> {{ isBuy ? $t('trade.buy') : $t('trade.sale') }}</a-button>
        </div>
      </div>
      <template v-if="closeData.data.length !== 0">
        <div class="bg-gray_800 text-white text-center mt-10 rounded-xl py-2">{{ $t('trade.temporaryClose') }}</div>
        <div v-for="item of closeData.data" class="opacity-50">
          <div class="border border-gray_800 rounded-xl p-3 mt-10">
            <div class="flex items-center xs:justify-center justify-between">
              <div class="flex xs:mx-3 mx-0">
                <img :src="getAssetsFile(`${item.baseSymbol}.svg`)" class="w-[30px]" />
                <div class="flex flex-col ml-2">
                  <span class="font-bold text-normal">{{ item.baseSymbol }}</span>
                  <span class="text-small text-gray_500">{{ item.baseName }}</span>
                </div>
              </div>
              <div class="flex xs:mx-3 mx-0">
                <template v-if="isBuy">
                  <caret-left-outlined class="text-gray_500" />
                  <caret-left-outlined class="text-gray_500" />
                  <caret-left-outlined class="text-gray_500" />
                </template>
                <template v-else>
                  <caret-right-outlined class="text-gray_500" />
                  <caret-right-outlined class="text-gray_500" />
                  <caret-right-outlined class="text-gray_500" />
                </template>
              </div>
              <div class="flex xs:mx-3 mx-0">
                <img :src="getAssetsFile(`${item.quoteSymbol}.svg`)" class="w-[30px]" />
                <div class="flex flex-col ml-2">
                  <span class="font-bold text-normal">{{ item.quoteSymbol }}</span>
                  <span class="text-small text-gray_500">{{ item.quoteName }}</span>
                </div>
              </div>
            </div>
            <div class="flex items-center justify-center mt-2 flex-wrap">
              <span class="font-bold text-subText mr-1">1</span>
              <span class="text-gray_500 text-small sm:mx-2 mx-1">{{ isBuy ? item.quoteSymbol : item.baseSymbol }}</span>
              <span class="font-bold text-subText mr-1">≈ {{ isBuy ? formatValueByDigits(item.buyPrice, 8) : formatValueByDigits(item.sellPrice, 8) }}</span>
              <span class="text-gray_500 text-small sm:mx-2 mx-1">{{ isBuy ? item.baseSymbol : item.quoteSymbol }}</span>
              <img src="/assets/img/twoWay.svg" class="sm:mx-2 mx-1 xxs:w-auto w-[15px]" />
              <span v-if="Number(item.upsAndDowns) > 0" class="xxs:text-subText text-small text-orange sm:mx-2 mx-1">{{ item.upsAndDowns }}%</span>
              <span v-else class="xxs:text-subText text-small text-green sm:mx-2 mx-1">{{ item.upsAndDowns }}%</span>
            </div>
          </div>
          <div class="flex justify-center mt-5 xs:mb-[60px] mb-10">
            <a-button class="bg-waterBlue text-white rounded-2xl hover:!text-white" :disabled="true"> {{ isBuy ? $t('trade.buy') : $t('trade.sale') }}</a-button>
          </div>
        </div>
      </template>
    </div>
  </section>
  <img src="/assets/img/Group16.png" class="absolute right-0 top-0 md:w-[600px] w-[400px] z-10" />
  <img src="/assets/img/Group37.png" class="absolute left-0 top-0 md:w-[800px] w-[200px] z-10" />
</template>
<script setup>
import useAssetsStore from '@/stores/assets';
import { formatValue, formatValueByDigits, getAssetsFile } from '@/config/config';

const localePath = useLocalePath();
const { t } = useI18n();
useHead({
  title: t('title.trade'),
  meta: [{ name: 'description', content: '' }],
});
const assetsStore = useAssetsStore();
const selectCoinType = ref('legalTender');
const coinTypeSelect = (type) => {
  if (type === 'cryptocurrency') {
    message.warning(t('trade.comingSoon'));
    return;
  }
  selectCoinType.value = type;
};
const isBuy = ref(true);
const marketAction = (type) => {
  if (type === 'buy') {
    isBuy.value = true;
    columns.value[1] = {
      title: 'trade.buyValue',
      dataIndex: 'buyPrice',
      key: 'buyPrice',
      width: 250,
    };
    getTrend();
  } else {
    isBuy.value = false;
    columns.value[1] = {
      title: 'trade.sellValue',
      dataIndex: 'sellPrice',
      key: 'sellPrice',
      width: 250,
    };
    getTrend();
  }
};
const data = ref([]);
const columns = ref([
  {
    title: 'trade.tradePair',
    dataIndex: 'symbol',
    key: 'symbol',
    width: 250,
  },
  {
    title: 'trade.buyValue',
    dataIndex: 'buyPrice',
    key: 'buyPrice',
    width: 250,
  },
  {
    title: 'trade.upAndDown24h',
    dataIndex: 'upsAndDowns',
    key: 'upsAndDowns',
    width: 80,
  },
  {
    title: 'trade.priceTrend',
    dataIndex: 'trend',
    key: 'trend',
  },
  {
    title: 'trade.action',
    dataIndex: 'operate',
    key: 'operate',
  },
]);
const closeColumns = [
  {
    title: 'trade.temporaryClose',
    dataIndex: 'symbol',
    key: 'symbol',
    width: 250,
  },
  {
    title: 'trade.space',
    dataIndex: 'buyPrice',
    key: 'buyPrice',
    width: 250,
  },
  {
    title: 'trade.space',
    dataIndex: 'upsAndDowns',
    key: 'upsAndDowns',
    width: 80,
  },
  {
    title: 'trade.space',
    dataIndex: 'trend',
    key: 'trend',
  },
  {
    title: 'trade.space',
    dataIndex: 'operate',
    key: 'operate',
  },
];
const tableData = reactive({ data: [] });
const closeData = reactive({ data: [] });
const totalCount = ref(0);

const getTrend = async () => {
  const result = await assetsStore.getTrend(t);
  if (result.data.value === null) return;
  closeData.data = [];
  for (let index = result.data.value.length - 1; index >= 0; index--) {
    const item = result.data.value[index];
    item.trend = item.trend.map((str) => parseFloat(str));
    item.upsAndDowns = parseFloat(item.upsAndDowns).toFixed(2);
    item.max = Math.max(...item.trend);
    item.min = Math.min(...item.trend);
    if (item.status === 0) {
      closeData.data.push(result.data.value.splice(index, 1)[0]);
    }
  }
  tableData.data = result.data.value;
};
const getAssets = async () => {
  const result = await assetsStore.getAssets(t);
  totalCount.value = Number(result.data.value.cryptocurrencyValuation) + Number(result.data.value.fiatBalance);
  data.value = result.data.value.assetInfos;
  data.value.unshift({
    currenciesName: 'TWD',
    currenciesSymbol: 'TWD',
    freeAmount: result.data.value.fiatBalance,
  });
};
const tradeRouter = async (record, side) => {
  const tradeCouple = {
    baseName: record.baseName,
    baseSymbol: record.baseSymbol,
    quoteName: record.quoteName,
    quoteSymbol: record.quoteSymbol,
    side: side,
    handlingChargeRate: record.handlingChargeRate,
  };
  await navigateTo({
    path: localePath('/Trade/cryptocurrencyTrade'),
    query: tradeCouple,
  });
};
const mobileTradeRouter = async (index, side) => {
  const tradeCouple = {
    baseName: tableData.data[index].baseName,
    baseSymbol: tableData.data[index].baseSymbol,
    quoteName: tableData.data[index].quoteName,
    quoteSymbol: tableData.data[index].quoteSymbol,
    side: side,
    handlingChargeRate: tableData.data[index].handlingChargeRate,
  };
  await navigateTo({
    path: localePath('/Trade/cryptocurrencyTrade'),
    query: tradeCouple,
  });
};
const dataReady = ref(false);
let timer = ref(null);
onMounted(async () => {
  dataReady.value = true;
  await getTrend();
  await getAssets();
  timer = setInterval(() => {
    getTrend();
  }, 5000);
});
onUnmounted(() => {
  clearTimeout(timer);
});
</script>

<style>
.tradeTable .ant-table table {
  background: #ffffff !important;
}
.tradeTable .ant-table-thead > tr > th {
  padding: 10px 12px !important;
  background: #7589a4 !important;
  text-align: center !important;
  color: #fff !important;
  border: none !important;
}
.tradeTable .ant-table-tbody > tr > td {
  padding: 10px 12px !important;
  color: #181b22 !important;
  font-weight: 400 !important;
  text-align: center !important;
  color: #1d1d1d !important;
  border-bottom: 1px solid rgb(241, 241, 241) !important;
}
/* .tradeTable .ant-table-tbody > tr > td:last-child {
  display: flex;
  justify-content: center;
} */
.tradeTable .ant-table-tbody .textRight {
  text-align: right !important;
}
.tradeTable .ant-table-tbody > tr > td.ant-table-cell-row-hover {
  background-color: #f2f2f2 !important;
}
.tradeTable :where(.css-dev-only-do-not-override-eq3tly).ant-table-wrapper .ant-table-tbody > tr.ant-table-row:hover > td,
.tradeTable :where(.css-dev-only-do-not-override-eq3tly).ant-table-wrapper .ant-table-tbody > tr > td.ant-table-cell-row-hover {
  background: #f2f2f2 !important;
}
:where(.css-dev-only-do-not-override-eq3tly).ant-table-wrapper .ant-table-container table > thead > tr:first-child > *:first-child {
  border-end-start-radius: 8px;
}
:where(.css-dev-only-do-not-override-eq3tly).ant-table-wrapper .ant-table-container table > thead > tr:first-child > *:last-child {
  border-end-end-radius: 8px;
}
</style>
