<script setup>
import useAPI from "@/hooks/useAPI";
import useNotification from "@/hooks/useNotification";
import usePagination from "@/hooks/usePagination";

const {
  getWithdrawalWhitelist,
  getExportWithdrawalWhitelistUrl,
  deleteWithdrawalWhitelist,
} = useAPI();
const route = useRoute();
const toast = useNotification();

const { uid } = route.params;
const loading = ref(false);
const deleteDialog = ref(null);
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
  let result = await getWithdrawalWhitelist(uid, page.value, perPage.value);
  loading.value = false;
  if (result.error) return;

  updatePagination(result);
}
async function exportWhitelist() {
  window.open(getExportWithdrawalWhitelistUrl(uid));
}
async function deleteWhitelistItem() {
  let id = deleteDialog.value.id;
  loading.value = true;
  let result = await deleteWithdrawalWhitelist(uid, id);
  loading.value = false;
  if (result.error) {
    toast.error("刪除失敗");
    return;
  }
  deleteDialog.value = null;
  getData();
}
</script>
<template>
  <TablePaginationContainer
    cols="1fr 1fr 1fr 3em"
    v-model:perPage="perPage"
    v-model:page="page"
    :totalPages="totalPages"
    v-if="!loading"
  >
    <TableHeader>
      <div>白名單地址</div>
      <div></div>
      <div></div>
      <div class="flex items-center justify-end">
        <Button size="mini" @click="exportWhitelist"> 匯出 </Button>
      </div>
    </TableHeader>
    <TableHeader class="bg-[#E7E7E7] border-[#E7E7E7] text-black">
      <div>區塊鏈</div>
      <div>白名單地址</div>
      <div>備註</div>
    </TableHeader>
    <TableItem v-for="item of data">
      <div>{{ item.mainnet }}</div>
      <div>{{ item.address }}</div>
      <div>{{ item.extra.memo }}</div>
      <div class="flex items-center justify-end">
        <Button size="mini" variant="danger" @click="deleteDialog = item">
          刪除
        </Button>
      </div>
    </TableItem>
  </TablePaginationContainer>
  <Loader v-if="loading" />

  <Dialog v-model="deleteDialog" :loading="loading" title="確認刪除">
    <div>確定要刪除此白名單地址嗎？</div>
    <template #actions>
      <Button variant="outline" @click="deleteDialog = false"> 取消 </Button>
      <Button variant="danger" @click="deleteWhitelistItem"> 確定 </Button>
    </template>
  </Dialog>
</template>
