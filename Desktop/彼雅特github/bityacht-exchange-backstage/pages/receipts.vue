<script setup>
import { FileTextOutlined, ReloadOutlined } from "@ant-design/icons-vue";
import useAPI from "@/hooks/useAPI";
import useNotification from "@/hooks/useNotification";
import usePagination from "@/hooks/usePagination";
import useParseDate from "@/hooks/useParseDate";

const { getReceipts, getExportReceiptsUrl, issueReceipts } = useAPI();
const toast = useNotification();
const parseDate = useParseDate();

const form = ref({
  status: "",
  startAt: "",
  endAt: "",
  search: "",
});
const issueDialog = ref(false);
const detailId = ref(null);
const detailDialog = computed({
  get() {
    return !!detailId.value;
  },
  set(val) {
    if (!val) {
      detailId.value = null;
    }
  },
});
const exportDialog = ref(false);
const loading = ref(true);
const mutipleSelection = ref([]);

const { data, page, perPage, totalPages, updatePagination } = usePagination();
const issueableItems = computed(() => {
  return data.value.filter((item) => item.status == 1) || [];
});

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
  let result = await getReceipts(page.value, perPage.value, options);
  loading.value = false;
  if (result.error) return;

  updatePagination(result);
}
function onSelectAll() {
  if (mutipleSelection.value.length !== issueableItems.value.length) {
    mutipleSelection.value = issueableItems.value.map((item) => item.id);
  } else {
    mutipleSelection.value = [];
  }
}
async function exportReceiptsData() {
  if (!document.getElementById("export-receipt-form").reportValidity()) {
    return;
  }
  window.open(
    getExportReceiptsUrl({
      statuses: exportDialog.value.statusList,
      startAt: parseDate(exportDialog.value.startAt, { dateOnly: true }),
      endAt: parseDate(exportDialog.value.endAt, { dateOnly: true }),
    })
  );
  exportDialog.value = false;
}
async function issueReceipt() {
  if (mutipleSelection.value.length === 0) {
    toast.error("請先選取發票");
    return;
  }
  loading.value = true;
  let result = await issueReceipts(mutipleSelection.value);
  if (!result.error) {
    toast.success("開立成功");
    issueDialog.value = false;
    mutipleSelection.value = [];
    getData();
  } else {
    toast.error("開立失敗");
  }
  loading.value = false;
}
</script>
<template>
  <div>
    <PageTitle> 發票管理 </PageTitle>
    <div class="m-6 flex gap-2 justify-between items-end">
      <div class="flex gap-2 items-end">
        <FormSelect label="狀態" v-model="form.status">
          <option value="">全部</option>
          <option value="1">尚未開立</option>
          <option value="2">開立中</option>
          <option value="3">已開立</option>
          <option value="4">開立失敗</option>
        </FormSelect>
        <FormDate
          label="訂單日期"
          v-model:start-date="form.startAt"
          v-model:end-date="form.endAt"
        />
        <FormInput
          label="關鍵字"
          placeholder="訂單編號/發票號碼"
          v-model="form.search"
        />
        <Button
          @click="
            form = {
              status: '',
              startAt: '',
              endAt: '',
              search: '',
            }
          "
        >
          重設
        </Button>
        <Button
          variant="export"
          @click="
            exportDialog = {
              startAt: form.startAt,
              endAt: form.endAt,
              statusList: [1, 2, 3, 4],
            }
          "
        >
          匯出
        </Button>
      </div>
      <div class="flex gap-2 items-end">
        <Button @click="getData"> <reload-outlined /> 重新整理 </Button>
        <Button @click="issueDialog = true"> 批次開立 </Button>
      </div>
    </div>
    <TablePaginationContainer
      cols="2em 0.5fr 1fr 5em 0.7fr 0.5fr 1fr 2em"
      v-model:perPage="perPage"
      v-model:page="page"
      :totalPages="totalPages"
      v-if="!loading"
    >
      <TableHeader>
        <div>
          <a-checkbox
            :checked="
              mutipleSelection.length === issueableItems.length &&
              data.length !== 0
            "
            @change="onSelectAll"
          />
        </div>
        <div>訂單編號</div>
        <div>訂單日期</div>
        <div>發票狀態</div>
        <div>發票號碼</div>
        <div>發票金額</div>
        <div>開立時間</div>
        <div>明細</div>
      </TableHeader>
      <TableItem v-for="item in data" :key="item">
        <div>
          <a-checkbox
            v-if="item.status === 1"
            :checked="mutipleSelection.includes(item.id)"
            @change="
              mutipleSelection.includes(item.id)
                ? (mutipleSelection = mutipleSelection.filter(
                    (id) => id !== item.id
                  ))
                : mutipleSelection.push(item.id)
            "
          >
          </a-checkbox>
        </div>
        <div>{{ item.id }}</div>
        <div>{{ item.createdAt }}</div>
        <div>
          <StatusBadge
            :status="
              ['', 'in-progress', 'frozen', 'completed', 'cancelled'][
                item.status
              ]
            "
            variant="filled"
          >
            {{ ["", "未開立", "開立中", "已開立", "失敗"][item.status] }}
          </StatusBadge>
        </div>
        <div>{{ item.invoiceID }}</div>
        <ParsePrice :price="item.invoiceAmount" unit="TWD" />
        <div>{{ item.invoiceIssuedAt }}</div>
        <div class="cursor-pointer" @click="detailId = item.id">
          <file-text-outlined />
        </div>
      </TableItem>
    </TablePaginationContainer>
    <Loader v-if="loading" />
    <Dialog v-model="exportDialog" title="匯出報表" :loading="loading">
      <form
        @submit.prevent
        class="grid grid-cols-[4rem_1fr] gap-4 items-center justify-end"
        id="export-receipt-form"
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
              label: '尚未開立',
              value: 1,
            },
            {
              label: '開立中',
              value: 2,
            },
            {
              label: '已開立',
              value: 3,
            },
            {
              label: '開立失敗',
              value: 4,
            },
          ]"
          required
        />
      </form>
      <template #actions>
        <Button variant="outline" @click="exportDialog = false"> 取消 </Button>
        <Button @click="exportReceiptsData"> 確定 </Button>
      </template>
    </Dialog>
    <Dialog v-model="issueDialog" title="開立發票" :loading="loading">
      <p>
        請確認是否進行批次開立發票動作？
        <br />
        批次開立數量：{{ mutipleSelection.length }} 筆
      </p>
      <template #actions>
        <Button variant="outline" @click="issueDialog = false"> 取消 </Button>
        <Button @click="issueReceipt"> 確定 </Button>
      </template>
    </Dialog>
    <Dialog v-model="detailDialog" title="發票明細" :loading="loading">
      <ReceiptsDetail :id="detailId" v-if="detailDialog" />
      <template #actions>
        <Button @click="detailDialog = false"> 確定 </Button>
      </template>
    </Dialog>
  </div>
</template>
