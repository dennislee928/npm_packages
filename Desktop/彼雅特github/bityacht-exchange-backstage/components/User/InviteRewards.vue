<script setup>
import useAPI from "@/hooks/useAPI";
import usePagination from "@/hooks/usePagination";

const { getUserInviteRewards, getExportUserInviteRewardsUrl } = useAPI();
const route = useRoute();

const { uid } = route.params;
const exportDialog = ref(false);
const loading = ref(false);
const { data, page, perPage, totalPages, updatePagination } = usePagination();
const form = ref({
  action: "0",
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
  let result = await getUserInviteRewards(
    uid,
    page.value,
    perPage.value,
    form.value
  );
  loading.value = false;
  if (result.error) return;

  updatePagination(result);
}
function exportRewards() {
  window.open(getExportUserInviteRewardsUrl(uid, exportDialog.value.action));
  exportDialog.value = false;
}
</script>
<template>
  <div class="m-6 flex gap-2 justify-between items-end">
    <div class="flex gap-2 items-end">
      <FormSelect label="類型" v-model="form.action">
        <option value="0">全部</option>
        <option value="1">返佣</option>
        <option value="2">提領</option>
      </FormSelect>
    </div>
    <div class="flex gap-2 items-end">
      <Button variant="export" @click="exportDialog = { action: 0 }">
        匯出
      </Button>
    </div>
  </div>

  <TablePaginationContainer
    cols="1fr 1fr 0.5fr 2fr"
    v-model:perPage="perPage"
    v-model:page="page"
    :totalPages="totalPages"
    v-if="!loading"
  >
    <TableHeader>
      <div>日期</div>
      <div>類型</div>
      <div class="text-right">數量</div>
      <div>帳號</div>
    </TableHeader>
    <TableItem v-for="item of data">
      <div>{{ item.createdAt }}</div>
      <div>{{ item.action == 1 ? "返佣" : "提領" }}</div>
      <ParsePrice :price="item.amount" unit="USDT" />
      <div>{{ item.account }}</div>
    </TableItem>
  </TablePaginationContainer>
  <Loader v-if="loading" />
  <Dialog v-model="exportDialog" title="匯出邀請獎勵">
    <form
      @submit.prevent
      class="grid grid-cols-[4rem_1fr] gap-4 items-center justify-end"
      id="export-users-form"
      v-if="exportDialog"
    >
      <FormLabel label="類別" />
      <FormSelect v-model="exportDialog.action">
        <option value="0">全部</option>
        <option value="1">返佣</option>
        <option value="2">提領</option>
      </FormSelect>
    </form>
    <template #actions>
      <Button variant="outline" @click="exportDialog = false"> 取消 </Button>
      <Button @click="exportRewards"> 送出 </Button>
    </template>
  </Dialog>
</template>
