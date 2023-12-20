<script setup>
import useAPI from "@/hooks/useAPI";
import useNotification from "@/hooks/useNotification";
const { updateSuspiciousTransactionDetail } = useAPI();
const toast = useNotification();
const route = useRoute();

const { id } = route.params;
const dedicatedReviewResult = ref(0);
const comment = ref("");
const { data, getData } = defineProps({
  data: {
    type: Object,
  },
  getData: {
    type: Function,
  },
});
onMounted(async () => {
  dedicatedReviewResult.value = data.dedicatedReviewResult;
  comment.value = data.dedicatedReviewComment;
});
async function updateDetail() {
  let result = await updateSuspiciousTransactionDetail(id, {
    dedicatedReviewResult: dedicatedReviewResult.value,
    comment: comment.value,
    type: 3,
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
      <div>專責審查</div>
    </TableHeader>
    <TableItemDual>
      <TableSubItem align-top>
        <div>結果</div>
        <div>
          <textarea
            class="w-full h-32 border border-gray-300 rounded-md p-2"
            placeholder="請輸入備註"
            v-model="comment"
          />
        </div>
      </TableSubItem>
      <TableSubItem align-top>
        <div>備註</div>
        <div>
          <div class="flex w-full flex-col gap-2">
            <div>
              <a-radio-group
                :options="[
                  {
                    label: '通過',
                    value: 2,
                  },
                ]"
                v-model:value="dedicatedReviewResult"
              />
              <br />
              <a-radio-group
                :options="[
                  {
                    label: '駁回',
                    value: 3,
                  },
                ]"
                v-model:value="dedicatedReviewResult"
              />
            </div>
            <Button class="self-end" @click="updateDetail">送出</Button>
          </div>
        </div>
      </TableSubItem>
    </TableItemDual>
  </TableContainer>
</template>
