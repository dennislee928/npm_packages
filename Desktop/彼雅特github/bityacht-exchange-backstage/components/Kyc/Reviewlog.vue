<script setup>
import useAPI from "@/hooks/useAPI";
import usePagination from "@/hooks/usePagination";
import useReviewStatus from "@/hooks/useReviewStatus";

const { status, badgeStatus } = useReviewStatus();
const { getUserKycReviewslogs, getExportUserKycReviewslogsUrl } = useAPI();
const { data, page, perPage, totalPages, updatePagination } = usePagination();
const route = useRoute();

const { uid } = route.params;
const loading = ref(false);
onMounted(async () => {
  getData();
});
watch([page, perPage], () => {
  getData();
});
async function getData() {
  loading.value = true;
  let result = await getUserKycReviewslogs(uid, page.value, perPage.value);
  loading.value = false;
  if (result.error) return;

  updatePagination(result);
}
async function exportData() {
  window.open(getExportUserKycReviewslogsUrl(uid));
}
</script>
<template>
  <div>
    <div class="flex item-center justify-end">
      <Button @click="exportData">匯出</Button>
    </div>
    <TablePaginationContainer
      v-model:perPage="perPage"
      v-model:page="page"
      :totalPages="totalPages"
      v-if="!loading"
    >
      <TableHeader>
        <div>日期</div>
        <div>類別</div>
        <div>狀態</div>
        <div>備註</div>
        <div>管理者</div>
      </TableHeader>
      <TableItem v-for="item of data">
        <div>{{ item.createdAt }}</div>
        <div>
          {{
            [
              "KryptoGO 複核",
              "姓名檢核排除評估",
              "內部風險審核",
              "法遵審查",
              "最終審查",
              "上傳認證結果(外籍人士)",
              "認證結果(外籍人士)",
            ][item.type - 1]
          }}
        </div>
        <div>
          <StatusBadge
            :status="badgeStatus[item.status]"
            v-if="item.status > 0"
          >
            {{ status[item.status] }}
          </StatusBadge>
          <span v-else class="opacity-50">-</span>
        </div>
        <div>
          <div v-if="item.comment">{{ item.comment }}</div>
          <div v-else class="opacity-50">-</div>
        </div>
        <div>
          <div v-if="item.managersName">{{ item.managersName }}</div>
          <div v-else class="opacity-50">-</div>
        </div>
      </TableItem>
    </TablePaginationContainer>
    <Loader v-if="loading" />
  </div>
</template>
