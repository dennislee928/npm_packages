<script setup>
import { EditOutlined } from "@ant-design/icons-vue";
import useAPI from "@/hooks/useAPI";
import useNotification from "@/hooks/useNotification";
import usePagination from "@/hooks/usePagination";

const { updateTransactionpairs, getMainnets, updateMainnets } = useAPI();
const toast = useNotification();

const editDialog = ref(false);
const loading = ref(true);

const { data, page, perPage, totalPages, updatePagination } = usePagination();

onMounted(async () => {
  getData();
});
watch([page, perPage], () => {
  getData();
});
async function getData() {
  loading.value = true;
  let result = await getMainnets(page.value, perPage.value);
  loading.value = false;
  if (result.error) return;

  updatePagination(result);
}
async function editMainnets() {
  loading.value = true;
  let {
    mainnet,
    currenciesSymbol: currency,
    withdrawFee,
    withdrawMin,
  } = editDialog.value;
  let result = await updateMainnets({
    mainnet,
    currency,
    withdrawFee,
    withdrawMin,
  });
  if (!result.error) {
    toast.success("修改成功");
    editDialog.value = false;
    await getData();
  } else {
    toast.error("修改失敗");
  }
  loading.value = false;
}

const deepCopy = (x) => JSON.parse(JSON.stringify(x));
</script>
<template>
  <div>
    <PageTitle>提幣/手續費設定</PageTitle>
    <TablePaginationContainer
      cols="repeat(4, 1fr) 2em"
      v-model:perPage="perPage"
      v-model:page="page"
      :totalPages="totalPages"
      v-if="!loading"
    >
      <TableHeader>
        <div>幣種</div>
        <div>鏈</div>
        <div>最小提幣數量</div>
        <div>手續費</div>
      </TableHeader>
      <TableItem v-for="item in data" :key="item.id">
        <div>{{ item.currenciesSymbol }}</div>
        <div>{{ item.name }}</div>
        <div>
          {{ item.withdrawMin }}
          <span class="text-xs opacity-60 ml-[2px]">
            {{ item.currenciesSymbol }}
          </span>
        </div>
        <div>
          {{ item.withdrawFee }}
          <span class="text-xs opacity-60 ml-[2px]">
            {{ item.currenciesSymbol }}
          </span>
        </div>
        <div class="flex justify-end">
          <EditOutlined @click="editDialog = deepCopy(item)" />
        </div>
      </TableItem>
    </TablePaginationContainer>
    <Loader v-else />
    <Dialog v-model="editDialog" :loading="loading" title="提幣/手續費設定">
      <div class="flex flex-col gap-4" v-if="editDialog">
        <div
          class="grid grid-cols-[7rem_1fr_4em] gap-4 items-center justify-end"
        >
          <FormLabel label="幣種" />
          <div>{{ editDialog.currenciesSymbol }}</div>
          <div></div>
          <FormLabel label="鏈" />
          <div>{{ editDialog.name }}</div>
          <div></div>
          <FormLabel label="最小提幣數量" />
          <FormInput type="number" v-model="editDialog.withdrawMin" />
          <div>{{ editDialog.currenciesSymbol }}</div>
          <FormLabel label="手續費" />
          <FormInput type="number" v-model="editDialog.withdrawFee" />
          <div>{{ editDialog.currenciesSymbol }}</div>
        </div>
      </div>
      <template #actions>
        <Button variant="outline" @click="editDialog = false"> 取消 </Button>
        <Button @click="editMainnets"> 確定 </Button>
      </template>
    </Dialog>
  </div>
</template>
