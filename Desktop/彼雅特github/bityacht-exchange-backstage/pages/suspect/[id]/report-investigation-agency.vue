<script setup>
import useAPI from "@/hooks/useAPI";
import useNotification from "@/hooks/useNotification";
import useParseDate from "@/hooks/useParseDate";
const { updateSuspiciousTransactionDetail } = useAPI();
const toast = useNotification();
const route = useRoute();
const parseDate = useParseDate();

const { id } = route.params;
const reportMJIBAt = ref("");
const { data, getData } = defineProps({
  data: {
    type: Object,
  },
  getData: {
    type: Function,
  },
});
onMounted(async () => {
  reportMJIBAt.value = data.reportMJIBAt.replace(" ", "T").replaceAll("/", "-");
});
async function updateDetail() {
  let result = await updateSuspiciousTransactionDetail(id, {
    reportMJIBAt: parseDate(reportMJIBAt.value),
    type: 4,
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
      <div>呈報調查局</div>
    </TableHeader>
    <TableItemDual>
      <TableSubItem align-top>
        <div>檔案</div>
        <Files
          :data="data"
          :getData="getData"
          file-type="report-investigation-agency"
        />
      </TableSubItem>
      <TableSubItem align-top>
        <div>呈報日期</div>
        <div>
          <div class="flex w-full flex-col gap-2">
            <FormInput type="datetime-local" v-model="reportMJIBAt" />
            {{ data.reportMJIBAt }}
            <br />
            <Button class="self-end" @click="updateDetail">送出</Button>
          </div>
        </div>
      </TableSubItem>
    </TableItemDual>
  </TableContainer>
</template>
