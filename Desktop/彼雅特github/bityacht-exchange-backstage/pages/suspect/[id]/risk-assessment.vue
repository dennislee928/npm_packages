<script setup>
import useAPI from "@/hooks/useAPI";
import useNotification from "@/hooks/useNotification";
const { updateSuspiciousTransactionDetail } = useAPI();
const toast = useNotification();
const route = useRoute();

const { id } = route.params;
const riskReviewResult = ref(0);
onMounted(async () => {
  riskReviewResult.value = data.riskReviewResult;
});
const { data, getData } = defineProps({
  data: {
    type: Object,
  },
  getData: {
    type: Function,
  },
});
async function updateDetail() {
  let result = await updateSuspiciousTransactionDetail(id, {
    riskReviewResult: riskReviewResult.value,
    type: 2,
  });
  if (result?.code) {
    toast.error("更新失敗");
  } else {
    toast.success("更新成功");
    getData();
  }
}
</script>
<template>
  <TableContainer cols="1fr 1fr">
    <TableHeader>
      <div>風控審查</div>
    </TableHeader>
    <TableItemDual>
      <TableSubItem align-top>
        <div>檔案</div>
        <Files :data="data" :getData="getData" file-type="risk-assessment" />
      </TableSubItem>
      <TableSubItem align-top>
        <div>結果</div>
        <div>
          <div class="flex w-full flex-col gap-2">
            <div>
              <a-radio-group
                :options="[
                  {
                    label: '非可疑交易',
                    value: 2,
                  },
                ]"
                v-model:value="riskReviewResult"
              />
              <br />
              <a-radio-group
                :options="[
                  {
                    label: '有可疑交易風險',
                    value: 3,
                  },
                ]"
                v-model:value="riskReviewResult"
              />
            </div>
            <Button class="self-end" @click="updateDetail">送出</Button>
          </div>
        </div>
      </TableSubItem>
    </TableItemDual>
  </TableContainer>
</template>
