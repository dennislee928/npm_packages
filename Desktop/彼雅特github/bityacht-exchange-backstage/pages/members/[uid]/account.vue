<script setup>
import useAPI from "@/hooks/useAPI";
import useNotification from "@/hooks/useNotification";
import useUserInfo from "@/hooks/useUserInfo";
const route = useRoute();
const toast = useNotification();
const { getUserBank, updateUserBank } = useAPI();
const { uid } = route.params;
const { userInfo } = useUserInfo();
const bankInfo = ref(null);
const updateStatusDialog = ref(false);
const statuslogDialog = ref(false);
const imageDialog = ref(false);
const loading = ref(false);
onMounted(async () => {
  await getData();
});
async function getData() {
  let result = await getUserBank(uid);
  bankInfo.value = result;
}
async function updateStatus() {
  if (document.getElementById("update-status-form").reportValidity()) {
    loading.value = true;
    let result = await updateUserBank(uid, {
      ...updateStatusDialog.value,
      id: bankInfo.value.id,
    });
    if (!result.error) {
      updateStatusDialog.value = false;
      await getData();
    }
    loading.value = false;
  }
}
</script>
<template>
  <TableContainer v-if="bankInfo">
    <TableHeader>
      <div>欄位</div>
      <div class="flex items-center justify-between">
        <span> 資料 </span>
        <Button size="mini" @click="statuslogDialog = true"> 異動紀錄 </Button>
      </div>
    </TableHeader>
    <TableItem>
      <div>帳戶狀態</div>
      <div class="flex items-center justify-between max-w-[200px]">
        <StatusBadge
          :status="
            ['in-progress', 'in-progress', 'completed', 'cancelled'][
              bankInfo.status
            ]
          "
        >
          {{ ["未綁定", "審查中", "已通過", "未通過"][bankInfo.status] }}
        </StatusBadge>
        <Button
          size="mini"
          @click="updateStatusDialog = { status: bankInfo.status, comment: '' }"
          v-show="
            [4, 2, 1].includes(userInfo.managersRolesID) &&
            bankInfo.status !== 0
          "
        >
          修改
        </Button>
      </div>
    </TableItem>
    <template v-if="bankInfo.status !== 0">
      <TableItem>
        <div>銀行</div>
        <div>
          {{ bankInfo.bankInfo.code }}（{{ bankInfo.bankInfo.chinese }}）
        </div>
      </TableItem>
      <TableItem>
        <div>分行</div>
        <div>
          {{ bankInfo.branchInfo.code }}（{{ bankInfo.branchInfo.chinese }}）
        </div>
      </TableItem>
      <TableItem>
        <div>銀行帳戶</div>
        <div>{{ bankInfo.account }}</div>
      </TableItem>
      <TableItem>
        <div>帳戶名稱</div>
        <div>{{ bankInfo.name }}</div>
      </TableItem>
      <TableItem v-if="bankInfo.coverImage && bankInfo.coverImage != ''">
        <div>封面照片</div>
        <div class="flex gap-2">
          <img
            :src="bankInfo.coverImage"
            @click="imageDialog = bankInfo.coverImage"
            class="h-16 cursor-pointer rounded-sm hover:opacity-80"
          />
        </div>
      </TableItem>
      <TableItem>
        <div>綁定時間</div>
        <div>{{ bankInfo.createdAt }}</div>
      </TableItem>
      <TableItem>
        <div>審核時間</div>
        <div>{{ bankInfo.auditTime }}</div>
      </TableItem>
    </template>
  </TableContainer>
  <Loader v-if="!bankInfo" />
  <Dialog v-model="updateStatusDialog" :loading="loading" title="修改狀態">
    <form
      @submit.prevent
      class="grid grid-cols-[4rem_1fr] gap-4 items-center justify-end"
      id="update-status-form"
      v-if="updateStatusDialog"
    >
      <FormLabel label="狀態" />
      <div>
        <a-radio-group
          v-model:value="updateStatusDialog.status"
          :options="
            [`已通過`, `未通過`].map((key, index) => ({
              label: key,
              value: index + 2,
            }))
          "
        />
      </div>
      <FormLabel label="緣由備註" />
      <FormInput v-model="updateStatusDialog.comment" />
    </form>
    <template #actions>
      <Button variant="outline" @click="updateStatusDialog = false">
        取消
      </Button>
      <Button @click="updateStatus"> 確定 </Button>
    </template>
  </Dialog>
  <Dialog
    v-model="statuslogDialog"
    title="異動紀錄"
    max-width="1100px"
    variant="sidebar"
  >
    <UserAccountLog v-if="statuslogDialog" />
    <template #actions>
      <Button @click="statuslogDialog = false"> 關閉 </Button>
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
</template>
