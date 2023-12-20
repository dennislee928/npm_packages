<script setup>
import useUserInfo from "@/hooks/useUserInfo";
import useSuspiciousReviewResult from "@/hooks/useSuspiciousReviewResult";
import useSuspiciousType from "@/hooks/useSuspiciousType";
import useParseDate from "@/hooks/useParseDate";
import usePagination from "@/hooks/usePagination";
import useAPI from "@/hooks/useAPI";
const { types } = useSuspiciousType();
const { status, badgeStatus } = useSuspiciousReviewResult();
const { userInfo } = useUserInfo();
const { getSuspiciousTransactions, getExportSuspiciousTransactionsUrl } =
  useAPI();
const router = useRouter();
const parseDate = useParseDate();

const loading = ref(false);
const form = ref({
  dedicatedReviewResult: "0",
  type: "0",
  date: {
    startDate: "",
    endDate: "",
  },
});
const exportDialog = ref(false);

const { data, page, perPage, totalPages, totalRecord, updatePagination } =
  usePagination();

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
  let result = await getSuspiciousTransactions(
    page.value,
    perPage.value,
    options
  );
  loading.value = false;
  if (result.error) return;

  updatePagination(result);
}
function exportSuspect() {
  if (!document.getElementById("export-suspect-form").reportValidity()) {
    return;
  }
  window.open(
    getExportSuspiciousTransactionsUrl({
      dedicatedReviewResult: exportDialog.value.dedicatedReviewResult,
      type: exportDialog.value.type,
      startAt: parseDate(exportDialog.value.startAt, { dateOnly: true }),
      endAt: parseDate(exportDialog.value.endAt, { dateOnly: true }),
    })
  );
  exportDialog.value = false;
}
</script>
<template>
  <div>
    <PageTitle> 可疑態樣管理 </PageTitle>
    <div class="m-6 flex gap-2 justify-between items-end">
      <div class="flex gap-2 items-end">
        <FormSelect label="狀態" v-model="form.dedicatedReviewResult">
          <option v-for="(item, i) of status" :value="i.toString()">
            {{ item }}
          </option>
        </FormSelect>
        <FormSelect label="可疑樣態" v-model="form.type">
          <option v-for="(item, i) of types" :value="i.toString()">
            {{ item }}
          </option>
        </FormSelect>
        <FormDate
          label="掃描日期"
          v-model:start-date="form.date.startDate"
          v-model:end-date="form.date.endDate"
        />
        <FormInput label="搜尋" placeholder="名/UID/案件單號" />
        <Button> 重設 </Button>
      </div>
      <div class="flex gap-2 items-end">
        <Button
          variant="export"
          @click="
            exportDialog = {
              dedicatedReviewResult: [1, 2, 3],
              type: [1, 2, 3, 4, 5, 6],
              date: {
                startDate: '',
                endDate: '',
              },
            }
          "
          v-show="[1, 2, 4, 5].includes(userInfo.managersRolesID)"
        >
          匯出
        </Button>
      </div>
    </div>
    <TablePaginationContainer
      cols="0.7fr 0.5fr 1fr 0.25fr 0.5fr 1fr 1fr 1fr"
      v-model:perPage="perPage"
      v-model:page="page"
      :totalPages="totalPages"
      v-if="!loading"
    >
      <TableHeader>
        <div>案件單號</div>
        <div>狀態</div>
        <div>訂單編號</div>
        <div>姓</div>
        <div>名</div>
        <div>E-MAIL</div>
        <div>可疑態樣</div>
        <div>掃描時間</div>
      </TableHeader>
      <TableItem
        v-for="item in data"
        :key="item"
        class="cursor-pointer"
        @click="router.push(`/suspect/${item.id}/transaction-info`)"
      >
        <div>{{ item.id }}</div>
        <div>
          <StatusBadge
            :status="badgeStatus[item.dedicatedReviewResult]"
            variant="filled"
          >
            {{ status[item.dedicatedReviewResult] }}
          </StatusBadge>
        </div>
        <div>{{ item.orderID }}</div>
        <div>{{ item.lastName }}</div>
        <div>{{ item.firstName }}</div>
        <div>{{ item.email }}</div>
        <div>{{ types[item.type] }}</div>
        <div>{{ item.createdAt }}</div>
      </TableItem>
    </TablePaginationContainer>
    <Loader v-else />
    <Dialog v-model="exportDialog" title="匯出報表">
      <form
        @submit.prevent
        class="grid grid-cols-[4rem_1fr] gap-4 items-center justify-end"
        id="export-suspect-form"
        v-if="exportDialog"
      >
        <FormLabel label="掃描時間" />
        <FormDate
          v-model:start-date="exportDialog.startAt"
          v-model:end-date="exportDialog.endAt"
          required
        />

        <FormLabel label="狀態" />
        <a-checkbox-group
          v-model:value="exportDialog.dedicatedReviewResult"
          :options="
            status.slice(1).map((item, i) => ({
              label: item,
              value: i + 1,
            }))
          "
          required
        />
        <FormLabel label="可疑態樣" />
        <a-checkbox-group
          v-model:value="exportDialog.type"
          :options="
            types.slice(1).map((item, i) => ({
              label: item,
              value: i + 1,
            }))
          "
          required
        />
      </form>
      <template #actions>
        <Button variant="outline" @click="exportDialog = false"> 取消 </Button>
        <Button @click="exportSuspect"> 確認 </Button>
      </template>
    </Dialog>
  </div>
</template>
