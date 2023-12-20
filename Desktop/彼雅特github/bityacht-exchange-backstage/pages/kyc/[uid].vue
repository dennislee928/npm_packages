<script setup>
import { twMerge } from "tailwind-merge";
import useAPI from "@/hooks/useAPI";
import useUserOptions from "@/hooks/useUserOptions";
const { getUserInfo } = useAPI();
const route = useRoute();
const { uid } = route.params;
const { result: countryList, loading: countryLoading } =
  useUserOptions("countries");

const memberInfo = ref(null);

onMounted(async () => {
  memberInfo.value = await getUserInfo(uid);
});
</script>
<template>
  <div
    class="flex gap-2 py-2 px-4 items-center justify-center text-[#394B6A] mt-4"
  >
    <div class="flex items-center gap-3 justify-center">
      <NuxtLink
        to="/kyc"
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
  </div>
  <KycPage />
</template>
