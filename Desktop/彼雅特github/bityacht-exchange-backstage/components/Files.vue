<script setup>
import {
  FileOutlined,
  PlusOutlined,
  DeleteOutlined,
} from "@ant-design/icons-vue";
import useAPI from "@/hooks/useAPI";
import useNotification from "@/hooks/useNotification";

const route = useRoute();
const toast = useNotification();
const { id } = route.params;
const {
  downloadSuspiciousTransactionFiles,
  uploadSuspiciousTransactionFiles,
  deleteSuspiciousTransactionFiles,
} = useAPI();
const uploadFileDialog = ref(false);
const loading = ref(false);
const file = ref(null);
const { data, getData, fileType } = defineProps({
  data: {
    type: Object,
  },
  getData: {
    type: Function,
  },
  // fileType: "info-check" | "risk-assessment" | "report-investigation-agency"
  fileType: {
    type: String,
    default: "",
  },
});

async function uploadFile() {
  let uploadType = 0;
  if (fileType === "info-check") {
    uploadType = 1;
  }
  if (fileType === "risk-assessment") {
    uploadType = 2;
  }
  if (fileType === "report-investigation-agency") {
    toast.info("無法在此上傳檔案");
    return;
  }
  let formData = new FormData();
  formData.append("id", id);
  formData.append("uploadType", uploadType);
  formData.append("file", file.value);
  loading.value = true;
  let result = await uploadSuspiciousTransactionFiles(formData);
  loading.value = false;
  if (!result.error) {
    toast.success("上傳成功");
    await getData();
    uploadFileDialog.value = false;
  } else {
    toast.error("上傳失敗");
  }
}
async function deleteFile(file) {
  if (confirm("確定要刪除嗎？")) {
    let result = await deleteSuspiciousTransactionFiles({
      id,
      filename: file.filename,
      fileType: file.type,
    });
    if (!result.error) {
      toast.success("刪除成功");
      await getData();
    } else {
      toast.error("刪除失敗");
    }
  }
}
async function downloadFile(file) {
  let buffer = await downloadSuspiciousTransactionFiles({
    id,
    filename: file.filename,
    fileType: file.type,
  });
  let blob = new Blob([buffer]);
  let url = window.URL.createObjectURL(blob);
  let a = document.createElement("a");
  a.href = url;
  a.download = file.filename;
  a.click();
  toast.success("下載成功");
}
</script>
<template>
  <div class="flex flex-col gap-2 items-start justify-start w-full">
    <div
      class="flex items-center gap-2"
      v-for="item of [
        ...(fileType !== 'risk-assessment' && data.informationReviewFiles
          ? data.informationReviewFiles.map((x) => ({
              filename: x,
              type: 1,
            }))
          : []),
        ...(fileType !== 'info-check' && data.riskReviewFiles
          ? data.riskReviewFiles.map((x) => ({
              filename: x,
              type: 2,
            }))
          : []),
      ]"
    >
      <div class="flex items-center gap-2">
        <FileOutlined />
        <div
          class="link break-all"
          @click="downloadFile(item)"
          :title="`${item.filename}&#10;${
            item.type === 1 ? '資訊審核' : '風控審查'
          }檔案`"
        >
          {{ item.filename }}
        </div>
      </div>
      <button
        class="text-sm text-[#7589A4] h-6 w-6 rounded-full flex items-center justify-center cursor-pointer hover:bg-[#7589A4] hover:text-white hover:bg-opacity-80 active:bg-opacity-100"
        @click="deleteFile(item)"
      >
        <DeleteOutlined />
      </button>
    </div>
    <Button
      size="mini"
      variant="outline"
      @click="uploadFileDialog = true"
      v-if="
        fileType != 'report-investigation-agency' &&
        ((fileType === 'info-check' &&
          (data.informationReviewFiles?.length || 0) < 10) ||
          (fileType === 'risk-assessment' && (data.riskReviewFiles?.length || 0) < 10))
      "
    >
      <PlusOutlined />上傳檔案
    </Button>
  </div>

  <Dialog v-model="uploadFileDialog" title="上傳檔案" :loading="loading">
    <form
      @submit.prevent
      class="grid grid-cols-[4rem_1fr] gap-x-4 gap-y-2 items-center justify-end"
      id="import-form"
      v-if="uploadFileDialog"
    >
      <FormLabel label="檔案" />
      <FormInput
        type="file"
        id="file"
        @change="
          (e) => {
            file = e.target.files[0];
          }
        "
        required
      />
    </form>
    <template #actions>
      <Button variant="outline" @click="uploadFileDialog = false">
        取消
      </Button>
      <Button @click="uploadFile"> 上傳 </Button>
    </template>
  </Dialog>
</template>
