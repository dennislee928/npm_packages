<script setup>
import useAPI from "@/hooks/useAPI";
import useNotification from "@/hooks/useNotification";
import useReviewStatus from "@/hooks/useReviewStatus";
import useUserInfo from "@/hooks/useUserInfo";
import useUserOptions from "@/hooks/useUserOptions";
import useIdvStatus from "@/hooks/useIdvStatus";
import {
  FileOutlined,
  CheckCircleOutlined,
  ExclamationCircleOutlined,
  CloseCircleOutlined,
  DownloadOutlined,
} from "@ant-design/icons-vue";
const {
  getUserKycInfo,
  updateUserKycComplianceReview,
  updateUserKycFinalReview,
  updateUserKycKryptoReview,
  resendUserKycKryptoReview,
  updateKycNameCheck,
  downloadKycNameCheck,
  updateUserKycKryptoTaskID,
  updateUserKycIdvAuditStatus,
} = useAPI();
const route = useRoute();
const toast = useNotification();
const { idvAuditStatus, idvAuditStatusBadge, idvState, idvStateBadge } =
  useIdvStatus();
const { uid } = route.params;
const { status, badgeStatus } = useReviewStatus();
const { result: industrialList, loading: industrialLoading } =
  useUserOptions("ics");
const { result: countryList, loading: countryLoading } =
  useUserOptions("countries");
const { userInfo } = useUserInfo();

const data = ref(null);
const kryptoReviewDialog = ref(false);
const kryptoResendDialog = ref(false);
const kryptoTaskIDDialog = ref(false);
const reviewDialog = ref(false);
const complianceReviewDialog = ref(false);
const finalReviewDialog = ref(false);
const internalRisksDialog = ref(false);
const IdvAuditStatusDialog = ref(false);
const nameCheckDialog = ref(false);
const imageDialog = ref(null);
const loading = ref(false);

const isForeigner = computed(
  () => data.value?.countriesCode !== "TWN" && data.value?.countriesCode !== ""
);

onMounted(() => {
  getData();
});
watch(internalRisksDialog, async () => {
  getData();
});

