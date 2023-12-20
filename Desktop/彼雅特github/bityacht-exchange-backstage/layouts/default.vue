<script setup>
import {
  IdcardOutlined,
  DollarOutlined,
  MoneyCollectOutlined,
  FileProtectOutlined,
  ShakeOutlined,
  FileTextOutlined,
  SettingOutlined,
  UserOutlined,
} from "@ant-design/icons-vue";
import useAPI from "@/hooks/useAPI";
import useUserInfo from "@/hooks/useUserInfo";
import useNotification from "@/hooks/useNotification";
const toast = useNotification();
const router = useRouter();
const { userInfo } = useUserInfo();
const { checkAuth } = useAPI();
onMounted(async () => {
  let result = await checkAuth();
  if (!userInfo || userInfo.isExpired || !result.isValid) {
    router.push("/");
  }
});
router.beforeEach(async (to, from, next) => {
  let result = await checkAuth();
  if (result.isValid || to.path == "/") {
    next();
  } else {
    toast.error("登入狀態已過期，可能是因為您已在其他裝置登入。");
    next("/");
  }
});
</script>
<template>
  <ClientOnly>
    <template #fallback>
      <Loader />
    </template>
    <a-config-provider
      :theme="{
        token: {
          colorPrimary: '#00b96b',
          fontSize: 16,
        },
      }"
    >
      <div
        class="grid grid-cols-[200px_1fr] xl:grid-cols-[240px_1fr] bg-[#F9FAFB] h-[100svh]"
      >
        <div
          class="bg-[#19253A] text-white flex flex-col overflow-y-auto"
          v-if="userInfo"
        >
          <img src="/img/logo.svg" class="mx-auto" />
          <LayoutNavItem href="/members">
            <template #icon>
              <idcard-outlined />
            </template>
            會員列表
          </LayoutNavItem>
          <LayoutNestedNavItems
            v-if="[4, 3, 2, 1].includes(userInfo.managersRolesID)"
          >
            <template #icon>
              <user-outlined />
            </template>
            <template #name> KYC 認證審查 </template>
            <LayoutNestedNavItem href="/kyc" exact>
              審查列表
            </LayoutNestedNavItem>
            <LayoutNestedNavItem href="/kyc/annual" exact>
              年度覆核名單
            </LayoutNestedNavItem>
          </LayoutNestedNavItems>
          <LayoutNavItem href="/order-list">
            <template #icon>
              <dollar-outlined />
            </template>
            交易訂單
          </LayoutNavItem>
          <LayoutNavItem href="/wallet">
            <template #icon>
              <shake-outlined />
            </template>
            提/入幣列表
          </LayoutNavItem>
          <!-- <LayoutNavItem href="/money-list">
            <template #icon>
              <money-collect-outlined />
            </template>
            出/入金列表
          </LayoutNavItem> -->
          <LayoutNavItem
            href="/suspect"
            v-if="[4, 2, 1].includes(userInfo.managersRolesID)"
          >
            <template #icon>
              <file-protect-outlined />
            </template>
            可疑交易管理
          </LayoutNavItem>
          <LayoutNavItem
            href="/receipts"
            v-if="[5, 2, 1].includes(userInfo.managersRolesID)"
          >
            <template #icon>
              <file-text-outlined />
            </template>
            發票管理
          </LayoutNavItem>
          <LayoutNestedNavItems>
            <template #icon>
              <setting-outlined />
            </template>
            <template #name> 平台管理 </template>
            <LayoutNestedNavItem
              href="/settings/trade"
              v-if="[5, 2, 1].includes(userInfo.managersRolesID)"
            >
              交易對/手續費設定
            </LayoutNestedNavItem>
            <LayoutNestedNavItem
              href="/settings/admin"
              v-if="[1].includes(userInfo.managersRolesID)"
            >
              管理者帳號設定
            </LayoutNestedNavItem>
            <LayoutNestedNavItem
              href="/settings/banner"
              v-if="[3, 2, 1].includes(userInfo.managersRolesID)"
            >
              Banner 管理
            </LayoutNestedNavItem>
            <LayoutNestedNavItem
              href="/settings/withdraw"
              v-if="[5, 2, 1].includes(userInfo.managersRolesID)"
            >
              提幣/手續費設定
            </LayoutNestedNavItem>
          </LayoutNestedNavItems>
          <div class="flex-1" />
          <div class="border-t border-gray-700 py-2 px-4">
            <div class="text-gray-300 tracking-tighter">Welcome</div>
            <div class="flex justify-between items-center gap-1">
              <div class="flex items-center gap-1" v-if="userInfo">
                <div class="truncate text-xl tracking-tighter">
                  {{ userInfo.name }}
                </div>
                <ChangePasswordDialog />
              </div>
              <router-link
                class="text-sm active:opacity-90 hover:underline min-w-max tracking-tighter"
                to="/logout"
              >
                登出 ❯
              </router-link>
            </div>
          </div>
        </div>
        <div class="overflow-y-auto tabular-nums">
          <slot />
        </div>
      </div>
    </a-config-provider>
  </ClientOnly>
</template>
