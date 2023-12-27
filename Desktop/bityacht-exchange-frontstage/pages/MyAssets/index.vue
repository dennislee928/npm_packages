<template>
  <section class="relative 3xl:px-[300px] xl:px-[200px] lg:px-[100px] px-[30px] bg-white pb-[100px]">
    <div class="relative z-30">
      <div class="flex justify-center items-center">
        <h1 class="md:text-h2 text-h4 text-gray_800 text-center pt-[88px] md:pb-[60px] pb-[20px]">{{ $t('myAssets.myAssets') }}</h1>
        <eye-invisible-outlined v-if="!eyes" class="text-[26px] text-grayBorder md:mt-7 mt-[70px] ml-3" @click="eyes = !eyes" />
        <eye-outlined v-else class="text-[26px] text-gray_800 md:mt-7 mt-[70px] ml-3" @click="eyes = !eyes" />
      </div>
      <NuxtLink :to="localePath('/myAssets/walletRecord')">
        <div class="w-[150px] mx-auto border border-gray_800 rounded-2xl py-1 cursor-pointer md:mb-2 mb-10 md:hidden block">
          <div class="flex justify-center items-center">
            <img src="/icon/walletGray.svg" />
            <span class="text-gray_800 md_lg:text-[26px] text-[18px] ml-3">{{ $t('myAssets.walletRecord') }}</span>
          </div>
        </div>
      </NuxtLink>
      <div class="flex justify-between md:flex-row flex-col">
        <div class="md:w-[26%] w-full flex flex-col md:order-1 order-2">
          <div class="bg-gray_800 rounded-xl py-3 cursor-pointer mb-2 md:block hidden">
            <NuxtLink :to="localePath('/myAssets/walletRecord')">
              <div class="flex justify-center items-center">
                <img src="/icon/wallet.svg" />
                <span class="text-white md_lg:text-[26px] text-[20px] ml-3">{{ $t('myAssets.walletRecord') }}</span>
              </div>
            </NuxtLink>
          </div>
          <div class="border border-gray_800 2xl:px-6 px-3 py-6 rounded-xl md:mb-2 mb-8">
            <div class="flex flex-col">
              <div class="flex items-center">
                <img src="/icon/hexagon.svg" alt="icon" class="mr-2 w-[25px]" />
                <span class="text-black text-subTitle font-bold">{{ $t('myAssets.legelCount') }}</span>
                <span class="text-grayBorder text-normal"> (TWD)</span>
              </div>
              <p class="text-waterBlue mt-5 ml-8">{{ eyes ? formatValue(legelCount) : '********' }}</p>
              <div class="flex md_lg:mt-10 md:mt-5 mt-2 md_lg:justify-end md:justify-center justify-end md_lg:flex-row md:flex-col flex-row items-center opacity-0">
                <a-button class="md_lg:w-auto w-[80px] flex items-center rounded-3xl text-small bg-waterBlue text-white md_lg:mr-2 md:mr-0 mr-2 md_lg:mb-0 md:mb-2 mb-0 hover:!text-white"> <download-outlined />{{ $t('myAssets.inputMoney') }} </a-button>
                <a-button class="md_lg:w-auto w-[80px] flex items-center rounded-3xl text-small"> <upload-outlined />{{ $t('myAssets.outputMoney') }} </a-button>
              </div>
            </div>
          </div>
          <div class="border border-gray_800 2xl:px-6 px-3 py-6 rounded-xl">
            <div class="flex flex-col">
              <div class="flex items-center">
                <img src="/icon/hexagon.svg" alt="icon" class="mr-2 w-[25px]" />
                <span class="text-black text-subTitle font-bold">{{ $t('myAssets.digitAssets') }}</span>
                <span class="text-grayBorder text-normal"> (TWD)</span>
              </div>
              <p class="text-waterBlue mt-5 ml-8">{{ eyes ? formatValue(digitCount) : '********' }}</p>
            </div>
          </div>
        </div>
        <div class="relative md:w-[73%] w-full md:order-2 order-1 md:mb-0 mb-8">
          <div class="border border-gray_800 pt-6 md_lg:pb-[50px] pb-[120px] xs:px-6 px-3 rounded-xl">
            <div class="flex flex-col">
              <div class="flex items-center">
                <img src="/icon/hexagon.svg" alt="icon" class="mr-2 w-[25px]" />
                <span class="text-black text-subTitle font-bold">{{ $t('myAssets.totalAssets') }}</span>
                <span class="text-grayBorder text-normal"> (TWD)</span>
              </div>
              <p class="text-waterBlue mt-5 ml-8">{{ eyes ? formatValue(totalCount) : '********' }}</p>
            </div>
            <div v-if="dataReady">
              <Echarts id="haha" :series="trend" :max="max" :min="min" location="myAssets" :show="false" symbol="circle" :inverse="true" />
            </div>
          </div>
          <div class="flex justify-center w-full absolute top-[80%] left-[50%] -translate-x-1/2">
            <a-button :class="`md:w-auto xs:w-[100px] w-[60px] rounded-3xl text-small mr-3 ${timeSelect === 'day' ? 'bg-waterBlue text-white hover:!text-white' : ''}`" @click="selectTime('day')">{{ $t('myAssets.day') }}</a-button>
            <a-button :class="`md:w-auto xs:w-[100px] w-[60px] rounded-3xl text-small mr-3 ${timeSelect === 'week' ? 'bg-waterBlue text-white hover:!text-white' : ''}`" @click="selectTime('week')">{{ $t('myAssets.week') }}</a-button>
            <a-button :class="`md:w-auto xs:w-[100px] w-[60px] rounded-3xl text-small mr-3 ${timeSelect === 'month' ? 'bg-waterBlue text-white hover:!text-white' : ''}`" @click="selectTime('month')">{{ $t('myAssets.month') }}</a-button>
            <a-button :class="`md:w-auto xs:w-[100px] w-[60px] rounded-3xl text-small  ${timeSelect === 'year' ? 'bg-waterBlue text-white hover:!text-white' : ''}`" @click="selectTime('year')">{{ $t('myAssets.years') }}</a-button>
          </div>
        </div>
      </div>
      <div class="myAssetsTable md:block hidden border border-gray_800 p-6 rounded-xl 2xl:mt-4 xl:mt-3 mt-2">
        <Table :columns="columns" :data="data" :pagination="false">
          <template v-slot:currenciesSymbol="{ record }">
            <template v-if="record.currenciesSymbol === 'BTC'">
              <div class="flex items-center">
                <img src="/assets/img/BTC.svg" class="w-10 h-10 mr-3" />
                <div class="flex flex-col items-start">
                  <p class="font-bold text-subTitle text-black">{{ record.currenciesSymbol }}</p>
                  <p class="text-subText text-grayText">{{ record.currenciesName }}</p>
                </div>
              </div>
            </template>
            <template v-if="record.currenciesSymbol === 'ETH'">
              <div class="flex items-center">
                <img src="/assets/img/ETH.svg" class="w-10 h-10 mr-3" />
                <div class="flex flex-col items-start">
                  <p class="font-bold text-subTitle text-black">{{ record.currenciesSymbol }}</p>
                  <p class="text-subText text-grayText">{{ record.currenciesName }}</p>
                </div>
              </div>
            </template>
            <template v-if="record.currenciesSymbol === 'USDT'">
              <div class="flex items-center">
                <img src="/assets/img/USDT.svg" class="w-10 h-10 mr-3" />
                <div class="flex flex-col items-start">
                  <p class="font-bold text-subTitle text-black">{{ record.currenciesSymbol }}</p>
                  <p class="text-subText text-grayText">{{ record.currenciesName }}</p>
                </div>
              </div>
            </template>
            <template v-if="record.currenciesSymbol === 'USDC'">
              <div class="flex items-center">
                <img src="/assets/img/USDC.svg" class="w-10 h-10 mr-3" />
                <div class="flex flex-col items-start">
                  <p class="font-bold text-subTitle text-black">{{ record.currenciesSymbol }}</p>
                  <p class="text-subText text-grayText">{{ record.currenciesName }}</p>
                </div>
              </div>
            </template>
          </template>
          <template v-slot:valuation="{ text }">
            <p v-if="eyes" class="text-right">≈ ${{ formatValue(text) }}</p>
            <p v-else class="text-right">********</p>
          </template>
          <template v-slot:operate="{ record }">
            <a-button :class="` my-3 rounded-3xl text-small mr-3`" @click="goInputCount(record)">{{ $t('table.inputCount') }}</a-button>

            <a-button :class="` my-3 rounded-3xl text-small mr-3`" @click="goOutputCount(record)">{{ $t('table.outputCount') }}</a-button>
            <NuxtLink :to="localePath('/Trade')">
              <a-button :class="` my-3 rounded-3xl text-small bg-waterBlue text-white hover:!text-white`">{{ $t('table.trade') }}</a-button>
            </NuxtLink>
          </template>
        </Table>
      </div>
      <div class="md:hidden block my-8">
        <div class="bg-gray_800 rounded-xl py-1 flex items-center justify-center">
          <span class="text-normal text-white font-semibold">{{ $t('table.coinType') }}</span
          ><span class="text-small text-grayBorder ml-2">{{ $t('table.nowValue') }}</span>
        </div>
      </div>
      <div v-for="(item, index) of data" class="md:hidden block" :key="`${item.currenciesSymbol}+${index}`">
        <div class="flex justify-between border border-gray_800 rounded-xl p-3">
          <div class="flex items-center">
            <img :src="getAssetsFile(`${item.currenciesSymbol}.svg`)" class="w-10 h-10 mr-3" />
            <div class="flex flex-col items-start">
              <p class="font-bold text-subTitle text-black">{{ item.currenciesSymbol }}</p>
              <p class="text-subText text-grayText">{{ item.currenciesName }}</p>
            </div>
          </div>
          <div class="flex flex-col">
            <p v-if="eyes" class="font-bold">≈ ${{ formatValue(item.valuation) }}</p>
            <p v-else class="font-bold">********</p>
            <p class="text-small text-orange text-right">{{ item.lockedAmount }}</p>
          </div>
        </div>
        <div class="flex justify-center mb-5 mt-3">
          <a-button :class="`sm:w-[130px] xs:w-[100px] w-[85px]  my-3 rounded-3xl text-small mr-3`" @click="goInputCount(item)">{{ $t('table.inputCount') }}</a-button>
          <a-button :class="`sm:w-[130px] xs:w-[100px] w-[85px]  my-3 rounded-3xl text-small mr-3`" @click="goOutputCount(item)">{{ $t('table.outputCount') }}</a-button>

          <NuxtLink :to="localePath('/Trade')">
            <a-button :class="`sm:w-[130px] xs:w-[100px] w-[85px]   my-3 rounded-3xl text-small bg-waterBlue text-white hover:!text-white`">{{ $t('table.trade') }}</a-button>
          </NuxtLink>
        </div>
      </div>
    </div>
  </section>
  <img src="/assets/img/Group16.png" class="absolute right-0 top-0 md:w-[600px] w-[400px] z-10" />
  <img src="/assets/img/Group37.png" class="absolute left-0 top-0 md:w-[800px] w-[200px] z-10" />
