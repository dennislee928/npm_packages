<template>
  <NuxtLayout name="members">
    <section class="relative bg-white 3xl:px-[300px] xl:px-[200px] lg:px-[100px] px-[30px] pb-10">
      <div class="relative z-30 border border-gray_800 rounded-xl 3xl:px-[120px] xl:px-[60px] px-[30px] 3xl:py-[48px] py-[25px]">
        <div class="flex md:justify-between justify-start md:flex-nowrap flex-wrap mb-7">
          <div class="flex md:w-[25%] w-[50%] flex-col items-center relative after:content-[''] xl:after:h-[100px] after:h-[60px] after:border after:border-waterBlue after:absolute after:right-0 after:top-[50%] after:-translate-y-1/2">
            <p class="xl:text-h5 md:text-[26px] text-gray_500">{{ $t('memberCenter.invitePeopleTotal') }}</p>
            <div>
              <span class="2xl:text-h2 xl:text-h4 md:text-[30px] text-[36px] text-waterBlue">{{ inviteData.totalInvited }}</span>
              <span class="text-[24px] text-gray_500 ml-1">{{ $t('memberCenter.people') }}</span>
            </div>
          </div>
          <div class="flex md:w-[25%] w-[50%] flex-col items-center relative after:content-[''] md:after:inline-block after:hidden xl:after:h-[100px] md:after:h-[60px] after:border after:border-waterBlue after:absolute after:right-0 after:top-[50%] after:-translate-y-1/2">
            <p class="xl:text-h5 md:text-[26px] text-gray_500">{{ $t('memberCenter.successPeopleNumber') }}</p>
            <div>
              <span class="2xl:text-h2 xl:text-h4 md:text-[30px] text-[36px] text-waterBlue">{{ inviteData.totalSucceed }}</span>
              <span class="text-[24px] text-gray_500 ml-1">{{ $t('memberCenter.people') }}</span>
            </div>
          </div>
          <div class="flex flex-col items-center md:w-[25%] w-full relative after:content-[''] xl:after:h-[100px] md:after:h-[60px] after:h-[1px] md:after:w-[0px] after:w-[180px] md:after:border after:border-b-2 after:border-waterBlue md:after:absolute after:relative after:right-0 md:after:top-[50%] after:-translate-y-1/2 md:mt-0 mt-5">
            <p class="xl:text-h5 md:text-[26px] text-gray_500">{{ $t('memberCenter.totalBonus') }}</p>
            <div class="md:mb-0 mb-5">
              <span class="2xl:text-h2 xl:text-h4 md:text-[30px] text-[36px] text-waterBlue">{{ commissionsData.totalReward }}</span>
              <span class="text-[24px] text-gray_500 ml-1">USDT</span>
            </div>
          </div>
          <div class="flex flex-col items-center md:w-[25%] w-full md:mt-0 mt-5">
            <p class="xl:text-h5 md:text-[26px] text-gray_500">{{ $t('memberCenter.notYetGetMoney') }}</p>
            <div>
              <span class="2xl:text-h2 xl:text-h4 md:text-[30px] text-[36px] text-waterBlue">{{ commissionsData.notWithdrew }}</span>
              <span class="text-[24px] text-gray_500 ml-1">USDT</span>
            </div>
            <a-tooltip>
              <template v-if="commissionsData.notWithdrew < 10" #title>{{ $t('memberCenter.mustBeNeedTen') }}</template>
              <a-button v-if="userData.status !== 3" class="text-subText !px-4 rounded-3xl" :disabled="commissionsData.notWithdrew < 10" @click="dialogOpen = !dialogOpen">{{ $t('memberCenter.gotoGetMoney') }}</a-button>
            </a-tooltip>
          </div>
        </div>
        <div class="flex md:flex-row flex-col justify-between">
          <div class="md:w-[28%] w-full bg-skyBlue2 py-5 3xl:px-[80px] xl:px-[40px] px-[20px] rounded-xl md:mb-0 mb-5">
            <div class="flex flex-col items-center">
              <p class="text-black text-normal font-bold">{{ $t('memberCenter.myInviteCode') }}</p>
              <div class="text-yellow xs:text-subTitle text-small flex items-center mt-1">
                <span>{{ inviteData.inviteCode }}</span>
                <span class="cursor-pointer ml-2 text-gray_500 flex" @click="copyText(inviteData.inviteCode, t)"><CopyOutlined class="text-normal" /></span>
              </div>
            </div>
          </div>
          <div class="md:w-[68%] w-full bg-skyBlue2 py-5 3xl:px-[80px] xl:px-[40px] px-[20px] rounded-xl">
            <div class="flex flex-col items-center">
              <p class="text-black text-normal font-bold">{{ $t('memberCenter.inviteHref') }}</p>
              <div class="text-yellow xs:text-subTitle text-small flex items-center mt-1">
                <a-tooltip>
                  <template #title>{{ inviteData.inviteHref }}</template>
                  <span>{{ shortHref(inviteData.inviteHref) }}</span>
                </a-tooltip>
                <span class="cursor-pointer ml-2 text-gray_500 flex" @click="copyText(inviteData.inviteHref, t)"><CopyOutlined class="text-normal" /></span>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="relative z-30 border border-gray_800 rounded-xl xs:p-[55px] p-[25px] mt-10">
        <h1 class="flex items-center xl:mb-10 mb-5">
          <img src="/icon/hexagon.svg" alt="icon" class="mr-2 w-[18px]" />
          <p class="text-black text-subTitle font-bold">
            {{ $t('memberCenter.commissionsRecord') }}
          </p>
        </h1>
        <template v-if="tableData.length === 0">
          <div class="flex justify-center items-center flex-col">
            <img src="/assets/img/noData.png" />
            <p class="text-[26px] text-gray_500 mt-5">{{ $t('table.noData') }}</p>
          </div>
        </template>
        <template v-else>
          <div v-for="item of tableData" class="flex md:flex-row flex-col md:mb-4 mb-8 justify-between">
            <span class="flex justify-between md:mr-5 mr-0 text-small font-bold md:w-[20%] w-auto md:mb-0 mb-2"
              ><span class="md:hidden block">{{ $t('memberCenter.commissionsTime') }}</span
              >{{ item.createdAt }}</span
            >
            <span v-if="item.action === 1" class="flex justify-between md:mr-5 mr-0 text-small md:w-[20%] w-auto md:mb-0 mb-2"
              ><span class="md:hidden block">{{ $t('memberCenter.commissionsAction') }}</span
              >{{ $t('memberCenter.backMoney') }}</span
            >
            <span v-else class="flex justify-between md:mr-5 mr-0 text-small text-red md:w-[20%] w-auto md:mb-0 mb-2"
              ><span class="md:hidden block text-black">{{ $t('memberCenter.commissionsAction') }}</span
              >{{ $t('memberCenter.getMoney') }}</span
            >
            <span class="flex md:mr-5 mr-0 md:justify-start justify-between text-small md:w-[20%] w-auto md:mb-0 mb-2"
              ><span class="md:hidden block">{{ $t('memberCenter.commissionsAmount') }}</span> <span> {{ item.amount }} <span class="ml-2 text-gray_500">USDT</span> </span></span
            >
            <span v-if="item.account === ''" class="flex justify-between md:mr-5 mr-0 text-small md:w-[20%] w-auto md:mb-0 mb-2"
              ><span class="md:hidden block">{{ $t('memberCenter.commissionsAccount') }}</span
              >-</span
            >
            <span v-else class="flex justify-between md:mr-5 mr-0 text-small md:w-[20%] w-auto md:mb-0 mb-2"
              ><span class="md:hidden block">{{ $t('memberCenter.commissionsAccount') }}</span
              >{{ item.account }}</span
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
  <Dialog v-model="dialogOpen">
    <div class="flex flex-col mt-[60px] md:w-[300px] w-[240px]">
      <p class="md:text-[22px] text-subTitle text-center text-gray_800 font-bold">{{ $t('memberCenter.confirmGetMoney') }}</p>
      <div class="flex justify-center mt-5">
        <img src="/assets/img/money.png" class="w-[170px]" />
      </div>
      <div class="bg-[#e5e5e5] p-3 rounded-3xl mt-5 font-bold text-center">{{ commissionsData.notWithdrew }} USDT</div>
      <div class="flex justify-center mt-10">
        <a-button class="w-[120px] flex items-center justify-center py-5 bg-waterBlue rounded-3xl text-white hover:!text-white" @click="withdraw">{{ $t('memberCenter.confirm') }}</a-button>
      </div>
      <close-circle-outlined class="md:text-h5 text-[22px] text-gray_500 absolute top-5 right-5" @click="dialogOpen = !dialogOpen" />
    </div>
  </Dialog>
