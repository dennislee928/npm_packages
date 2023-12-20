<script setup>
import useAPI from "@/hooks/useAPI";
import usePagination from "@/hooks/usePagination";
import useOrderStatus from "@/hooks/useOrderStatus";
import useUserInfo from "@/hooks/useUserInfo";
import useParseDate from "@/hooks/useParseDate";
import useNotification from "@/hooks/useNotification";
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
const {
  getSpots,
  getExportSpotsUrl,
  getTransactionpairs,
  getExportSpotsAegisUrl,
  importSpotsAegis,
} = useAPI();
const { status, badgeStatus } = useOrderStatus();
const { userInfo } = useUserInfo();
const parseDate = useParseDate();
const toast = useNotification();

const form = ref({
  status: "0",
  coin: "",
  startAt: "",
  endAt: "",
});
const exportDialog = ref(false);
const exportAegisDialog = ref(false);
const importDialog = ref(false);
const loading = ref(false);
const coinOptions = ref([]);

const isSingleUser = computed(() => id != 0);

const { data, page, perPage, totalPages, totalRecord, updatePagination } =
  usePagination();

onMounted(async () => {
  getData();
  getCoinOptions();
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
    options.startAt = parseDate(
      options.startAt,
      { dateOnly: true },
      { dateOnly: false }
    );
    options.endAt = parseDate(
      options.endAt,
      { dateOnly: true },
      { dateOnly: false }
    );
  }
  let result = await getSpots(id, page.value, perPage.value, options);
  loading.value = false;
  if (result.error) return;

  updatePagination(result);
}
async function getCoinOptions() {
  let result = await getTransactionpairs(1, 100);
  let coinOptionsList = new Set();
  result.data.map((item) => {
    coinOptionsList.add(item.baseCurrenciesSymbol);
    coinOptionsList.add(item.quoteCurrenciesSymbol);
  });
  coinOptions.value = [...coinOptionsList];
}

