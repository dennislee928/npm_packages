<script setup>
import useAPI from "@/hooks/useAPI";
import useNotification from "@/hooks/useNotification";
import { SettingOutlined } from "@ant-design/icons-vue";

const { changePassword } = useAPI();
const toast = useNotification();

const passwordDialog = ref(false);
const loading = ref(false);
const newPassword = ref("");
const confirmPassword = ref("");

async function updatePassword() {
  if (!document.getElementById("update-password-form").reportValidity()) {
    return;
  }
  if (newPassword.value != confirmPassword.value) {
    toast.info("兩次密碼不一致");
    return;
  }
  if (
    !newPassword.value.match(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$/)
  ) {
    toast.info("密碼格式錯誤");
    return;
  }
  loading.value = true;
  let result = await changePassword(newPassword.value);
  if (!result.error) {
    passwordDialog.value = false;
    newPassword.value = "";
    confirmPassword.value = "";
    toast.success("修改成功");
  } else {
    toast.error("修改失敗");
  }
  loading.value = false;
}
</script>
<template>
  <div>
    <Dialog v-model="passwordDialog" :loading="loading" title="修改密碼">
      <form
        @submit.prevent
        class="grid grid-cols-[4rem_1fr] gap-4 items-center justify-end"
        id="update-password-form"
      >
        <FormLabel label="新密碼" />
        <FormInput v-model="newPassword" required type="password" />
        <div />

        <div class="text-sm opacity-75 -mt-2">
          密碼至少 8
          碼，必須包含至少一個英文大寫字母、一個英文小寫字母以及一個數字。
        </div>
        <FormLabel label="確認密碼" />
        <FormInput v-model="confirmPassword" required type="password" />
      </form>
      <template #actions>
        <Button variant="outline" @click="passwordDialog = false">
          取消
        </Button>
        <Button @click="updatePassword"> 確定 </Button>
      </template>
    </Dialog>
    <a-tooltip :mouseEnterDelay="0" :mouseLeaveDelay="0" color="white">
      <template #title>
        <span class="text-gray-950"> 修改密碼 </span>
      </template>

      <div
        class="text-base flex items-center justify-center p-1 rounded bg-white bg-opacity-0 hover:bg-opacity-10 active:bg-opacity-30 transition-colors cursor-pointer hover:shadow"
        @click="passwordDialog = true"
      >
        <SettingOutlined />
      </div>
    </a-tooltip>
  </div>
</template>
