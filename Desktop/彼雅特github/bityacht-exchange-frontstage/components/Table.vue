<template>
  <a-table :columns="columns" :data-source="data" class="px-10 py-6" :show-header="showHeader" :pagination="page" :loading="isLoading" @change="tableOnChange">
    <template #headerCell="{ column }">
      <template v-if="column">
        <span class="text-white text-content2">{{ $t(`${column.title}`) }}</span>
      </template>
    </template>
    <template #bodyCell="{ column, record, text }">
      <template v-if="column.type === 'number'">
        <p class="textRight">≈ $ {{ formatValue(text) }}</p>
      </template>
      <template v-if="column.type === 'noSymbolNumber'">
        <p class="textRight">{{ formatValueByDigits(text, 8) }}</p>
      </template>
      <template v-if="column.key === 'symbol'">
        <slot name="symbol" :record="record"></slot>
      </template>
      <template v-if="column.key === 'buyPrice'">
        <slot name="buyPrice" :record="record"></slot>
      </template>
      <template v-if="column.key === 'sellPrice'">
        <slot name="sellPrice" :record="record"></slot>
      </template>
      <template v-if="column.key === 'baseSymbol'">
        <slot name="baseSymbol" :record="record"></slot>
      </template>
      <template v-if="column.key === 'currenciesSymbol'">
        <slot name="currenciesSymbol" :record="record"></slot>
      </template>
      <template v-if="column.key === 'upsAndDowns'">
        <slot name="updown" :text="text"></slot>
      </template>
      <template v-if="column.key === 'trend'">
        <slot name="trend" :record="record"></slot>
      </template>
      <template v-if="column.key === 'nowValue'">
        <slot name="nowValue" :text="text"></slot>
      </template>
      <template v-if="column.key === 'valuation'">
        <slot name="valuation" :text="text"></slot>
      </template>
      <!-- <template v-if="column.key === 'amount'">
        <slot name="amount" :text="text"></slot>
      </template>
      <template v-if="column.key === 'handlingCharge'">
        <slot name="handlingCharge" :text="text"></slot>
      </template> -->

      <template v-if="column.key === 'status'">
        <slot name="status" :text="text"></slot>
      </template>
      <template v-if="column.key === 'side'">
        <slot name="side" :text="text"></slot>
      </template>
      <template v-if="column.key === 'action'">
        <slot name="action" :text="text"></slot>
      </template>
      <template v-if="column.key === 'txID'">
        <slot name="txid" :text="text"></slot>
      </template>
      <template v-if="column.key === 'operate'">
        <slot name="operate" :record="record"></slot>
      </template>
      <template v-if="column.key === 'handlingCharge'">
        <slot name="handlingCharge" :record="record"></slot>
      </template>
    </template>
    <template #emptyText>
      <div class="flex justify-center items-center mt-10">
        <img src="/assets/img/noData.png" />
        <p class="text-h3 text-gray_500 ml-2">{{ $t('table.noData') }}</p>
      </div>
    </template>
  </a-table>
</template>
<script setup>
import useAssetsStore from '@/stores/assets';
import { formatValue, formatValueByDigits } from '@/config/config';

const { data, columns } = defineProps({
  columns: {
    type: Array,
  },
  data: {
    type: Array || Object,
  },
  showHeader: {
    type: Boolean,
    default: true,
  },
  page: {
    type: Object,
  },
});
const emit = defineEmits(['changePage']);

const tableOnChange = (pagination) => {
  emit('changePage', pagination);
};
const assetsStore = useAssetsStore();
const isLoading = computed(() => assetsStore.isLoadingTable);
</script>
