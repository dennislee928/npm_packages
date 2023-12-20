<script setup>
import {
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  CloseCircleOutlined,
} from "@ant-design/icons-vue";
import useAPI from "@/hooks/useAPI";
import usePagination from "@/hooks/usePagination";
import useReviewStatus from "@/hooks/useReviewStatus";
import useParseDate from "@/hooks/useParseDate";
import useUserInfo from "@/hooks/useUserInfo";
import useUserOptions from "@/hooks/useUserOptions";

const { getKycList, getExportKycListUrl } = useAPI();
const router = useRouter();
const parseDate = useParseDate();
const { userInfo } = useUserInfo();
const { result: countryList, loading: countryLoading } =
  useUserOptions("countries");

const loading = ref(true);
const exportDialog = ref(null);

const { data, page, perPage, totalPages, updatePagination } = usePagination();
const { status, badgeStatus } = useReviewStatus();
const form = ref({
  complianceReview: "0",
  finalReview: "0",
  startAt: "",
  endAt: "",
  search: "",
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
  for (let key of Object.keys(form.value)) {
    if (form.value[key] !== "") {
      options[key] = form.value[key];
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
  let result = await getKycList(page.value, perPage.value, options);
  loading.value = false;
  if (result.error) return;

  updatePagination(result);
}
async function exportKycData() {
  if (!document.getElementById("export-kyc-form").reportValidity()) {
    return;
  }
  window.open(
    getExportKycListUrl({
      statusList: exportDialog.value.statusList,
      startAt: parseDate(exportDialog.value.startAt, { dateOnly: true }),
      endAt: parseDate(exportDialog.value.endAt, { dateOnly: true }),
    })
  );
  exportDialog.value = false;
}
</script>
<template>
  <PageTitle> 審查列表 </PageTitle>
  <div class="m-6 flex gap-2 justify-between items-end">
    <div class="flex gap-2 items-end">
      <FormSelect label="審查結果" v-model="form.finalReview">
        <option :value="i" v-for="(item, i) of status" :key="i">
          {{ item }}
        </option>
      </FormSelect>
      <FormSelect label="法遵審查" v-model="form.complianceReview">
        <option :value="i" v-for="(item, i) of status" :key="i">
          {{ item }}
        </option>
      </FormSelect>
      <FormDate
        label="送審日期"
        v-model:start-date="form.startAt"
        v-model:end-date="form.endAt"
      />
      <FormInput label="搜尋" v-model="form.search" placeholder="UID / 姓名" />
      <Button
        @click="
          form = {
            complianceReview: '0',
            finalReview: '0',
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
            startAt: new Date(new Date() - 30 * 24 * 60 * 60 * 1000)
              .toISOString()
              .slice(0, 10),
            endAt: new Date().toISOString().slice(0, 10),
            statusList: [1, 2, 3, 4],
          }
        "
        v-show="[1, 2, 4].includes(userInfo.managersRolesID)"
      >
        匯出
      </Button>
    </div>
  </div>
  <TablePaginationContainer
    cols="1fr 1fr 1fr 1fr 1fr 1fr 1fr 10em 10em"
    v-model:perPage="perPage"
    v-model:page="page"
    :totalPages="totalPages"
    v-if="!loading"
  >
    <TableHeader>
      <div>UID</div>
      <div>審核結果</div>
      <div>姓</div>
      <div>名</div>
      <div>國籍</div>
      <div>風險評估</div>
      <div>法遵審查</div>
      <div>送審日期</div>
      <div>審查日期</div>
    </TableHeader>
    <TableItem
      v-for="item in data"
      :key="item.id"
      class="cursor-pointer"
      @click="router.push(`/kyc/${item.usersID}`)"
    >
      <div>{{ item.usersID }}</div>
      <StatusBadge :status="badgeStatus[item.finalReview]">
        {{ status[item.finalReview] }}
      </StatusBadge>
      <div>{{ item.lastName }}</div>
      <div>{{ item.firstName }}</div>
      <div v-if="countryLoading">{{ item.countriesCode }}</div>
      <div v-else>
        {{
          countryList.filter((x) => x.code === item.countriesCode)[0]?.chinese
        }}
      </div>
      <div
        v-if="item.internalRisksTotal <= 7"
        class="text-[#02C879] flex items-center gap-1"
      >
        <CheckCircleOutlined />{{ item.internalRisksTotal }}
      </div>
      <div
        v-else-if="item.internalRisksTotal <= 11"
        class="text-[#FF9D18] flex items-center gap-1"
      >
        <ExclamationCircleOutlined />{{ item.internalRisksTotal }}
      </div>
      <div
        v-else-if="item.internalRisksTotal > 12"
        class="text-[#FF574C] flex items-center gap-1"
      >
        <CloseCircleOutlined />{{ item.internalRisksTotal }}
      </div>
      <div v-else>未審查</div>
      <StatusBadge :status="badgeStatus[item.complianceReview]">
        {{ status[item.complianceReview] }}
      </StatusBadge>
      <div>{{ item.createdAt }}</div>
      <div>{{ item.finalReviewTime }}</div>
    </TableItem>
  </TablePaginationContainer>
  <Loader v-else />
  <Dialog v-model="exportDialog" title="匯出報表" :loading="loading">
    <form
      @submit.prevent
      class="grid grid-cols-[4rem_1fr] gap-4 items-center justify-end"
      id="export-kyc-form"
      v-if="exportDialog"
    >
      <FormLabel label="註冊日期" />
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
      <Button @click="exportKycData"> 送出 </Button>
    </template>
  </Dialog>
</template>
