<script setup>
import useOrderStatus from "@/hooks/useOrderStatus";
const { status, badgeStatus } = useOrderStatus();
const { data, getData } = defineProps({
  data: {
    type: Object,
  },
  getData: {
    type: Function,
  },
});
</script>
<template>
  <TableContainer :cols="`8em 4em 3em 4em 4em 6em 1.2fr 1fr 1fr`">
    <TableHeader>
      <div>交易編號</div>
      <div>狀態</div>
      <div>類型</div>
      <div>幣種</div>
      <div>主網</div>
      <div class="text-right">數量</div>
      <div class="text-center">TXID</div>
      <div>交易時間</div>
      <div>完成時間</div>
    </TableHeader>
    <TableItem v-for="item in data.informations.spotTransfers" :key="item">
      <div>{{ item.id }}</div>
      <div>
        <StatusBadge :status="badgeStatus[item.status]" variant="filled">
          {{ status[item.status] }}
        </StatusBadge>
      </div>
      <div>{{ item.action === 1 ? "入幣" : "提幣" }}</div>
      <div>{{ item.currenciesSymbol }}</div>
      <div>{{ item.mainnet }}</div>
      <ParsePrice :price="item.amount" />
      <ParseTxid :txid="item.txID" />
      <div>{{ item.createdAt }}</div>
      <div>{{ item.finishedAt }}</div>
    </TableItem>
  </TableContainer>
</template>
