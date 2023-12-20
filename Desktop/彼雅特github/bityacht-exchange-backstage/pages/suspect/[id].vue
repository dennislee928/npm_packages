<script setup>
import { WarningFilled } from "@ant-design/icons-vue";
import useSuspiciousReviewResult from "@/hooks/useSuspiciousReviewResult";
import useSuspiciousType from "@/hooks/useSuspiciousType";
import useUserOptions from "@/hooks/useUserOptions";
import useAPI from "@/hooks/useAPI";
const { result: countryList, loading: countryLoading } =
  useUserOptions("countries");
const { result: industrialList, loading: industrialLoading } =
  useUserOptions("ics");
const route = useRoute();
const { getSuspiciousTransactionDetail } = useAPI();
const { id } = route.params;
const { types } = useSuspiciousType();
const { status, badgeStatus } = useSuspiciousReviewResult();

const data = ref(null);
onMounted(async () => {
  getData();
});
async function getData() {
  let result = await getSuspiciousTransactionDetail(id);
  if (result.error) return;
  data.value = result;
}
</script>
<template>
  <div class="flex gap-3 m-6 items-center justify-start text-[#394B6A]">
    <NuxtLink
      to="/suspect"
      class="bg-[#394B6A] text-white py-1 px-3 rounded-full hover:opacity-80 active:opacity-60"
    >
      ❮ 返回
    </NuxtLink>
    <div class="text-2xl font-bold">可疑交易管理</div>
    <div class="font-bold">案件單號：{{ id }}</div>
    <StatusBadge
      :status="badgeStatus[data.dedicatedReviewResult]"
      v-if="data"
      variant="filled"
    >
      {{ status[data.dedicatedReviewResult] }}
    </StatusBadge>
  </div>
  <div class="m-6" v-if="data">
    <Alert>
      <WarningFilled class="text-red-500 text-xl" />
      {{ types[data.type] }}
      <template #content> </template>
    </Alert>
  </div>
  <TableContainer cols="1fr 1fr" v-if="data">
    <TableHeader>
      <div>會員資料</div>
    </TableHeader>
    <TableItemDual>
      <TableSubItem>
        <div>UID</div>
        <div>{{ data.usersID }}</div>
      </TableSubItem>
      <TableSubItem>
        <div>E-MAIL</div>
        <div>{{ data.email }}</div>
      </TableSubItem>
      <TableSubItem>
        <div>姓 名</div>
        <div>{{ data.lastName }} {{ data.firstName }}</div>
      </TableSubItem>
      <TableSubItem>
        <div>手機</div>
        <div>{{ data.phone }}</div>
      </TableSubItem>
      <TableSubItem>
        <div>身份證字號</div>
        <div>{{ data.nationalID }}</div>
      </TableSubItem>
      <TableSubItem>
        <div>國籍/雙國籍</div>
        <div v-if="!countryLoading && data.countriesCode != ''">
          <span>
            {{
              countryList.filter((x) => x.code === data.countriesCode)[0]
                ?.chinese
            }}
          </span>
          <span v-if="data.dualNationalityCode != ''"> / {{ " " }}</span>
          <span v-if="data.dualNationalityCode != ''">
            {{
              countryList.filter((x) => x.code === data.dualNationalityCode)[0]
                ?.chinese
            }}
          </span>
        </div>
        <div v-else>讀取中</div>
      </TableSubItem>
      <TableSubItem>
        <div>行業別</div>
        <div v-if="!industrialLoading">
          {{
            industrialList.filter(
              (x) => x.id === data.industrialClassificationsID
            )[0]?.chinese
          }}
        </div>
        <div v-else>讀取中</div>
      </TableSubItem>
      <TableSubItem>
        <div>年收入</div>
        <div>{{ data.annualIncome }}</div>
      </TableSubItem>
      <TableSubItem>
        <div>資金來源</div>
        <div>{{ data.fundsSources }}</div>
      </TableSubItem>
      <TableSubItem>
        <div>使用目的</div>
        <div>{{ data.purposeOfUse }}</div>
      </TableSubItem>
      <TableSubItem>
        <div>註冊時間</div>
        <div>{{ data.registerAt }}</div>
      </TableSubItem>
      <TableSubItem>
        <div>註冊 IP</div>
        <div>{{ data.registerIP }}</div>
      </TableSubItem>
    </TableItemDual>
  </TableContainer>

  <div class="flex m-6 gap-4" v-if="data">
    <NestedNavItem :href="`/suspect/${id}/transaction-info`">
      交易資訊
    </NestedNavItem>
    <NestedNavItem :href="`/suspect/${id}/info-check`">
      資訊審核
    </NestedNavItem>
    <NestedNavItem :href="`/suspect/${id}/risk-assessment`">
      風控審查
    </NestedNavItem>
    <NestedNavItem :href="`/suspect/${id}/review-specialist`">
      專責審查
    </NestedNavItem>
    <NestedNavItem :href="`/suspect/${id}/report-investigation-agency`">
      呈報調查局
    </NestedNavItem>
  </div>
  <NuxtPage :data="data" :getData="getData" v-if="data" />
  <Loader v-if="!data" />
</template>
