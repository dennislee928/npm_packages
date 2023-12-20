<template>
  <section class="relative 3xl:px-[300px] xl:px-[200px] lg:px-[100px] px-[30px] bg-white pb-10">
    <h1 class="relative z-30 flex items-center justify-center md:text-h2 text-h4 text-gray_800 pt-[88px] md:pb-[60px] pb-[40px]"><left-circle-outlined class="mr-4 text-h5" @click="back" />{{ $t('trade.tradeRecord') }}</h1>
    <div class="relative z-30 flex justify-center md:mb-[60px] mb-6">
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
            <!-- <a-select-option :value="1">{{ $t('myAssets.handle') }}</a-select-option> -->
            <a-select-option :value="1">{{ $t('myAssets.finish') }}</a-select-option>
            <a-select-option :value="2">{{ $t('myAssets.fail') }}</a-select-option>
            <!-- <a-select-option :value="4">{{ $t('myAssets.cancel') }}</a-select-option> -->
          </FormAntdSelect>
        </div>
        <div class="flex flex-col mb-[50px]">
          <p class="xs:text-subText text-small font-bold flex items-center mb-2"><img src="/icon/hexagon.svg" class="mr-2 w-[20px]" />{{ $t('table.side') }}</p>
          <FormAntdSelect v-model="selectSide">
            <a-select-option :value="1">{{ $t('trade.buy') }}</a-select-option>
            <a-select-option :value="2">{{ $t('trade.sale') }}</a-select-option>
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
        <FormAntdSelect v-model="selectStatus" class="2xl:w-[300px] xl:w-[270px] md:w-[200px] text-center">
          <a-select-option :value="0">{{ $t('myAssets.all') }}</a-select-option>
          <!-- <a-select-option :value="1">{{ $t('myAssets.handle') }}</a-select-option> -->
          <a-select-option :value="1">{{ $t('myAssets.finish') }}</a-select-option>
          <a-select-option :value="2">{{ $t('myAssets.fail') }}</a-select-option>
          <!-- <a-select-option :value="4">{{ $t('myAssets.cancel') }}</a-select-option> -->
        </FormAntdSelect>
      </div>
      <div class="flex flex-col items-center">
        <p class="text-normal mb-2">{{ $t('table.side') }}</p>
        <FormAntdSelect v-model="selectSide" class="2xl:w-[300px] xl:w-[270px] md:w-[200px] text-center">
          <a-select-option :value="1">{{ $t('trade.buy') }}</a-select-option>
          <a-select-option :value="2">{{ $t('trade.sale') }}</a-select-option>
        </FormAntdSelect>
      </div>
      <div class="flex flex-col items-center">
        <p class="text-normal mb-2">{{ $t('myAssets.tradeTime') }}</p>
        <a-range-picker v-model:value="selectTradeTime" :format="tradeTimeFormat" :locale="locale === 'zh-TW' ? localeTW : ''"></a-range-picker>
      </div>
    </div>
    <div class="relative z-30 tradeRecordTable md:block hidden border border-gray_800 rounded-xl mt-10">
      <Table :columns="columns" :data="data" :page="tableConfig" @change-page="changePage">
        <template v-slot:baseSymbol="{ record }">
          <span class="text-subText">{{ record.baseSymbol }} / </span>
          <span class="text-subText">{{ record.quoteSymbol }}</span>
        </template>
        <template v-slot:status="{ text }">
          <span v-if="text === 0">{{ $t('myAssets.all') }}</span>
          <!-- <span v-if="text === 1">{{ $t('myAssets.handle') }}</span> -->
          <span v-if="text === 1">{{ $t('myAssets.finish') }}</span>
          <span v-if="text === 2">{{ $t('myAssets.fail') }}</span>
          <!-- <span v-if="text === 4">{{ $t('myAssets.cancel') }}</span> -->
        </template>
        <template v-slot:side="{ text }">
          <span v-if="text === 1">{{ $t('trade.buy') }}</span>
          <span v-if="text === 2">{{ $t('trade.sale') }}</span>
        </template>
        <template v-slot:handlingCharge="{ record }">
          <span>{{ record.handlingCharge }}</span>
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
        <div v-for="item of data" class="border-b border-gray_300 mb-[55px]" :key="item.transactionsID">
          <div class="flex justify-between mb-2">
            <p class="xxs:text-normal text-subText font-bold">{{ $t('table.orderNumber') }}</p>
            <p class="xxs:text-normal text-subText font-bold">{{ item.transactionsID }}</p>
          </div>
          <div class="flex justify-between mb-2">
            <p class="xxs:text-normal text-subText">{{ $t('table.status') }}</p>
            <template v-if="item.status">
              <p v-if="item.status === 0" class="xxs:text-normal text-subText">{{ $t('myAssets.all') }}</p>
              <!-- <p v-else-if="item.status === 1" class="xxs:text-normal text-subText">{{ $t('myAssets.handle') }}</p> -->
              <p v-else-if="item.status === 1" class="xxs:text-normal text-subText">{{ $t('myAssets.finish') }}</p>
              <p v-else-if="item.status === 2" class="xxs:text-normal text-subText">{{ $t('myAssets.fail') }}</p>
              <!-- <p v-else class="xxs:text-normal text-subText">{{ $t('myAssets.cancel') }}</p> -->
            </template>
          </div>
          <div class="flex justify-between mb-2">
            <p class="xxs:text-normal text-subText">{{ $t('table.tradeCouple') }}</p>
            <p>
              <span class="xxs:text-normal text-subText">{{ item.baseSymbol }} / </span>
              <span class="xxs:text-normal text-subText">{{ item.quoteSymbol }}</span>
            </p>
          </div>
          <div class="flex justify-between mb-2">
            <p class="xxs:text-normal text-subText">{{ $t('table.side') }}</p>
            <template v-if="item.side">
              <p v-if="item.side === 1" class="xxs:text-normal text-subText">{{ $t('trade.buy') }}</p>
              <p v-else class="xxs:text-normal text-subText">{{ $t('trade.sale') }}</p>
            </template>
          </div>
          <div class="flex justify-between mb-2">
            <p class="xxs:text-normal text-subText">{{ $t('table.count') }}</p>
            <p class="xxs:text-normal text-subText">{{ formatValue(item.quantity) }}</p>
          </div>
          <div class="flex justify-between mb-2">
            <p class="xxs:text-normal text-subText">{{ $t('table.price') }}</p>
            <p class="xxs:text-normal text-subText">{{ formatValue(item.price) }}</p>
          </div>
          <div class="flex justify-between mb-2">
            <p class="xxs:text-normal text-subText">{{ $t('table.amount') }}</p>
            <p class="xxs:text-normal text-subText">{{ formatValue(item.amount) }}</p>
          </div>
          <div class="flex justify-between mb-2">
            <p class="xxs:text-normal text-subText">{{ $t('table.redemptionFee') }}</p>
            <p class="xxs:text-normal text-subText">{{ item.handlingCharge }}</p>
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
import { shortString, formatValue } from '@/config/config';
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
    dataIndex: 'transactionsID',
    key: 'transactionsID',
    width: 150,
  },
  {
    title: 'table.status',
    dataIndex: 'status',
    key: 'status',
  },
  {
    title: 'table.tradeCouple',
    dataIndex: 'baseSymbol',
    key: 'baseSymbol',
  },
  {
    title: 'table.side',
    dataIndex: 'side',
    key: 'side',
  },
  {
    title: 'table.count',
    dataIndex: 'quantity',
    key: 'quantity',
    type: 'noSymbolNumber',
  },
  {
    title: 'table.price',
    dataIndex: 'price',
    key: 'price',
    type: 'noSymbolNumber',
  },
  {
    title: 'table.amount',
    dataIndex: 'amount',
    key: 'amount',
    type: 'noSymbolNumber',
  },
  {
    title: 'table.redemptionFee',
    dataIndex: 'handlingCharge',
    key: 'handlingCharge',
  },
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
const selectStatus = ref(0); // 0:all,1:handle,2:finish,3:fail,4:cancel
const selectSide = ref(1);
const selectTradeTime = ref([]);
const tradeTimeFormat = ['YYYY/MM/DD', 'YYYY/MM/DD'];
const mobileDialog = ref(false);
const back = async () => {
  await navigateTo(localePath('/Trade'));
};