</template>
<script setup>
import useUserStore from '@/stores/user';
import { copyText, shortHref } from '@/config/config';

const userStore = useUserStore();
const { t } = useI18n();
const userData = ref();
const userInfo = JSON.parse(localStorage.getItem('userInfo'));
userData.value = userInfo;
const dialogOpen = ref(false);
const inviteData = ref({
  inviteCode: '',
  totalInvited: '',
  totalSucceed: '',
  inviteHref: '',
});
const commissionsData = ref({
  notWithdrew: '',
  totalReward: '',
});
const tableData = ref([]);
const tableConfig = ref({
  page: 1,
  pageSize: 25,
  total: 0,
});

const getUserInfo = async () => {
  const result = await userStore.getUserInfo(t, true);
  inviteData.value = { ...result.data.value };
};
const getInviteHref = async() => {
  const href = window.location.href;
  const firstHref = href.split('Members');
  console.log('inviteData.value :>> ', inviteData.value);
  inviteData.value.inviteHref = firstHref[0] + `signUp?inviteCode=${inviteData.value.inviteCode}`;
};
const getCommissions = async () => {
  const result = await userStore.commissions(tableConfig.value.page, tableConfig.value.pageSize, t, true);
  tableData.value = result.data.value.data;
  commissionsData.value.notWithdrew = result.data.value.notWithdrew;
  commissionsData.value.totalReward = result.data.value.totalReward;
  tableConfig.value.total = result.data.value.totalRecord;
};
const pageChange = async (page, pageSize) => {
  const result = await userStore.commissions(page, pageSize, t, true);
  tableData.value = result.data.value.data;
  commissionsData.value.notWithdrew = result.data.value.notWithdrew;
  commissionsData.value.totalReward = result.data.value.totalReward;
  tableConfig.value.total = result.data.value.totalRecord;
};
const withdraw = async () => {
  try {
    const result = await userStore.commissionsWithdraw(t);
    if (result.status.value === 'success') {
      message.success(t('memberCenter.getMoneySuccess'));
      dialogOpen.value = false;
      await getCommissions();
    }
  } catch (e) {
    message.error('error');
  }
};

onMounted(async () => {
  await getCommissions();
  setTimeout(async () => {
    await getUserInfo();
    await getInviteHref();
  }, 100);
});
</script>
