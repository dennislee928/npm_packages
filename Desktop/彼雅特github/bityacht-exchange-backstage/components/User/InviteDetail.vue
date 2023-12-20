<script setup>
import useAPI from "@/hooks/useAPI";
import usePagination from "@/hooks/usePagination";

const { getUserInviteRecords } = useAPI();
const route = useRoute();

const { uid } = route.params;
const loading = ref(false);
const { data, page, perPage, totalPages, updatePagination } = usePagination();
onMounted(async () => {
  getData();
});
watch(
  [page, perPage],
  () => {
    getData();
  },
  { deep: true }
);

async function getData() {
  loading.value = true;
  let result = await getUserInviteRecords(uid, page.value, perPage.value);
  loading.value = false;
  if (result.error) return;

  updatePagination(result);
}
</script>
<template>
  <div class="-m-4">
    <TablePaginationContainer
      cols="1fr 1fr 1fr"
      v-model:perPage="perPage"
      v-model:page="page"
      :totalPages="totalPages"
      v-if="!loading"
    >
      <TableHeader>
        <div>帳號</div>
        <div>狀態</div>
        <div>註冊日期</div>
      </TableHeader>
      <TableItem v-for="item of data">
        <div>{{ item.account }}</div>
        <div>{{ item.status == 1 ? "未完成" : "已完成" }}</div>
        <div>{{ item.createdAt }}</div>
      </TableItem>
    </TablePaginationContainer>
    <Loader v-if="loading" />
  </div>
</template>
