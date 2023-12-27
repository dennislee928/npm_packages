<template>
  <NuxtLayout name="members">
    <section class="relative bg-white 3xl:px-[300px] xl:px-[200px] lg:px-[100px] px-[30px] pb-10">
      <div class="relative z-30 border border-gray_800 rounded-xl xs:p-[55px] p-[25px]">
        <h1 class="flex items-center xl:mb-10 mb-5">
          <img src="/icon/hexagon.svg" alt="icon" class="mr-2 w-[18px]" />
          <p class="text-black text-subTitle font-bold">
            {{ $t('memberCenter.recentLogin') }}
          </p>
        </h1>
        <template v-if="tableData.length === 0">
          <div class="flex justify-center items-center flex-col">
            <img src="/assets/img/noData.png" />
            <p class="text-[26px] text-gray_500 mt-5">{{ $t('table.noData') }}</p>
          </div>
        </template>
        <template v-else>
          <div v-for="item of tableData" class="flex md:flex-row flex-col md:mb-4 mb-8">
            <span class="flex justify-between md:mr-5 mr-0 text-subText font-bold md:w-[20%] w-auto md:mb-0 mb-2"
              ><span class="md:hidden block">{{ $t('memberCenter.loginTime') }}</span
              >{{ item.createdAt }}</span
            >
            <span class="flex justify-between md:mr-5 mr-0 text-small md:w-[20%] w-auto md:mb-0 mb-2"
              ><span class="md:hidden block">{{ $t('memberCenter.ipAddress') }}</span
              >{{ item.ip }}</span
            >
            <span class="flex justify-between md:mr-5 mr-0 text-small md:w-[20%] w-auto md:mb-0 mb-2"
              ><span class="md:hidden block">{{ $t('memberCenter.loginLocation') }}</span
              >{{ item.location }}</span
            >
            <span class="flex justify-between md:mr-5 mr-0 text-small md:w-[20%] w-auto md:mb-0 mb-2"
              ><span class="md:hidden block">{{ $t('memberCenter.browser') }}</span
              >{{ item.browser }}</span
            >
            <span class="flex justify-between md:mr-5 mr-0 text-small md:w-[20%] w-auto md:mb-0 mb-2"
              ><span class="md:hidden block">{{ $t('memberCenter.loginDevice') }}</span
              >{{ item.device }}</span
            >
          </div>
          <a-divider></a-divider>
          <div class="text-center mt-5">
            <a-pagination v-model:current="tableConfig.page" v-model:pageSize="tableConfig.pageSize" :total="tableConfig.total" :showSizeChanger="true" :pageSizeOptions="['25', '50', '100']" @change="pageChange" />
          </div>
        </template>
      </div>
    </section>
  </NuxtLayout>
</template>
<script setup>
import useUserStore from '@/stores/user';

const { t } = useI18n();
const tableConfig = ref({
  page: 1,
  pageSize: 25,
  total: 0,
});

const userStore = useUserStore();
const tableData = ref([]);
const getLoginRecordList = async () => {
  const result = await userStore.getLoginRecord(tableConfig.value.page, tableConfig.value.pageSize, t);
  tableData.value = result.data.value.data;
  tableConfig.value.total = result.data.value.totalRecord;
};
const pageChange = async (page, pageSize) => {
  const result = await userStore.getLoginRecord(page, pageSize);
  tableData.value = result.data.value.data;
  tableConfig.value.total = result.data.value.totalRecord;
};

onMounted(() => {
  getLoginRecordList();
});
</script>