</template>
<script setup>
import useAssetsStore from '@/stores/assets';
import { formatValue,getAssetsFile } from '@/config/config';

const localePath = useLocalePath();
const assetsStore = useAssetsStore();
const { t } = useI18n();
useHead({
  title: t('title.myAssets'),
  meta: [{ name: 'description', content: '' }],
});
const eyes = ref(false);
const legelCount = ref('');
const digitCount = ref('');
const totalCount = computed(() => Number(legelCount.value) + Number(digitCount.value));
const timeSelect = ref('day');
const selectTime = (value) => {
  timeSelect.value = value;
  trend.value = [];
  const trendCopy = orignalTrend.slice();
  if (value === 'week') {
    trend.value = trendCopy.slice(-7);
    trend.value = trend.value.reverse();
    max.value = Math.max(...trend.value).toFixed(0) * 1;
    min.value = Math.min(...trend.value).toFixed(0) * 1;
    trend.value.forEach((item, index) => {
      trend.value[index] = Number(item).toFixed(2);
    });
  } else if (value === 'month') {
    trend.value = trendCopy.slice(-30);
    trend.value = trend.value.reverse();
    max.value = Math.max(...trend.value).toFixed(0) * 1;
    min.value = Math.min(...trend.value).toFixed(0) * 1;
    trend.value.forEach((item, index) => {
      trend.value[index] = Number(item).toFixed(2);
    });
  } else if (value === 'year') {
    trend.value = trendCopy;
    trend.value = trend.value.reverse();
    max.value = Math.max(...trend.value).toFixed(0) * 1;
    min.value = Math.min(...trend.value).toFixed(0) * 1;
    trend.value.forEach((item, index) => {
      trend.value[index] = Number(item).toFixed(2);
    });
  } else {
    trend.value = [trendCopy[trendCopy.length - 1]];
    max.value = Number((trend.value[0] * 1.1).toFixed(0));
    min.value = Number((trend.value[0] * 0.9).toFixed(0));
    trend.value = [Number(trend.value[0]).toFixed(2)];
  }
};
const columns = [
  {
    title: 'table.coinType',
    dataIndex: 'currenciesSymbol',
    key: 'currenciesSymbol',
    width: 150,
  },
  {
    title: 'table.nowValue',
    dataIndex: 'valuation',
    key: 'valuation',
  },
  {
    title: 'table.haveValue',
    dataIndex: 'freeAmount',
    key: 'freeAmount',
    type: 'noSymbolNumber',
  },
  {
    title: 'table.handleValue',
    dataIndex: 'lockedAmount',
    key: 'lockedAmount',
    type: 'noSymbolNumber',
  },
  {
    title: 'table.action',
    dataIndex: 'operate',
    key: 'operate',
    width: 250,
  },
];
const data = ref([]);
let orignalTrend = [];
const trend = ref([]);
const max = ref(0);
const min = ref(0);
const goOutputCount = async (record) => {
  await navigateTo({
    path: localePath('/myAssets/cryptocurrencyGet'),
    query: {
      coinType: record.currenciesSymbol,
    },
  });
};
const goInputCount = async (record) => {
  await navigateTo({
    path: localePath('/myAssets/cryptocurrencyTakeOver'),
    query: {
      coinType: record.currenciesSymbol,
    },
  });
};
const getAssets = async () => {
  const result = await assetsStore.getAssets(t);
  digitCount.value = result.data.value.cryptocurrencyValuation;
  legelCount.value = result.data.value.fiatBalance;
  orignalTrend = result.data.value.histories;
  trend.value = [orignalTrend[orignalTrend.length - 1]];
  max.value = Number((trend.value[0] * 1.1).toFixed(0));
  min.value = Number((trend.value[0] * 0.9).toFixed(0));
  trend.value = [Number(trend.value[0]).toFixed(2)];
  data.value = result.data.value.assetInfos;
};
const dataReady = ref(false);
onMounted(async () => {
  await getAssets();
  dataReady.value = true;
});
</script>