async function getData() {
  loading.value = true;
  let result = await getUserKycInfo(uid);
  loading.value = false;
  if (!result.error) {
    data.value = result;
  } else {
    toast.error("發生錯誤，請稍後再試。");
  }
}
async function updateComplianceReview() {
  loading.value = true;
  let result = await updateUserKycComplianceReview(
    uid,
    complianceReviewDialog.value.result,
    complianceReviewDialog.value.comment
  );
  if (!result.error) {
    toast.success("更新成功");
    complianceReviewDialog.value = false;
    await getData();
  } else {
    toast.error("更新失敗");
  }
  loading.value = false;
}
async function updateFinalReview() {
  if (!document.getElementById("update-final-review").reportValidity()) {
    return;
  }
  loading.value = true;
  let result = await updateUserKycFinalReview(
    uid,
    finalReviewDialog.value.result,
    finalReviewDialog.value.comment
  );
  if (!result.error) {
    toast.success("更新成功");
    finalReviewDialog.value = false;
    await getData();
  } else {
    toast.error("更新失敗");
  }
  loading.value = false;
}
async function updateIdvAuditStatus() {
  if (!document.getElementById("update-idv-audit-status").reportValidity()) {
    return;
  }
  loading.value = true;
  let result = await updateUserKycIdvAuditStatus(
    uid,
    IdvAuditStatusDialog.value.auditStatus,
    IdvAuditStatusDialog.value.comment
  );
  if (!result.error) {
    toast.success("更新成功");
    IdvAuditStatusDialog.value = false;
    await getData();
  } else {
    toast.error("更新失敗");
  }
  loading.value = false;
}
async function updateKrptoReview() {
  if (!document.getElementById("update-krypto-review").reportValidity()) {
    return;
  }
  loading.value = true;
  let result = await updateUserKycKryptoReview(
    uid,
    kryptoReviewDialog.value.result,
    kryptoReviewDialog.value.comment
  );
  if (!result.error) {
    toast.success("更新成功");
    kryptoReviewDialog.value = false;
    await getData();
  } else {
    toast.error("更新失敗");
  }
  loading.value = false;
}
async function resendKryptoReview() {
  loading.value = true;
  let result = await resendUserKycKryptoReview(uid);
  if (!result.error) {
    toast.success("重送成功");
    kryptoResendDialog.value = false;
    await getData();
  } else {
    toast.error("重送失敗");
  }
  loading.value = false;
}
async function updateNameCheck() {
  // chcek form
  let formData = new FormData();
  {
    let { file, result } = nameCheckDialog.value;
    if (!file || !result) {
      toast.error("請填寫完整");
      return;
    }
    formData.append("file", file);
    formData.append("result", result);
  }
  loading.value = true;
  let result = await updateKycNameCheck(uid, formData);
  loading.value = false;
  if (!result.error) {
    toast.success("更新成功");
    nameCheckDialog.value = false;
    await getData();
  } else {
    toast.error("更新失敗");
  }
}
async function downloadNameCheckPdf() {
  let buffer = await downloadKycNameCheck(uid);
  let blob = new Blob([buffer], { type: "application/pdf" });
  let url = window.URL.createObjectURL(blob);
  let a = document.createElement("a");
  a.href = url;
  a.download = `${data.value.firstName}_${data.value.lastName}_姓名檢核排除評估.pdf`;
  a.click();
  toast.success("下載成功");
}
async function uploadKryptoTaskIDDialog() {
  loading.value = true;
  let result = await updateUserKycKryptoTaskID(
    uid,
    kryptoTaskIDDialog.value.taskID
  );
  if (!result.error) {
    toast.success("更新成功");
    finalReviewDialog.value = false;
    await getData();
  } else {
    toast.error("更新失敗");
  }
  loading.value = false;
}
</script>
<template>
  <TableContainer v-if="data" cols="10em 1fr 5em">
    <TableHeader>
      <div>審查資料</div>
      <div></div>
      <Button size="mini" @click="reviewDialog = true"> 異動紀錄 </Button>
    </TableHeader>
    <TableItem>
      <div>
        {{ !isForeigner ? "身分證字號" : "居留證號碼" }}
      </div>
      <div>{{ data.nationalID }}</div>
    </TableItem>
    <TableItem v-if="data.passportNumber != ''">
      <div>護照號碼</div>
      <div>{{ data.passportNumber }}</div>
    </TableItem>
    <TableItem>
      <div>姓</div>
      <div>{{ data.lastName }}</div>
    </TableItem>
    <TableItem>
      <div>名</div>
      <div>{{ data.firstName }}</div>
    </TableItem>
    <TableItem>
      <div>{{ !isForeigner ? "證件" : "居留證" }}</div>
      <div class="flex gap-2">
        <img
          :src="data.idImage"
          @click="imageDialog = data.idImage"
          v-if="data.idImage != ''"
          class="h-16 cursor-pointer rounded-sm hover:opacity-80"
        />
        <img
          :src="data.idBackImage"
          @click="imageDialog = data.idBackImage"
          v-if="data.idBackImage != ''"
          class="h-16 cursor-pointer rounded-sm hover:opacity-80"
        />
        <img
          :src="data.faceImage"
          @click="imageDialog = data.faceImage"
          v-if="data.faceImage != ''"
          class="h-16 cursor-pointer rounded-sm hover:opacity-80"
        />
        <img
          :src="data.idAndFaceImage"
          @click="imageDialog = data.idAndFaceImage"
          v-if="data.idAndFaceImage != ''"
          class="h-16 cursor-pointer rounded-sm hover:opacity-80"
        />
      </div>
    </TableItem>
    <TableItem v-if="data.passportImage != ''">
      <div>護照</div>
      <div class="flex gap-2">
        <img
          :src="data.passportImage"
          @click="imageDialog = data.passportImage"
          class="h-16 cursor-pointer rounded-sm hover:opacity-80"
        />
      </div>
    </TableItem>
    <TableItem v-if="isForeigner">
      <div v-if="isForeigner">驗證結果</div>
      <div>
        <StatusBadge :status="idvAuditStatusBadge[data.idvAuditStatus]">
          {{ idvAuditStatus[data.idvAuditStatus] }}
        </StatusBadge>
      </div>
      <Button
        size="mini"
        @click="
          IdvAuditStatusDialog = {
            auditStatus: data.idvAuditStatus,
            comment: '',
          }
        "
      >
        修改
      </Button>
    </TableItem>
    <template v-if="!isForeigner">
      <TableHeader>
        <div>IDV</div>
      </TableHeader>
      <TableItem>
        <div>單號</div>
        <div>{{ data.idvTaskID }}</div>
      </TableItem>
      <TableItem>
        <div>自動驗證狀態</div>
        <div>
          <StatusBadge :status="idvStateBadge[data.idvState]">
            {{ idvState[data.idvState] }}
          </StatusBadge>
        </div>
      </TableItem>
      <TableItem>
        <div>驗證結果</div>
        <div>
          <StatusBadge :status="idvAuditStatusBadge[data.idvAuditStatus]">
            {{ idvAuditStatus[data.idvAuditStatus] }}
          </StatusBadge>
        </div>
      </TableItem>
    </template>
    <TableHeader>
      <div>KryptoGO</div>
    </TableHeader>
    <TableItem>
      <div>KryptoGO 單號</div>
      <div>{{ data.kryptoID }}</div>
      <Button size="mini" @click="kryptoTaskIDDialog = { taskID: '' }">
        編輯
      </Button>
    </TableItem>
    <TableItem>
      <div>KryptoGO 風險評分</div>
      <div>{{ data.kryptoPotentialRisk }}</div>
      <Button size="mini" @click="kryptoResendDialog = true"> 重送 </Button>
    </TableItem>
    <TableItem>
      <div>管制名單</div>
      <div>{{ ["未知", "未命中", "命中"][data.kryptoSanctionMatched] }}</div>
    </TableItem>
    <TableItem>
      <div>KryptoGO 複核</div>
      <div>{{ ["未複核", "拒絕", "通過"][data.kryptoAuditAccepted] }}</div>
      <Button
        size="mini"
        @click="
          kryptoReviewDialog = {
            result: data.reviewKrypto,
            comment: '',
          }
        "
        v-show="
          [1, 2, 4].includes(userInfo.managersRolesID) &&
          data.kryptoID != '' &&
          data.kryptoAuditAccepted == 0
        "
      >
        審查
      </Button>
    </TableItem>
    <TableHeader>
      <div>姓名檢核排除評估</div>
    </TableHeader>
    <TableItem>
      <div>姓名檢核排除評估</div>
      <div class="flex items-center gap-8">
        <StatusBadge :status="badgeStatus[data.nameCheck]">
          {{ status[data.nameCheck] }}
        </StatusBadge>
        <button
          @click="downloadNameCheckPdf"
          class="text-[#0F62AE] hover:underline active:opacity-80"
          v-if="data.nameCheckPdfName != ''"
        >
          <FileOutlined />
          下載 PDF
        </button>
      </div>
      <Button
        size="mini"
        @click="nameCheckDialog = { result: null, file: null }"
        v-show="[1, 2, 4].includes(userInfo.managersRolesID)"
      >
        上傳
      </Button>
    </TableItem>
    <TableHeader>
      <div>內部風險審核</div>
    </TableHeader>
    <TableItem>
      <div>內部風險審核</div>
      <div
        v-if="data.internalRisksTotal <= 7"
        class="text-[#02C879] flex items-center gap-1"
      >
        <CheckCircleOutlined />{{ data.internalRisksTotal }}
      </div>
      <div
        v-else-if="data.internalRisksTotal <= 11"
        class="text-[#FF9D18] flex items-center gap-1"
      >
        <ExclamationCircleOutlined />{{ data.internalRisksTotal }}
      </div>
      <div
        v-else-if="data.internalRisksTotal >= 12"
        class="text-[#FF574C] flex items-center gap-1"
      >
        <CloseCircleOutlined />{{ data.internalRisksTotal }}
      </div>
      <Button
        size="mini"
        @click="internalRisksDialog = true"
        v-show="[1, 2, 4].includes(userInfo.managersRolesID)"
      >
        審查
      </Button>
    </TableItem>
    <TableHeader>
      <div>法遵審查</div>
    </TableHeader>
    <TableItem>
      <div>法遵審查</div>
      <StatusBadge :status="badgeStatus[data.complianceReview]">
        {{ status[data.complianceReview] }}
      </StatusBadge>
      <Button
        size="mini"
        @click="
          complianceReviewDialog = {
            result: data.complianceReview,
            comment: '',
          }
        "
        v-show="[1, 2].includes(userInfo.managersRolesID)"
      >
        審查
      </Button>
    </TableItem>
    <TableItem>
      <div>法遵審查備註</div>
      <div>{{ data.complianceReviewComment }}</div>
    </TableItem>
    <TableHeader>
      <div>最終審查</div>
    </TableHeader>
    <TableItem>
      <div>最終審查</div>
      <StatusBadge :status="badgeStatus[data.finalReview]">
        {{ status[data.finalReview] }}
      </StatusBadge>
      <Button
        size="mini"
        @click="
          finalReviewDialog = {
            result: data.finalReview,
            notice: '',
            comment: '',
          }
        "
        v-show="[1, 2].includes(userInfo.managersRolesID)"
      >
        審查
      </Button>
    </TableItem>
    <TableItem>
      <div>最終審查通知訊息</div>
      <div>{{ data.finalReviewNotice }}</div>
    </TableItem>
    <TableItem>
      <div>最終審查原由備註</div>
      <div>{{ data.finalReviewComment }}</div>
    </TableItem>
  </TableContainer>
  <TableContainer v-if="data">
    <TableHeader>
      <div>會員資料</div>
      <div>資料</div>
    </TableHeader>
    <TableItem>
      <div>手機</div>
      <div>{{ data.phone }}</div>
    </TableItem>
    <TableItem>
      <div>生日</div>
      <div>{{ data.birthDate }}</div>
    </TableItem>
    <TableItem>
      <div>國家</div>
      <div v-if="!countryLoading && data.countriesCode != ''">
        {{
          countryList.filter(({ code }) => code === data.countriesCode)[0]
            .chinese
        }}
      </div>
      <div class="opacity-50" v-else>-</div>
    </TableItem>
    <TableItem>
      <div>雙重國籍</div>
      <div v-if="!countryLoading && data.dualNationalityCode != ''">
        {{
          countryList.filter(({ code }) => code === data.dualNationalityCode)[0]
            .chinese
        }}
      </div>
      <div class="opacity-50" v-else>-</div>
    </TableItem>
    <TableItem>
      <div>地址</div>
      <div>{{ data.address }}</div>
    </TableItem>
    <TableItem>
      <div>行業別</div>
      <div v-if="data.industrialClassificationsID != 0 && !industrialLoading">
        {{
          industrialList.filter(
            (x) => x.id === data.industrialClassificationsID
          )[0]?.chinese
        }}
      </div>
      <div v-else class="opacity-50">-</div>
    </TableItem>
    <TableItem>
      <div>年收入</div>
      <div>{{ data.annualIncome }}</div>
    </TableItem>
    <TableItem>
      <div>資金來源</div>
      <div>{{ data.fundsSources }}</div>
    </TableItem>
    <TableItem>
      <div>使用目的</div>
      <div>{{ data.purposeOfUse }}</div>
    </TableItem>
    <TableItem>
      <div>申請時間</div>
      <div>{{ data.createdAt }}</div>
    </TableItem>
  </TableContainer>
  <Loader v-if="loading" />
  <Dialog
    v-model="reviewDialog"
    title="KYC 異動紀錄"
    max-width="1100px"
    variant="sidebar"
  >
    <KycReviewlog v-if="reviewDialog" />
    <template #actions>
      <Button @click="reviewDialog = false"> 關閉 </Button>
    </template>
  </Dialog>
  <KycInternalRisks v-model="internalRisksDialog" />
  <Dialog v-model="complianceReviewDialog" title="法遵審查" :loading="loading">
    <form
      @submit.prevent
      class="grid grid-cols-[5rem_1fr] gap-4 items-center justify-end"
      id="update-status-form"
      v-if="complianceReviewDialog"
    >
      <FormLabel label="狀態" />
      <div>
        <a-radio-group
          v-model:value="complianceReviewDialog.result"
          :options="[
            {
              label: '通過',
              value: 2,
            },
            {
              label: '拒絕',
              value: 3,
            },
          ]"
        />
      </div>
      <FormLabel label="緣由備註" required />
      <FormInput v-model="complianceReviewDialog.comment" required />
    </form>
    <template #actions>
      <Button variant="outline" @click="complianceReviewDialog = false">
        取消
      </Button>
      <Button @click="updateComplianceReview"> 確定 </Button>
    </template>
  </Dialog>
  <Dialog v-model="finalReviewDialog" title="最終審查" :loading="loading">
    <form
      @submit.prevent
      class="grid grid-cols-[5rem_1fr] gap-4 items-center justify-end"
      id="update-final-review"
      v-if="finalReviewDialog"
    >
      <FormLabel label="狀態" />
      <div>
        <a-radio-group v-model:value="finalReviewDialog.result">
          <a-radio
            class="block"
            v-for="item of [
              {
                label: '通過',
                value: 2,
              },
              {
                label: '拒絕 - 風險考量',
                value: 3,
              },
              {
                label: '拒絕 - 證件不符',
                value: 4,
              },
            ]"
            :value="item.value"
          >
            {{ item.label }}
          </a-radio>
        </a-radio-group>
      </div>
      <FormLabel label="緣由備註" required />
      <FormInput v-model="finalReviewDialog.comment" required />
    </form>
    <template #actions>
      <Button variant="outline" @click="finalReviewDialog = false">
        取消
      </Button>
      <Button @click="updateFinalReview"> 確定 </Button>
    </template>
  </Dialog>
  <Dialog v-model="kryptoReviewDialog" title="KryptoGO 複核" :loading="loading">
    <form
      @submit.prevent
      class="grid grid-cols-[5rem_1fr] gap-4 items-center justify-end"
      id="update-krypto-review"
      v-if="kryptoReviewDialog"
    >
      <FormLabel label="狀態" />
      <div>
        <a-radio-group
          v-model:value="kryptoReviewDialog.result"
          :options="[
            {
              label: '通過',
              value: 2,
            },
            {
              label: '拒絕',
              value: 3,
            },
          ]"
        />
      </div>
      <FormLabel label="備註" />
      <FormInput v-model="kryptoReviewDialog.comment" />
    </form>
    <template #actions>
      <Button variant="outline" @click="kryptoReviewDialog = false">
        取消
      </Button>
      <Button @click="updateKrptoReview"> 確定 </Button>
    </template>
  </Dialog>
  <Dialog
    v-model="IdvAuditStatusDialog"
    title="更新認證結果（外籍人士）"
    :loading="loading"
  >
    <form
      @submit.prevent
      class="grid grid-cols-[5rem_1fr] gap-4 items-center justify-end"
      id="update-idv-audit-status"
      v-if="IdvAuditStatusDialog"
    >
      <FormLabel label="狀態" />
      <div>
        <a-radio-group
          v-model:value="IdvAuditStatusDialog.auditStatus"
          :options="[
            {
              label: '通過',
              value: 2,
            },
            {
              label: '拒絕',
              value: 3,
            },
          ]"
        />
      </div>
      <FormLabel label="備註" required />
      <FormInput v-model="IdvAuditStatusDialog.comment" required />
    </form>
    <template #actions>
      <Button variant="outline" @click="IdvAuditStatusDialog = false">
        取消
      </Button>
      <Button @click="updateIdvAuditStatus"> 確定 </Button>
    </template>
  </Dialog>
  <Dialog v-model="nameCheckDialog" title="姓名檢核排除評估" :loading="loading">
    <form
      @submit.prevent
      class="grid grid-cols-[5rem_1fr] gap-4 items-center justify-end"
      id="name-check-form"
      v-if="nameCheckDialog"
    >
      <FormLabel label="姓名檢核" />
      <div>
        <a-radio-group
          v-model:value="nameCheckDialog.result"
          :options="[
            {
              label: '通過',
              value: 2,
            },
            {
              label: '拒絕',
              value: 3,
            },
          ]"
        />
      </div>
      <FormLabel label="檔案" />
      <FormInput
        type="file"
        id="name-check-pdf"
        accept="application/pdf"
        @change="
          (e) => {
            nameCheckDialog.file = e.target.files[0];
          }
        "
      />
    </form>
    <template #actions>
      <Button variant="outline" @click="nameCheckDialog = false"> 取消 </Button>
      <Button @click="updateNameCheck"> 確定 </Button>
    </template>
  </Dialog>
  <Dialog v-model="imageDialog" :loading="loading" variant="image">
    <img :src="imageDialog" class="w-full" />
    <template #actions>
      <a :href="imageDialog" download class="flex-1">
        <Button variant="outline"> <DownloadOutlined /> 下載 </Button>
      </a>
      <div class="flex-1"></div>
      <Button @click="imageDialog = false"> 關閉 </Button>
    </template>
  </Dialog>
  <Dialog v-model="kryptoResendDialog" :loading="loading" title="重送 KryptoGo">
    <p>
      確定要重新將會員 <span class="text-[#0F62AE]">{{ uid }}</span> <br />
      重新送審 KryptoGO 嗎？
    </p>
    <template #actions>
      <div class="flex-1"></div>
      <Button variant="outline" @click="kryptoResendDialog = false">
        取消
      </Button>
      <Button @click="resendKryptoReview"> 送出 </Button>
    </template>
  </Dialog>
  <Dialog v-model="kryptoTaskIDDialog" title="KryptoGO 單號" :loading="loading">
    <form
      @submit.prevent
      class="flex flex-col gap-4"
      id="update-kryptogo-task-id-form"
    >
      <div class="grid grid-cols-[8rem_1fr] gap-4 items-center justify-end">
        <FormLabel label="KryptoGO 單號" />
        <FormInput v-model="kryptoTaskIDDialog.taskID" />
      </div>
    </form>
    <template #actions>
      <Button variant="outline" @click="kryptoTaskIDDialog = false">
        取消
      </Button>
      <Button @click="uploadKryptoTaskIDDialog"> 確定 </Button>
    </template>
  </Dialog>
</template>
