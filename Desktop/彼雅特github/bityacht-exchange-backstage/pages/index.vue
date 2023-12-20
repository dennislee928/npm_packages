<script setup>
import { QuestionCircleOutlined } from "@ant-design/icons-vue";
import useNotification from "@/hooks/useNotification";
import useAPI from "@/hooks/useAPI";
const toast = useNotification();
const { login, forgotPassword } = useAPI();
definePageMeta({
  layout: "empty",
});
const router = useRouter();
const forgotPasswordDialog = ref(false);
const loading = ref(false);
const form = ref({
  account: "",
  password: "",
});
async function submit() {
  loading.value = true;
  let result = await login(form.value);
  if (result.accessToken) {
    localStorage.setItem("accessToken", result.accessToken);
    router.push("/members");
  } else {
    loading.value = false;
    toast.error("帳號或密碼錯誤");
  }
}
async function forgotPasswordSubmit() {
  loading.value = true;
  let result = await forgotPassword(form.value.account);
  loading.value = false;
  if (result) {
    toast.success("已發送臨時密碼至信箱");
  } else {
    toast.error("發生錯誤");
  }
}
</script>
<template>
  <div class="h-[100svh] w-full bg-[#19253A] flex items-center justify-center">
    <form class="w-[512px] flex flex-col items-center" @submit.prevent="submit">
      <img src="/img/logo.svg" class="mx-auto -my-14 w-[35%]" />
      <div class="text-center text-2xl font-bold text-white">
        交易所管理後台
      </div>
      <div
        class="grid grid-cols-[2rem_1fr] gap-y-4 gap-x-2 items-center justify-end mt-8 w-[350px]"
      >
        <div class="text-white">帳號</div>
        <FormInput v-model="form.account" />
        <div class="text-white">密碼</div>
        <FormInput type="password" v-model="form.password" />
        <FormLabel label="" />
        <div class="px-2">
          <a
            class="flex leading-4 items-center gap-2 text-sm text-white opacity-80 hover:opacity-90 active:opacity-100 cursor-pointer"
            @click="forgotPasswordDialog = true"
          >
            <QuestionCircleOutlined />
            忘記密碼
          </a>
        </div>
      </div>
      <Button
        type="submit"
        :loading="loading"
        class="bg-[#25BBEE] text-white rounded-full px-12 py-3 hover:bg-[#1D9AC6] active:bg-[#12849B] mt-4"
      >
        登入
      </Button>
    </form>
    <Dialog v-model="forgotPasswordDialog" title="忘記密碼" :loading="loading">
      <div class="flex flex-col gap-4 my-8">
        <div class="grid grid-cols-[4rem_1fr] gap-2 items-center justify-end">
          <FormLabel label="E-mail" required />
          <FormInput type="email" v-model="form.account" />
        </div>
        <div class="text-center text-[#394B6A]">
          將發送臨時密碼此信箱，請於登入後重新設定密碼。
        </div>
      </div>
      <template #actions>
        <Button variant="outline" @click="forgotPasswordDialog = false">
          取消
        </Button>
        <Button
          @click="
            forgotPasswordSubmit();
            forgotPasswordDialog = false;
          "
        >
          確定
        </Button>
      </template>
    </Dialog>
  </div>
</template>
