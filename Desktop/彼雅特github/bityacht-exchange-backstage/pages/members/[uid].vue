<script setup>
import useAPI from "@/hooks/useAPI";
import useUserInfo from "@/hooks/useUserInfo";
import useUserOptions from "@/hooks/useUserOptions";
import { twMerge } from "tailwind-merge";

const { getUserInfo } = useAPI();
const { userInfo } = useUserInfo();
const route = useRoute();
const { result: countryList, loading: countryLoading } =
  useUserOptions("countries");

const { uid } = route.params;
const memberInfo = ref(null);

onMounted(async () => {
  let result = await getUserInfo(uid);
  memberInfo.value = result;
});
</script>
<template>
  <div
    class="flex gap-2 py-2 px-4 items-center justify-center text-[#394B6A] mt-4"
  >
    <div class="flex items-center gap-3 justify-center">
      <NuxtLink
        to="/members"
        class="bg-[#394B6A] text-white py-1 px-3 rounded-full hover:opacity-80 active:opacity-60"
      >
        ❮ 返回
      </NuxtLink>
      <div class="flex items-center gap-3 justify-center" v-if="memberInfo">
        <div class="text-3xl font-bold truncate">
          {{ memberInfo.lastName }} {{ memberInfo.firstName }}
        </div>
        <div class="">UID: {{ memberInfo.id }}</div>
        <div
          :class="
            twMerge(`border px-2 rounded-sm`, 'border-[#52BBE8] text-[#52BBE8]')
          "
          v-if="memberInfo.type === 2"
        >
          法人
        </div>
        <div
          :class="
            twMerge(
              `border px-2 rounded-sm`,
              memberInfo.countriesCode === 'TWN'
                ? 'border-[#394B6A] text-[#394B6A]'
                : 'border-[#DD4949] text-[#DD4949]'
            )
          "
          v-else-if="!countryLoading && memberInfo.countriesCode != ''"
        >
          {{
            countryList.filter((x) => x.code === memberInfo.countriesCode)[0]
              ?.chinese
          }}
        </div>
      </div>
      <div class="text-3xl font-bold opacity-0" v-else>Loading</div>
    </div>
    <div class="flex-1"></div>
    <NestedNavItem :href="`/members/${uid}/info`"> 會員資料 </NestedNavItem>
    <NestedNavItem
      :href="`/members/${uid}/kyc`"
      v-if="
        [4, 2, 1].includes(userInfo.managersRolesID) && memberInfo?.type === 1
      "
    >
      KYC 驗證
    </NestedNavItem>
    <NestedNavItem :href="`/members/${uid}/account`"> 銀行帳戶 </NestedNavItem>
    <NestedNavItem :href="`/members/${uid}/wallet`">
      提/入幣紀錄
    </NestedNavItem>
    <!-- <NestedNavItem :href="`/members/${uid}/money`"> 出/入金紀錄 </NestedNavItem> -->
    <NestedNavItem :href="`/members/${uid}/trade-log`">
      交易紀錄
    </NestedNavItem>
    <NestedNavItem :href="`/members/${uid}/login-log`">
      登入紀錄
    </NestedNavItem>
    <NestedNavItem :href="`/members/${uid}/friends`"> 邀請獎勵 </NestedNavItem>
  </div>
  <NuxtPage />
</template>
