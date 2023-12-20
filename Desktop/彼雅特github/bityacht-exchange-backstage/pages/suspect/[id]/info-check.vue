<script setup>
import useAPI from "@/hooks/useAPI";
import useNotification from "@/hooks/useNotification";
const { data, getData } = defineProps({
  data: {
    type: Object,
  },
  getData: {
    type: Function,
  },
});
const route = useRoute();
const toast = useNotification();
const { id } = route.params;
const { updateSuspiciousTransactionDetail } = useAPI();
const informationReviewComment = ref("");
onMounted(async () => {
  informationReviewComment.value = data.informationReviewComment;
});
async function updateDetail() {
  let result = await updateSuspiciousTransactionDetail(id, {
    comment: informationReviewComment.value,
    type: 1,
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
      <div>資訊審核</div>
    </TableHeader>
    <TableItemDual>
      <TableSubItem align-top>
        <div>檔案</div>
        <Files :data="data" :getData="getData" file-type="info-check" />
      </TableSubItem>
      <TableSubItem align-top>
        <div>備註</div>
        <div>
          <div class="flex w-full flex-col gap-2">
            <textarea
              class="w-full h-32 border border-gray-300 rounded-md p-2"
              placeholder="請輸入備註"
              v-model="informationReviewComment"
            />
            <Button class="self-end" @click="updateDetail">送出</Button>
          </div>
        </div>
      </TableSubItem>
    </TableItemDual>
  </TableContainer>
</template>
