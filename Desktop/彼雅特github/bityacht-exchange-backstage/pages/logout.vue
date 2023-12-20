<script setup>
import useAPI from "@/hooks/useAPI";
import useNotification from "@/hooks/useNotification";
const { logout } = useAPI();
const toast = useNotification();
const router = useRouter();
onMounted(async () => {
  await logout();
  localStorage.removeItem("accessToken");
  toast.success("登出成功");
  router.push("/");
});
</script>
<template>
  <div class="flex flex-col items-center justify-center h-screen">
    <Loader />
  </div>
</template>
