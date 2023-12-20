<script setup>
import useAPI from "@/hooks/useAPI";
import usePagination from "@/hooks/usePagination";

const { getUserLoginLog, getExportUserLoginLogUrl } = useAPI();
const route = useRoute();

const { uid } = route.params;
const loading = ref(false);
const { data, page, perPage, totalPages, updatePagination } = usePagination();

onMounted(async () => {
  getData();
});
watch([page, perPage], () => {
  getData();
});
async function getData() {
  loading.value = true;
  let result = await getUserLoginLog(uid, page.value, perPage.value);
  loading.value = false;
  if (result.error) return;

  updatePagination(result);
}

async function exportLog() {
  window.open(getExportUserLoginLogUrl(uid));
}
</script>
<template>
  <div class="m-6 flex gap-2 justify-between items-end">
    <div class="flex gap-2 items-end"></div>
    <div class="flex gap-2 items-end">
      <Button variant="export" @click="exportLog"> 匯出 </Button>
    </div>
  </div>
  <TablePaginationContainer
    cols="1fr 1fr 1fr 10em 10em"
    v-model:perPage="perPage"
    v-model:page="page"
    :totalPages="totalPages"
    v-if="!loading"
  >
    <TableHeader>
      <div>位置</div>
      <div>瀏覽器</div>
      <div>裝置</div>
      <div>IP 位址</div>
      <div>登入時間</div>
    </TableHeader>
    <TableItem v-for="item in data">
      <div>{{ item.location }}</div>
      <div>{{ item.browser }}</div>
      <div>{{ item.device }}</div>
      <div>{{ item.ip }}</div>
      <div>{{ item.createdAt }}</div>
    </TableItem>
  </TablePaginationContainer>
  <Loader v-if="loading" />
</template>
