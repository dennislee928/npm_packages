<script setup>
import useAPI from "@/hooks/useAPI";
import usePagination from "@/hooks/usePagination";
import { FileTextOutlined } from "@ant-design/icons-vue";
import useUserInfo from "@/hooks/useUserInfo";
import useParseDate from "@/hooks/useParseDate";

const { id, showHeader } = defineProps({
  id: {
    type: String,
    default: "0",
  },
  showHeader: {
    type: Boolean,
    default: false,
  },
});
const { getTransactions, getExportTransactionsUrl } = useAPI();
const { data, page, perPage, totalPages, totalRecord, updatePagination } =
  usePagination();
const { userInfo } = useUserInfo();
const parseDate = useParseDate();

const form = ref({
  status: "0",
  side: "0",
  startAt: "",
  endAt: "",
  search: "",
});
const loading = ref(false);
const exportDialog = ref(false);
const orderIdDialog = ref(false);
const invoiceDialog = ref(false);

onMounted(async () => {
  getData();
});
watch(
  [page, perPage, form],
  () => {
    getData();
  },
  { deep: true }
);
async function getData() {
  loading.value = true;
  let options = {};
  for (let [key, value] of Object.entries(form.value)) {
    if (value !== "-1" && value !== "") {
      options[key] = value;
    }
  }
  if (options?.startAt || options?.endAt) {
    if (!options.startAt) {
      options.startAt = "1990/01/01";
    }
    if (!options.endAt) {
      options.endAt = new Date()
        .toISOString()
        .slice(0, 10)
        .replaceAll("-", "/");
    }
    options.startAt = parseDate(options.startAt, { dateOnly: true });
    options.endAt = parseDate(options.endAt, { dateOnly: true });
  }
  let result = await getTransactions(id, page.value, perPage.value, options);
  loading.value = false;
  if (result.error) return;

  updatePagination(result);
}
async function exportTransactionsData() {
  if (!document.getElementById("export-form").reportValidity()) {
    return;
  }
  window.open(
    getExportTransactionsUrl({
      usersID: id,
      statusList: exportDialog.value.statusList,
      startAt: parseDate(exportDialog.value.startAt, { dateOnly: true }),
      endAt: parseDate(exportDialog.value.endAt, { dateOnly: true }),
    })
  );
  exportDialog.value = false;
}
</script>
<template>
  <div>
    <PageTitle v-if="showHeader">
      交易訂單
      <template #action> 訂單數：{{ totalRecord.toLocaleString() }} </template>
    </PageTitle>
    <div class="m-6 flex gap-2 justify-between items-end">
      <div class="flex gap-2 items-end">
        <FormSelect label="狀態" v-model="form.status">
          <option value="0">全部</option>
          <option value="1">完成</option>
          <option value="2">失敗</option>
        </FormSelect>
        <FormSelect label="方向" v-model="form.side">
          <option value="0">全部</option>
          <option value="1">買</option>
          <option value="2">賣</option>
        </FormSelect>
        <FormDate
          label="交易日期"
          v-model:start-date="form.startAt"
          v-model:end-date="form.endAt"
        />
        <FormInput
          label="搜尋"
          placeholder="訂單編號 / UID"
          v-model="form.search"
        />
        <Button
          @click="
            form = {
              status: '0',
              side: '0',
              startAt: '',
              endAt: '',
              search: '',
            }
          "
        >
          重設
        </Button>
      </div>
      <div class="flex gap-2 items-end">
        <Button
          variant="export"
          @click="
            exportDialog = {
              statusList: [1, 2],
              startAt: '',
              endAt: '',
            }
          "
          v-show="[1, 2, 4, 5].includes(userInfo.managersRolesID)"
        >
          匯出
        </Button>
      </div>
    </div>
    <TablePaginationContainer
      :cols="`7.5em 3em 5em 2em repeat(5, 1fr) ${
        id == 0 ? '6em' : ''
      } 8.65em 2em 2em`"
      v-model:perPage="perPage"
      v-model:page="page"
      :totalPages="totalPages"
      v-if="!loading"
    >
      <TableHeader>
        <div>訂單編號</div>
        <div>狀態</div>
        <div>交易對</div>
        <div>方向</div>
        <div class="text-right">數量</div>
        <div class="text-right">價格</div>
        <div class="text-right">成交額</div>
        <div class="text-right">總價值</div>
        <div class="text-right">手續費</div>
        <div v-if="id == 0">UID</div>
        <div>交易時間</div>
        <div>幣安</div>
        <div>發票</div>
      </TableHeader>
      <TableItem v-for="item in data" :key="item">
        <div class="tracking-tighter">{{ item.transactionsID }}</div>
        <div>
          <StatusBadge
            :status="['in-progress', 'completed', 'cancelled'][item.status]"
            variant="filled"
          >
            {{ ["", "完成", "失敗"][item.status] }}
          </StatusBadge>
        </div>
        <div class="tracking-tighter">
          {{ item.baseSymbol }}/{{ item.quoteSymbol }}
        </div>
        <div
          :class="`text-center ${
            item.side == 1 ? `text-[#3663A7]` : `text-[#DD4949]`
          }`"
        >
          {{ ["", "買", "賣"][item.side] }}
        </div>
        <ParsePrice :price="item.quantity" :unit="item.baseSymbol" />
        <ParsePrice :price="item.price" />
        <ParsePrice :price="item.amount" :unit="item.quoteSymbol" />
        <a-tooltip :mouseEnterDelay="0" :mouseLeaveDelay="0" color="white">
          <template #title>
            <span class="text-black">匯率：{{ item.twdExchangeRate }} </span>
          </template>
          <div>
            <div class="pointer-events-none">
              <ParsePrice :price="item.twdTotalValue" unit="TWD" />
            </div>
          </div>
        </a-tooltip>
        <ParsePrice
          :price="item.handlingCharge * item.twdExchangeRate"
          unit="TWD"
        />
        <div v-if="id == 0">{{ item.usersID }}</div>
        <div class="tracking-tighter">{{ item.createdAt }}</div>
        <div
          class="flex items-center justify-center cursor-pointer"
          @click="orderIdDialog = item"
        >
          <file-text-outlined />
        </div>
        <div
          class="flex items-center justify-center cursor-pointer"
          @click="invoiceDialog = item"
          v-if="item.invoiceStatus != 0"
        >
          <file-text-outlined />
        </div>
      </TableItem>
    </TablePaginationContainer>
    <Loader v-if="loading" />
    <Dialog v-model="orderIdDialog" title="幣安訂單詳細">
      <div
        class="grid grid-cols-[6rem_1fr] gap-4 items-center justify-end"
        v-if="orderIdDialog"
      >
        <span>訂單編號：</span>
        <div>{{ orderIdDialog.binanceID }}</div>
        <span>交易方向：</span>
        <div>{{ ["", "買", "賣"][orderIdDialog.side] }}</div>
        <span>成交數量：</span>
        <span class="flex items-start">
          <ParsePrice
            :price="orderIdDialog.binanceQuantity"
            :unit="orderIdDialog.baseSymbol"
            show-all-numbers
          />
        </span>
        <span>成交價格：</span>
        <span class="flex items-start">
          <ParsePrice :price="orderIdDialog.binancePrice" />
        </span>
        <span>成交額：</span>
        <span class="flex items-start">
          <ParsePrice
            :price="orderIdDialog.binanceAmount"
            :unit="orderIdDialog.quoteSymbol"
            show-all-numbers
          />
        </span>
        <span>手續費：</span>
        <span class="flex items-start">
          <ParsePrice :price="orderIdDialog.binanceHandlingCharge" />
        </span>
      </div>
      <template #actions>
        <Button @click="orderIdDialog = false"> 確定 </Button>
      </template>
    </Dialog>
    <Dialog v-model="invoiceDialog" title="發票資訊">
      <div
        class="grid grid-cols-[6rem_1fr] gap-4 items-center justify-end"
        v-if="invoiceDialog"
      >
        <span> 發票號碼： </span>
        <div v-if="invoiceDialog.invoiceID !== ''">{{ invoiceDialog.invoiceID }}</div>
        <div v-else>-</div>
        <span> 發票金額： </span>
        <span class="flex items-start">
          <ParsePrice :price="invoiceDialog.invoiceAmount" unit="TWD" />
        </span>
        <span> 發票狀態： </span>
        <div v-if="invoiceDialog.invoiceStatus === 1">尚未開立</div>
        <div v-else-if="invoiceDialog.invoiceStatus === 2">開立中</div>
        <div v-else-if="invoiceDialog.invoiceStatus === 3">已開立</div>
        <div v-else-if="invoiceDialog.invoiceStatus === 4">開立失敗</div>
      </div>
      <template #actions>
        <Button @click="invoiceDialog = false"> 確定 </Button>
      </template>
    </Dialog>
    <Dialog v-model="exportDialog" title="匯出報表" :loading="loading">
      <form
        @submit.prevent
        class="grid grid-cols-[4rem_1fr] gap-4 items-center justify-end"
        id="export-form"
        v-if="exportDialog"
      >
        <FormLabel label="日期" />
        <FormDate
          v-model:start-date="exportDialog.startAt"
          v-model:end-date="exportDialog.endAt"
          required
        />
        <FormLabel label="狀態" />
        <a-checkbox-group
          v-model:value="exportDialog.statusList"
          :options="[
            {
              label: '完成',
              value: 1,
            },
            {
              label: '失敗',
              value: 2,
            },
          ]"
          required
        />
      </form>
      <template #actions>
        <Button variant="outline" @click="exportDialog = false"> 取消 </Button>
        <Button @click="exportTransactionsData"> 確認 </Button>
      </template>
    </Dialog>
  </div>
</template>
