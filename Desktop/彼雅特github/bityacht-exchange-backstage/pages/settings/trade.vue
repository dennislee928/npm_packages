<script setup>
import { EditOutlined } from "@ant-design/icons-vue";
import useAPI from "@/hooks/useAPI";
import useNotification from "@/hooks/useNotification";
import usePagination from "@/hooks/usePagination";

const { getTransactionpairs, updateTransactionpairs } = useAPI();
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
  let result = await getTransactionpairs(page.value, perPage.value);
  loading.value = false;
  if (result.error) return;

  updatePagination(result);
}
async function editTransactionpairs() {
  loading.value = true;
  let result = await updateTransactionpairs(editDialog.value);
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
    <PageTitle> 交易對/手續費設定 </PageTitle>
    <TablePaginationContainer
      cols="repeat(5, 1fr) 2em"
      v-model:perPage="perPage"
      v-model:page="page"
      :totalPages="totalPages"
      v-if="!loading"
    >
      <TableHeader>
        <div>交易對</div>
        <div>狀態</div>
        <div class="text-right">買幣匯差</div>
        <div class="text-right">賣幣匯差</div>
        <div class="text-right">交易手續費</div>
        <div></div>
      </TableHeader>
      <TableItem v-for="item in data" :key="item.id">
        <div>
          {{ item.baseCurrenciesSymbol }}/{{ item.quoteCurrenciesSymbol }}
        </div>
        <div>
          <StatusBadge status="completed" v-if="item.status">開啟</StatusBadge>
          <StatusBadge status="cancelled" v-else>關閉</StatusBadge>
        </div>
        <ParsePrice :price="item.spreadsOfBuy" unit="%" />
        <ParsePrice :price="item.spreadsOfsell" unit="%" />
        <ParsePrice e :price="item.handlingChargeRate" unit="%" />
        <div class="flex justify-end">
          <EditOutlined @click="editDialog = deepCopy(item)" />
        </div>
      </TableItem>
    </TablePaginationContainer>
    <Loader v-else />
    <Dialog v-model="editDialog" :loading="loading" title="交易對設定">
      <div class="flex flex-col gap-4" v-if="editDialog">
        <div
          class="grid grid-cols-[4rem_1fr_1em] gap-4 items-center justify-end"
        >
          <FormLabel label="交易對" />
          <div>
            {{ editDialog.baseCurrenciesSymbol }}/{{
              editDialog.quoteCurrenciesSymbol
            }}
          </div>
          <div></div>
          <FormLabel label="狀態" />
          <div>
            <a-radio-group
              v-model:value="editDialog.status"
              :options="[
                { label: '開啟', value: 1 },
                { label: '關閉', value: 0 },
              ]"
            />
          </div>
          <div></div>
          <FormLabel label="買幣匯差" />
          <FormInput type="number" v-model="editDialog.spreadsOfBuy" />
          <div>%</div>
          <FormLabel label="賣幣匯差" />
          <FormInput type="number" v-model="editDialog.spreadsOfsell" />
          <div>%</div>
          <FormLabel label="手續費" />
          <FormInput type="number" v-model="editDialog.handlingChargeRate" />
          <div>%</div>
        </div>
      </div>
      <template #actions>
        <Button variant="outline" @click="editDialog = false"> 取消 </Button>
        <Button @click="editTransactionpairs"> 確定 </Button>
      </template>
    </Dialog>
  </div>
</template>
