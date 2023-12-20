<script setup>
import useAPI from "@/hooks/useAPI";
import usePagination from "@/hooks/usePagination";
import useUserStatus from "@/hooks/useUserStatus";
import useUserInfo from "@/hooks/useUserInfo";

const { status, badgeStatus } = useUserStatus();
const { getUserBankLogs, getExportUserBankLogsUrl } = useAPI();
const { data, page, perPage, totalPages, updatePagination } = usePagination();
const route = useRoute();
const { userInfo } = useUserInfo();

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
  let result = await getUserBankLogs(uid, page.value, perPage.value);
  loading.value = false;
  if (result.error) return;

  updatePagination(result);
}
function exportData() {
  window.open(getExportUserBankLogsUrl(uid));
}
</script>
<template>
  <div>
    <div
      class="flex item-center justify-end"
      v-show="[1, 2, 4].includes(userInfo.managersRolesID)"
    >
      <Button @click="exportData">匯出</Button>
    </div>
    <TablePaginationContainer
      v-model:perPage="perPage"
      v-model:page="page"
      cols="1fr 1fr 3fr 1fr"
      :totalPages="totalPages"
      v-if="!loading"
    >
      <TableHeader>
        <div>日期</div>
        <div>狀態</div>
        <div>備註</div>
        <div>管理者</div>
      </TableHeader>
      <TableItem v-for="item of data">
        <div>{{ item.createdAt }}</div>
        <div>
          <StatusBadge
            :status="
              ['in-progress', 'in-progress', 'completed', 'cancelled'][
                item.status
              ]
            "
          >
            {{ ["未綁定", "審查中", "已通過", "未通過"][item.status] }}
          </StatusBadge>
        </div>
        <div>
          <div v-if="item.comment">{{ item.comment }}</div>
          <div v-else>-</div>
        </div>
        <div>
          <div v-if="item.managersName">{{ item.managersName }}</div>
          <div v-else>-</div>
        </div>
      </TableItem>
    </TablePaginationContainer>
    <Loader v-if="loading" />
  </div>
</template>
