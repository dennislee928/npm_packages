<template>
  <section class="relative 3xl:px-[300px] xl:px-[200px] lg:px-[100px] px-[30px] bg-white pb-10">
    <h1 class="relative z-30 flex items-center justify-center md:text-h2 text-h4 text-gray_800 pt-[88px] md:pb-[60px] pb-[40px]"><left-circle-outlined class="mr-4 text-h5" @click="back" />{{ $t('myAssets.walletHistoryRecord') }}</h1>
    <div class="relative z-30 flex justify-center md:mb-[60px] mb-6">
      <a-button :class="`md:w-[145px] w-[100px] md:text-subTitle text-subText md:py-5 py-3  rounded-3xl flex items-center justify-center md:mr-5 mr-4 ${selectCoinType === 'cryptocurrency' ? 'bg-waterBlue text-white hover:!text-white' : ''}`" @click="coinTypeSelect('cryptocurrency')">{{ $t('myAssets.cryptocurrency') }}</a-button>
      <a-button :class="`md:w-[145px] w-[100px] md:text-subTitle text-subText md:py-5 py-3 rounded-3xl flex items-center justify-center md:mr-0 mr-4 ${selectCoinType === 'legalTender' ? 'bg-waterBlue text-white hover:!text-white' : ''}`" :disabled="true" @click="coinTypeSelect('legalTender')">{{ $t('myAssets.legalTender') }}</a-button>
      <img src="/icon/filter.svg" class="cursor-pointer md:hidden block" @click="mobileDialog = !mobileDialog" />
    </div>
    <a-drawer v-model:open="mobileDialog" placement="bottom" :closable="false" height="450" class="relative border border-waterBlue rounded-xl px-5">
      <close-circle-outlined class="text-gray_800 text-[26px] absolute top-5 right-5" @click="mobileDialog = !mobileDialog"></close-circle-outlined>
      <h1 class="xs:text-subTitle text-normal font-bold text-center">{{ $t('table.filter') }}</h1>
      <div class="mt-4 flex flex-col">
        <div class="flex flex-col mb-6">
          <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('table.tradeTime') }}</p>
          <a-range-picker v-model:value="selectTradeTime" :format="tradeTimeFormat" :locale="locale === 'zh-TW' ? localeTW : ''"></a-range-picker>
        </div>
        <div class="flex flex-col mb-6">
          <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('table.status') }}</p>
          <FormAntdSelect v-model="selectStatus">
            <a-select-option :value="0">{{ $t('myAssets.all') }}</a-select-option>
            <a-select-option :value="1">{{ $t('myAssets.handle') }}</a-select-option>
            <a-select-option :value="2">{{ $t('myAssets.finish') }}</a-select-option>
            <a-select-option :value="3">{{ $t('myAssets.fail') }}</a-select-option>
            <a-select-option :value="4">{{ $t('myAssets.cancel') }}</a-select-option>
            <a-select-option :value="5">{{ $t('myAssets.checking') }}</a-select-option>
          </FormAntdSelect>
        </div>
        <div class="flex flex-col mb-6">
          <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('table.kind') }}</p>
          <FormAntdSelect v-model="selectKind">
            <a-select-option :value="1">{{ $t('myAssets.inputCoin') }}</a-select-option>
            <a-select-option :value="2">{{ $t('myAssets.outputCoin') }}</a-select-option>
          </FormAntdSelect>
        </div>
        <div class="flex flex-col mb-[50px]">
          <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('table.coinType') }}</p>
          <FormAntdSelect v-model="selectCoin">
            <a-select-option value="">{{ $t('myAssets.all') }}</a-select-option>
            <a-select-option value="BTC">BTC</a-select-option>
            <a-select-option value="ETH">ETH</a-select-option>
            <a-select-option value="USDC">USDC</a-select-option>
            <a-select-option value="USDT">USDT</a-select-option>
          </FormAntdSelect>
        </div>
        <div class="flex justify-center">
          <a-button :class="`md:w-[145px] w-[100px] md:text-subTitle text-subText md:py-5 py-3 rounded-3xl flex items-center justify-center md:mr-5 mr-4 `" @click="reset">{{ $t('table.resetFilter') }}</a-button>
          <!-- <a-button :class="`md:w-[145px] w-[100px] md:text-subTitle text-subText md:py-5 py-3 rounded-3xl flex items-center justify-center md:mr-5 mr-4 bg-waterBlue text-white active:!text-white hover:!text-white`">{{ $t('table.query') }}</a-button> -->
        </div>
      </div>
    </a-drawer>
    <div class="relative z-30 selectBox md:flex hidden justify-between">
      <div class="flex flex-col items-center">
        <p class="text-normal mb-2">{{ $t('myAssets.status') }}</p>
        <FormAntdSelect v-model="selectStatus" class="2xl:w-[225px] xl:w-[180px] lg:w-[150px] md:w-[125px] text-center">
          <a-select-option :value="0">{{ $t('myAssets.all') }}</a-select-option>
          <a-select-option :value="1">{{ $t('myAssets.handle') }}</a-select-option>
          <a-select-option :value="2">{{ $t('myAssets.finish') }}</a-select-option>
          <a-select-option :value="3">{{ $t('myAssets.fail') }}</a-select-option>
          <a-select-option :value="4">{{ $t('myAssets.cancel') }}</a-select-option>
          <a-select-option :value="5">{{ $t('myAssets.checking') }}</a-select-option>
        </FormAntdSelect>
      </div>
      <div class="flex flex-col items-center">
        <p class="text-normal mb-2">{{ $t('myAssets.kind') }}</p>
        <FormAntdSelect v-model="selectKind" class="2xl:w-[225px] xl:w-[180px] lg:w-[150px] md:w-[125px] text-center">
          <a-select-option :value="1">{{ $t('myAssets.inputCoin') }}</a-select-option>
          <a-select-option :value="2">{{ $t('myAssets.outputCoin') }}</a-select-option>
        </FormAntdSelect>
      </div>
      <div class="flex flex-col items-center">
        <p class="text-normal mb-2">{{ $t('myAssets.coinType') }}</p>
        <FormAntdSelect v-model="selectCoin" class="2xl:w-[225px] xl:w-[180px] lg:w-[150px] md:w-[125px] text-center">
          <a-select-option value="">{{ $t('myAssets.all') }}</a-select-option>
          <a-select-option value="BTC">BTC</a-select-option>
          <a-select-option value="ETH">ETH</a-select-option>
          <a-select-option value="USDC">USDC</a-select-option>
          <a-select-option value="USDT">USDT</a-select-option>
        </FormAntdSelect>
      </div>
      <div class="flex flex-col items-center">
        <p class="text-normal mb-2">{{ $t('myAssets.tradeTime') }}</p>
        <a-range-picker v-model:value="selectTradeTime" :format="tradeTimeFormat" :locale="locale === 'zh-TW' ? localeTW : ''"></a-range-picker>
      </div>
    </div>
    <div class="relative z-30 walletRecordTable md:block hidden border border-gray_800 rounded-xl mt-10">
      <Table :columns="columns" :data="data" :page="tableConfig" @change-page="changePage">
        <template v-slot:txid="{ text }">
          <a-tooltip v-if="text !== ''">
            <template #title>{{ text }}</template>
            <span class="text-subText">{{ shortString(text) }}</span>
            <span class="text-small text-gray_800 bg-gray_300 rounded-lg cursor-pointer p-1 ml-2" @click="copyText(text, t)">{{ $t('table.copy') }}</span>
          </a-tooltip>
        </template>
        <!-- <template v-slot:amount="{ text }">
          <span class="text-right">{{ formatValue(text) }}</span>
        </template>
        <template v-slot:handlingCharge="{ text }">
          <span class="text-right">{{ formatValue(text) }}</span>
        </template> -->
        <template v-slot:status="{ text }">
          <span v-if="text === 0">{{ $t('myAssets.all') }}</span>
          <span v-if="text === 1">{{ $t('myAssets.handle') }}</span>
          <span v-if="text === 2">{{ $t('myAssets.finish') }}</span>
          <span v-if="text === 3">{{ $t('myAssets.fail') }}</span>
          <span v-if="text === 4">{{ $t('myAssets.cancel') }}</span>
          <span v-if="text === 5">{{ $t('myAssets.checking') }}</span>
        </template>
        <template v-slot:action="{ text }">
          <span v-if="text === 1">{{ $t('myAssets.inputCoin') }}</span>
          <span v-if="text === 2">{{ $t('myAssets.outputCoin') }}</span>
        </template>
      </Table>
    </div>
    <div class="relative z-30 md:hidden block border border-gray_800 rounded-xl mt-5 xs:p-10 p-6">
      <template v-if="data.length === 0">
        <div class="flex justify-center items-center flex-col">
          <img src="/assets/img/noData.png" />
          <p class="text-[26px] text-gray_500 mt-5">{{ $t('table.noData') }}</p>
        </div>
      </template>
      <template v-else>
        <div v-for="item of data" class="border-b border-gray_300 mb-[55px]" :key="item.transfersID">
          <div class="flex justify-between mb-2">
            <p class="xxs:text-normal text-subText font-bold">{{ $t('table.orderNumber') }}</p>
            <p class="xxs:text-normal text-subText font-bold">{{ item.transfersID }}</p>
          </div>
          <div class="flex justify-between mb-2">
            <p class="xxs:text-normal text-subText">{{ $t('table.status') }}</p>
            <template v-if="item.status">
              <p v-if="item.status === 0" class="xxs:text-normal text-subText">{{ $t('myAssets.all') }}</p>
              <p v-else-if="item.status === 1" class="xxs:text-normal text-subText">{{ $t('myAssets.handle') }}</p>
              <p v-else-if="item.status === 2" class="xxs:text-normal text-subText">{{ $t('myAssets.finish') }}</p>
              <p v-else-if="item.status === 3" class="xxs:text-normal text-subText">{{ $t('myAssets.fail') }}</p>
              <p v-else-if="item.status === 4" class="xxs:text-normal text-subText">{{ $t('myAssets.cancel') }}</p>
              <p v-else class="xxs:text-normal text-subText">{{ $t('myAssets.checking') }}</p>
            </template>
          </div>
          <div class="flex justify-between mb-2">
            <p class="xxs:text-normal text-subText">{{ $t('table.txid') }}</p>
            <p class="xxs:text-normal text-subText">{{ shortString(item.txID) }}</p>
          </div>
          <div class="flex justify-between mb-2">
            <p class="xxs:text-normal text-subText">{{ $t('table.kind') }}</p>
            <template v-if="item.action">
              <p v-if="item.action === 1" class="xxs:text-normal text-subText">{{ $t('myAssets.inputCoin') }}</p>
              <p v-else class="xxs:text-normal text-subText">{{ $t('myAssets.outputCoin') }}</p>
            </template>
          </div>
          <div class="flex justify-between mb-2">
            <p class="xxs:text-normal text-subText">{{ $t('table.coinType') }}</p>
            <p class="xxs:text-normal text-subText">{{ item.currenciesSymbol }}</p>
          </div>
          <div class="flex justify-between mb-2">
            <p class="xxs:text-normal text-subText">{{ $t('table.count') }}</p>
            <p class="xxs:text-normal text-subText">{{ formatValue(item.amount) }}</p>
          </div>
          <div class="flex justify-between mb-2">
            <p class="xxs:text-normal text-subText">{{ $t('table.redemptionFee') }}</p>
            <p class="xxs:text-normal text-subText">{{ formatValue(item.handlingCharge) }}</p>
          </div>
          <div class="flex justify-between mb-5">
            <p class="xxs:text-normal text-subText">{{ $t('table.tradeTime') }}</p>
            <p class="xxs:text-normal text-subText">{{ item.createdAt }}</p>
          </div>
        </div>
        <div class="text-center mt-5">
          <a-pagination v-model:current="tableConfig.current" v-model:pageSize="tableConfig.pageSize" :total="tableConfig.total" :showSizeChanger="true" :pageSizeOptions="['25', '50', '100']" @change="pageChange" />
        </div>
      </template>
    </div>
  </section>
  <img src="/assets/img/Group16.png" class="absolute right-0 top-0 md:w-[600px] w-[400px] z-10" />
  <img src="/assets/img/Group37.png" class="absolute left-0 top-0 md:w-[800px] w-[200px] z-10" />
