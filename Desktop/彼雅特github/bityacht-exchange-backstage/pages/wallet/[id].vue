<script setup>
import useOrderStatus from "@/hooks/useOrderStatus";

const route = useRoute();
const { id } = route.params;
const { status, badgeStatus } = useOrderStatus();

const updateStatusDialog = ref(false);
</script>
<template>
  <PageTitle>
    <div class="flex gap-2 items-center justify-start">
      <button
        @click="$router.back()"
        class="bg-[#394B6A] text-white py-1 px-3 rounded-full hover:opacity-80 active:opacity-60 text-base font-normal"
      >
        ❮ 返回
      </button>
      提/入幣列表
    </div>
  </PageTitle>

  <TableContainer cols="1fr 1fr">
    <TableHeader>
      <div>訂單編號</div>
    </TableHeader>
    <TableItemDual>
      <TableSubItem>
        <div>交易編號</div>
        <div>{{ id }}</div>
      </TableSubItem>
      <TableSubItem>
        <div>狀態</div>
        <div class="flex justify-between items-center">
          <StatusBadge status="warn" variant="filled"> 可疑單 </StatusBadge>
          <Button size="mini" @click="updateStatusDialog = { status: 1 }">
            編輯
          </Button>
        </div>
      </TableSubItem>
      <TableSubItem>
        <div>方向</div>
        <div>提幣</div>
      </TableSubItem>
      <TableSubItem>
        <div>幣種</div>
        <div>USDT</div>
      </TableSubItem>
      <TableSubItem>
        <div>主網</div>
        <div>TRC20</div>
      </TableSubItem>
      <TableSubItem>
        <div>數量</div>
        <div>1500</div>
      </TableSubItem>
      <TableSubItem>
        <div>TXID</div>
        <div>-</div>
      </TableSubItem>
      <TableSubItem>
        <div>UID</div>
        <div>1500</div>
      </TableSubItem>
      <TableSubItem>
        <div>交易時間</div>
        <div>2023/08/17 13:30:33</div>
      </TableSubItem>
      <TableSubItem>
        <div>完成時間</div>
        <div>-</div>
      </TableSubItem>
      <TableSubItem>
        <div>交易IP位罝</div>
        <div>123.220.133.22</div>
      </TableSubItem>
      <TableSubItem>
        <div>交易IP國家</div>
        <div>台灣</div>
      </TableSubItem>
    </TableItemDual>
  </TableContainer>

  <TableContainer cols="220px 1fr 5em">
    <TableHeader>
      <div>訂單編號</div>
    </TableHeader>
    <TableItem>
      <div>2023/08/17 10:40:23</div>
      <div class="link">密集交易（230817003）</div>
      <div class="text-right">待審核</div>
    </TableItem>
  </TableContainer>
  <Dialog v-model="updateStatusDialog" title="狀態審核" :loading="loading">
    <form
      @submit.prevent
      class="grid grid-cols-[4rem_1fr] gap-4 items-center justify-end"
      id="update-status-form"
      v-if="updateStatusDialog"
    >
      <FormLabel label="狀態" />
      <a-radio-group
        v-model:value="updateStatusDialog.status"
        :options="
          status.slice(1).map((item, i) => ({
            label: item,
            value: i + 1,
          }))
        "
        required
      />
    </form>
    <template #actions>
      <Button variant="outline" @click="updateStatusDialog = false">
        取消
      </Button>
      <Button @click="updateStatusDialog = false"> 確認 </Button>
    </template>
  </Dialog>
</template>