<style>
.myAssetsTable .ant-table table {
  background: #ffffff !important;
}
.myAssetsTable .ant-table-thead > tr > th {
  padding: 10px 12px !important;
  background: #7589a4 !important;
  text-align: center !important;
  color: #fff !important;
  border: none !important;
}
.myAssetsTable .ant-table-tbody > tr > td {
  padding: 10px 12px !important;
  color: #181b22 !important;
  font-weight: 400 !important;
  text-align: center !important;
  color: #1d1d1d !important;
  border-bottom: 1px solid rgb(241, 241, 241) !important;
}
/* .myAssetsTable .ant-table-tbody > tr > td:last-child {
  display: flex;
  justify-content: center;
} */
.myAssetsTable .ant-table-tbody .textRight {
  text-align: right !important;
}
.myAssetsTable .ant-table-tbody > tr > td.ant-table-cell-row-hover {
  background-color: #f2f2f2 !important;
}
.myAssetsTable :where(.css-dev-only-do-not-override-eq3tly).ant-table-wrapper .ant-table-tbody > tr.ant-table-row:hover > td,
.myAssetsTable :where(.css-dev-only-do-not-override-eq3tly).ant-table-wrapper .ant-table-tbody > tr > td.ant-table-cell-row-hover {
  background: #f2f2f2 !important;
}
.myAssetsTable :where(.css-dev-only-do-not-override-eq3tly).ant-table-wrapper .ant-table-container table > thead > tr:first-child > *:first-child {
  border-end-start-radius: 10px;
}
.myAssetsTable :where(.css-dev-only-do-not-override-eq3tly).ant-table-wrapper .ant-table-container table > thead > tr:first-child > *:last-child {
  border-end-end-radius: 10px;
}
</style>