</template>
<script setup>
import useAssetsStore from '@/stores/assets';
import { copyText, shortString, formatValue } from '@/config/config';
import localeTW from '@/node_modules/ant-design-vue/es/date-picker/locale/zh_TW';

const localePath = useLocalePath();
const assetsStore = useAssetsStore();
const dayjs = useDayjs();
const { t, locale } = useI18n();
dayjs.locale('zh-tw');
useHead({
  title: t('title.myAssets_walletHistoryRecord'),
  meta: [{ name: 'description', content: '' }],
});
const columns = [
  {
    title: 'table.orderNumber',
    dataIndex: 'transfersID',
    key: 'transfersID',
    width: 150,
  },
  {
    title: 'table.status',
    dataIndex: 'status',
    key: 'status',
  },
  {
    title: 'table.txid',
    dataIndex: 'txID',
    key: 'txID',
  },
  {
    title: 'table.kind',
    dataIndex: 'action',
    key: 'action',
  },
  {
    title: 'table.coinType',
    dataIndex: 'currenciesSymbol',
    key: 'currenciesSymbol',
  },
  {
    title: 'table.count',
    dataIndex: 'amount',
    key: 'amount',
    type: 'noSymbolNumber',
  },
  // {
  //   title: 'table.redemptionFee',
  //   dataIndex: 'handlingCharge',
  //   key: 'handlingCharge',
  //   type: 'noSymbolNumber',
  // },
  {
    title: 'table.tradeTime',
    dataIndex: 'createdAt',
    key: 'createdAt',
    width: 200,
  },
];
const data = ref([]);
const tableConfig = ref({
  position: ['bottomCenter'],
  showSizeChanger: true,
  current: 1,
  pageSize: 25,
  total: 0,
  pageSizeOptions: ['25', '50', '100'],
});
const selectCoinType = ref('cryptocurrency');
const coinTypeSelect = (type) => {
  selectCoinType.value = type;
};
const selectStatus = ref(0); // 0:all,1:handle,2:finish,3:fail,4:cancel
const selectKind = ref(1);
const selectCoin = ref('');
const selectTradeTime = ref([]);
const tradeTimeFormat = ['YYYY/MM/DD', 'YYYY/MM/DD'];
const mobileDialog = ref(false);
const back = async () => {
  await navigateTo(localePath('/myAssets'));
};