async function exportSpotsData() {
  if (!document.getElementById("export-spots-form").reportValidity()) {
    return;
  }
  // window.open(
  //   getExportSpotsUrl({
  //     id,
  //     statusList: exportDialog.value.statusList,
  //     startAt: parseDate(exportDialog.value.startAt, { dateOnly: true }),
  //     endAt: parseDate(exportDialog.value.endAt, { dateOnly: true }),
  //   })
  // );
  let url = getExportSpotsUrl({
    id,
    statusList: exportDialog.value.statusList,
    startAt: parseDate(exportDialog.value.startAt, { dateOnly: true }),
    endAt: parseDate(exportDialog.value.endAt, { dateOnly: true }),
  });
  let response = await fetch(url);
  let status = response.status;
  let filename = response.headers
    .get("Content-Disposition")
    .split("=")[1]
    .replaceAll('"', "");
  let blob = await response.blob();
  let link = document.createElement("a");
  link.href = window.URL.createObjectURL(blob);
  link.download = filename;
  link.click();

  exportDialog.value = false;
}
async function exportAegis() {
  if (!document.getElementById("export-aegis-form").reportValidity()) {
    return;
  }
  let url = getExportSpotsAegisUrl({
    usersID: id,
    mainnet: exportAegisDialog.value.mainnet,
    startAt: parseDate(exportAegisDialog.value.startAt, { dateOnly: true }),
    endAt: parseDate(exportAegisDialog.value.endAt, { dateOnly: true }),
  });
  // use fetch & blob to download file
  // use header to set filename\
  // status: 200, 204
  let response = await fetch(url);
  let status = response.status;
  if (status === 204) {
    toast.info(`無資料可供匯出`);
    return;
  }
  let filename = response.headers
    .get("Content-Disposition")
    .split("=")[1]
    .replaceAll('"', "");
  let blob = await response.blob();
  let link = document.createElement("a");
  link.href = window.URL.createObjectURL(blob);
  link.download = filename;
  link.click();

  exportAegisDialog.value = false;
}
async function importAegis() {
  if (!document.getElementById("import-form").reportValidity()) {
    return;
  }
  let formData = new FormData();
  formData.append("file", importDialog.value.file);
  loading.value = true;
  let result = await importSpotsAegis(formData);
  loading.value = false;
  if (result.error) {
    toast.error(`匯入失敗，請確認檔案格式是否正確`);
  } else {
    toast.success(`匯入成功`);
    importDialog.value = false;
    await getData();
  }
}
</script>
<template>
  <div>
    <PageTitle v-if="showHeader">
      提/入幣列表
      <template #action> 訂單數：{{ totalRecord.toLocaleString() }} </template>
    </PageTitle>
    <div class="m-6 flex gap-2 justify-between items-end">
      <div class="flex gap-2 items-end">
        <FormSelect label="狀態" v-model="form.status">
          <option v-for="(item, i) of status" :key="i" :value="i">
            {{ item }}
          </option>
        </FormSelect>
        <FormSelect label="幣種" v-model="form.coin">
          <option value="">全部</option>
          <option v-for="item of coinOptions" :value="item">{{ item }}</option>
        </FormSelect>
        <FormDate
          label="日期"
          v-model:start-date="form.startAt"
          v-model:end-date="form.endAt"
        />
        <FormInput
          label="搜尋"
          v-model="form.search"
          placeholder="訂單編號 / UID"
        />
        <Button
          @click="
            form = {
              status: '0',
              coin: '',
              startAt: '',
              endAt: '',
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
              startAt: form.startAt,
              endAt: form.endAt,
              statusList: [1, 2, 3],
            }
          "
          v-show="[1, 2, 4, 5].includes(userInfo.managersRolesID)"
        >
          匯出
        </Button>
        <Button
          @click="
            exportAegisDialog = {
              startAt: form.startAt,
              endAt: form.endAt,
              mainnet: 1,
            }
          "
          v-show="[1, 2, 4, 5].includes(userInfo.managersRolesID)"
        >
          Aegis 匯出
        </Button>
        <Button
          @click="
            importDialog = {
              file: null,
            }
          "
          v-show="[1, 2, 4, 5].includes(userInfo.managersRolesID)"
        >
          匯入
        </Button>
      </div>
    </div>

    <TablePaginationContainer
      :cols="`8em 4em 3em 4em 4em 6em 1.2fr ${
        !isSingleUser ? '6em' : ''
      } 1fr 1fr`"
      v-model:perPage="perPage"
      v-model:page="page"
      :totalPages="totalPages"
      v-if="!loading"
    >
      <TableHeader>
        <div>交易編號</div>
        <div>狀態</div>
        <div>類型</div>
        <div>幣種</div>
        <div>主網</div>
        <div class="text-right">數量</div>
        <div class="text-center">TXID</div>
        <div v-if="!isSingleUser">UID</div>
        <div>交易時間</div>
        <div>完成時間</div>
      </TableHeader>
      <TableItem v-for="item in data" :key="item">
        <!-- <NuxtLink class="tracking-tighter link" :to="`/wallet/${item.id}`">
          {{ item.id }}
        </NuxtLink> -->
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
        <div v-if="!isSingleUser">{{ item.usersID }}</div>
        <div>{{ item.createdAt }}</div>
        <div>{{ item.finishedAt }}</div>
      </TableItem>
    </TablePaginationContainer>
    <Loader v-if="loading" />
    <Dialog v-model="exportDialog" title="匯出報表" :loading="loading">
      <form
        @submit.prevent
        class="grid grid-cols-[4rem_1fr] gap-4 items-center justify-end"
        id="export-spots-form"
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
        <Button variant="outline" @click="exportDialog = false"> 取消 </Button>
        <Button @click="exportSpotsData"> 確認 </Button>
      </template>
    </Dialog>
    <Dialog
      v-model="exportAegisDialog"
      title="匯出 Aegis 報表"
      :loading="loading"
    >
      <form
        @submit.prevent
        class="grid grid-cols-[4rem_1fr] gap-4 items-center justify-end"
        id="export-aegis-form"
        v-if="exportAegisDialog"
      >
        <FormLabel label="日期" />
        <FormDate
          v-model:start-date="exportAegisDialog.startAt"
          v-model:end-date="exportAegisDialog.endAt"
          required
        />
        <FormLabel label="主網" />
        <FormSelect v-model="exportAegisDialog.mainnet">
          <option value="1">BTC</option>
          <option value="2">ETH</option>
          <option value="3">ERC20</option>
          <option value="4">TRC20</option>
        </FormSelect>
      </form>
      <template #actions>
        <Button variant="outline" @click="exportAegisDialog = false">
          取消
        </Button>
        <Button @click="exportAegis"> 確認 </Button>
      </template>
    </Dialog>
    <Dialog v-model="importDialog" title="匯入" :loading="loading">
      <form
        @submit.prevent
        class="grid grid-cols-[4rem_1fr] gap-x-4 gap-y-2 items-center justify-end"
        id="import-form"
        v-if="importDialog"
      >
        <FormLabel label="檔案" />
        <FormInput
          type="file"
          id="file"
          @change="
            (e) => {
              importDialog.file = e.target.files[0];
            }
          "
          accept=".csv"
          required
        />
        <div></div>
        <div class="text-red-500 opacity-50 text-sm">僅可上傳 .csv 格式</div>
      </form>
      <template #actions>
        <Button variant="outline" @click="importDialog = false"> 取消 </Button>
        <Button @click="importAegis"> 上傳 </Button>
      </template>
    </Dialog>
  </div>
</template>