watch(
  () => [selectStatus.value, selectSide.value, selectTradeTime.value],
  async (newVal) => {
    assetsStore.isLoadingTable = true;
    let day = [];
    if (!selectTradeTime.value || !selectTradeTime.value.length) {
      day = ['', ''];
    } else {
      day = selectTradeTime.value.map((date) => dayjs(date).format('YYYY/MM/DD'));
    }
    tableConfig.value.current = 1;
    const result = await assetsStore.getTransactions(selectSide.value, selectStatus.value, `"${day[0]}"`, `"${day[1]}"`, tableConfig.value.current, tableConfig.value.pageSize, t);
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
  const result = await assetsStore.getTransactions(selectSide.value, selectStatus.value, `"${day[0]}"`, `"${day[1]}"`, paginations.current, paginations.pageSize, t);
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
  const result = await assetsStore.getTransactions(selectSide.value, selectStatus.value, `"${day[0]}"`, `"${day[1]}"`, page, pageSize, t);
  data.value = result.data.value.data;
  tableConfig.value.current = result.data.value.page;
  tableConfig.value.pageSize = result.data.value.pageSize;
  tableConfig.value.total = result.data.value.totalRecord;
  assetsStore.isLoadingTable = false;
  window.scrollTo(0, 0);
};
const reset = () => {
  selectStatus.value = 0;
  selectSide.value = 1;
  selectTradeTime.value = [];
  mobileDialog.value = false;
};
const getTransactions = async () => {
  assetsStore.isLoadingTable = true;
  const result = await assetsStore.getTransactions(selectSide.value, selectStatus.value, undefined, undefined, undefined, undefined, t);
  data.value = result.data.value.data;
  tableConfig.value.total = result.data.value.totalRecord;
  assetsStore.isLoadingTable = false;
};
onMounted(async () => {
  await getTransactions();
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
.tradeRecordTable .ant-table-wrapper {
  padding: 24px !important;
}
.tradeRecordTable .ant-table table {
  background: #ffffff !important;
}
.tradeRecordTable .ant-table-thead > tr > th {
  padding: 10px 12px !important;
  background: #7589a4 !important;
  text-align: center !important;
  color: #fff !important;
  border: none !important;
}
.tradeRecordTable .ant-table-tbody > tr > td {
  padding: 14px !important;
  color: #181b22 !important;
  font-weight: 400 !important;
  text-align: center !important;
  color: #1d1d1d !important;
  border: none !important;
}
/* .tradeRecordTable .ant-table-tbody > tr > td:last-child {
  display: flex;
  justify-content: center;
} */
.tradeRecordTable .ant-table-tbody .textRight {
  text-align: right !important;
}
.tradeRecordTable .ant-table-tbody > tr > td.ant-table-cell-row-hover {
  background-color: #f2f2f2 !important;
}
.tradeRecordTable :where(.css-dev-only-do-not-override-eq3tly).ant-table-wrapper .ant-table-tbody > tr.ant-table-row:hover > td,
.tradeRecordTable :where(.css-dev-only-do-not-override-eq3tly).ant-table-wrapper .ant-table-tbody > tr > td.ant-table-cell-row-hover {
  background: #f2f2f2 !important;
}
.tradeRecordTable :where(.css-dev-only-do-not-override-eq3tly).ant-table-wrapper .ant-table-container table > thead > tr:first-child > *:first-child {
  border-end-start-radius: 10px;
}
.tradeRecordTable :where(.css-dev-only-do-not-override-eq3tly).ant-table-wrapper .ant-table-container table > thead > tr:first-child > *:last-child {
  border-end-end-radius: 10px;
}
</style>