watch(
  () => [selectStatus.value, selectKind.value, selectCoin.value, selectTradeTime.value],
  async () => {
    assetsStore.isLoadingTable = true;
    let day = [];
    if (!selectTradeTime.value || !selectTradeTime.value.length) {
      day = ['', ''];
    } else {
      day = selectTradeTime.value.map((date) => dayjs(date).format('YYYY/MM/DD'));
    }
    tableConfig.value.current = 1;
    const result = await assetsStore.getHistories(selectKind.value, selectCoin.value, selectStatus.value, `"${day[0]}"`, `"${day[1]}"`, tableConfig.value.current, tableConfig.value.pageSize, t);
    data.value = result.data.value.data;
    tableConfig.value.current = result.data.value.page;
    tableConfig.value.pageSize = result.data.value.pageSize;
    tableConfig.value.total = result.data.value.totalRecord;
    assetsStore.isLoadingTable = false;
  },
  { deep: true }
);
const changePage = async (paginations) => {
  assetsStore.isLoadingTable = true;
  let day = [];
  if (!selectTradeTime.value || !selectTradeTime.value.length) {
    day = ['', ''];
  } else {
    day = selectTradeTime.value.map((date) => dayjs(date).format('YYYY/MM/DD'));
  }
  const result = await assetsStore.getHistories(selectKind.value, selectCoin.value, selectStatus.value, `"${day[0]}"`, `"${day[1]}"`, paginations.current, paginations.pageSize, t);
  data.value = result.data.value.data;
  tableConfig.value.current = result.data.value.page;
  tableConfig.value.pageSize = result.data.value.pageSize;
  tableConfig.value.total = result.data.value.totalRecord;
  assetsStore.isLoadingTable = false;
  window.scrollTo(0, 0);
};
const pageChange = async (page, pageSize) => {
  assetsStore.isLoadingTable = true;
  let day = [];
  if (!selectTradeTime.value || !selectTradeTime.value.length) {
    day = ['', ''];
  } else {
    day = selectTradeTime.value.map((date) => dayjs(date).format('YYYY/MM/DD'));
  }
  const result = await assetsStore.getHistories(selectKind.value, selectCoin.value, selectStatus.value, `"${day[0]}"`, `"${day[1]}"`, page, pageSize, t);
  data.value = result.data.value.data;
  tableConfig.value.current = result.data.value.page;
  tableConfig.value.pageSize = result.data.value.pageSize;
  tableConfig.value.total = result.data.value.totalRecord;
  assetsStore.isLoadingTable = false;
  window.scrollTo(0, 0);
};
const reset = () => {
  selectStatus.value = 0;
  selectKind.value = 1;
  selectCoin.value = '';
  selectTradeTime.value = [];
  mobileDialog.value = false;
};
const getHistories = async () => {
  assetsStore.isLoadingTable = true;
  const result = await assetsStore.getHistories(selectKind.value, selectCoin.value, selectStatus.value, undefined, undefined, undefined, undefined, t);
  data.value = result.data.value.data;
  tableConfig.value.total = result.data.value.totalRecord;
  assetsStore.isLoadingTable = false;
};
onMounted(() => {
  getHistories();
});
</script>
<style>
.selectBox .ant-select-selector {
  background-color: #efefef !important;
  border-radius: 25px;
  padding: 25px 0px !important;
  display: flex;
  align-items: center;
}
.selectBox :where(.css-dev-only-do-not-override-eq3tly).ant-select-single.ant-select-show-arrow .ant-select-selection-item {
  padding-inline-end: 0px !important;
}
.selectBox .ant-picker {
  background-color: #efefef !important;
  border-radius: 25px;
  padding: 14px !important;
}
.selectBox :where(.css-dev-only-do-not-override-eq3tly).ant-picker-range .ant-picker-range-separator {
  padding: 0px 8px 5px 8px !important;
}
.walletRecordTable .ant-table-wrapper {
  padding: 24px !important;
}
.walletRecordTable .ant-table table {
  background: #ffffff !important;
}
.walletRecordTable .ant-table-thead > tr > th {
  padding: 10px 12px !important;
  background: #7589a4 !important;
  text-align: center !important;
  color: #fff !important;
  border: none !important;
}
.walletRecordTable .ant-table-tbody > tr > td {
  padding: 14px !important;
  color: #181b22 !important;
  font-weight: 400 !important;
  text-align: center !important;
  color: #1d1d1d !important;
  border: none !important;
}
/* .walletRecordTable .ant-table-tbody > tr > td:last-child {
  display: flex;
  justify-content: center;
} */
.walletRecordTable .ant-table-tbody .textRight {
  text-align: right !important;
}
.walletRecordTable .ant-table-tbody > tr > td.ant-table-cell-row-hover {
  background-color: #f2f2f2 !important;
}
.walletRecordTable :where(.css-dev-only-do-not-override-eq3tly).ant-table-wrapper .ant-table-tbody > tr.ant-table-row:hover > td,
.walletRecordTable :where(.css-dev-only-do-not-override-eq3tly).ant-table-wrapper .ant-table-tbody > tr > td.ant-table-cell-row-hover {
  background: #f2f2f2 !important;
}
.walletRecordTable :where(.css-dev-only-do-not-override-eq3tly).ant-table-wrapper .ant-table-container table > thead > tr:first-child > *:first-child {
  border-end-start-radius: 10px;
}
.walletRecordTable :where(.css-dev-only-do-not-override-eq3tly).ant-table-wrapper .ant-table-container table > thead > tr:first-child > *:last-child {
  border-end-end-radius: 10px;
}
</style>
