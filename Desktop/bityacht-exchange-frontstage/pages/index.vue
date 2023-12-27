<template>
  <!-- <section class="md:bg-[url('assets/img/banner.png')] bg-[url('assets/img/m-banner.png')] bg-cover bg-bottom 2xl:h-[600px] xl:h-[500px] lg:h-[450px] md:h-[400px] h-[600px] 2xl:px-[300px] xl:px-[200px] lg:px-[100px] px-[50px]">
    <div class="2xl:pt-[135px] xl:pt-[100px] md:pt-[70px] pt-[120px]">
      <h1 class="2xl:text-h2 text-h4 text-white">累積你的財富</h1>
      <h1 class="2xl:text-h1 text-h3 text-white">就在 <span class="text-waterBlue">BitYacht</span></h1>
      <p class="text-subTitle text-white mt-2 font-light">Bityacht wasn’t built to make you rich,</p>
      <p class="text-subTitle text-white font-light">it was built to make you free.</p>
      <Button variant="waterBlue" class="mt-10">{{ $t('index.signUp') }}</Button>
    </div>
  </section> -->
  <section>
    <div v-if="bannerReady" class="relative">
      <a-carousel ref="carouselRef" effect="fade" autoplay>
        <div
          v-for="item of bannerList"
          :key="item.webImage"
          :class="`relative md:-mt-[40px] mt-[72px]`"
        >
          <img :src="isMobile ? item.appImage : item.webImage" class="w-full" />
          <div
            class="absolute 3xl:top-[30%] md:top-[25%] top-[20%] 3xl:left-[18%] md:left-[12%] left-[8%]"
          >
            <p
              class="3xl:w-[550px] xl:w-[450px] md:w-[400px] w-[90%] 2xl:text-h2 lg:text-h5 md:text-[26px] text-[24px] text-white"
            >
              {{ item.title }}
            </p>
            <p
              class="3xl:w-[300px] md:w-[250px] w-[85%] md:text-subTitle text-subText text-white mt-2 font-light"
            >
              {{ item.subTitle }}
            </p>
            <NuxtLink :to="item.buttonUrl" target="_blank">
              <Button variant="waterBlue" class="mt-6">{{
                item.buttonText
              }}</Button>
            </NuxtLink>
          </div>
        </div>
      </a-carousel>
    </div>
  </section>
  <section class="2xl:px-[300px] xl:px-[200px] lg:px-[100px] px-[50px]">
    <div class="md:flex hidden justify-center items-center mt-[60px]">
      <span class="block w-[40px] border border-waterBlue mr-5"></span>
      <span class="text-h2 text-white mr-3">{{ $t("index.whyChoice") }} </span>
      <span class="text-h2 text-waterBlue">{{ $t("index.bityacht") }}</span>
      <span class="block w-[40px] border border-waterBlue ml-5"></span>
    </div>
    <div class="relative md:hidden flex flex-col items-center mt-[42px]">
      <span class="text-h5 text-white">{{ $t("index.whyChoice") }} </span>
      <div class="flex items-center">
        <span class="block w-[40px] border-b-2 border-waterBlue mr-5"></span>
        <span class="text-h5 text-waterBlue">{{ $t("index.bityacht") }}</span>
        <span class="block w-[40px] border-b-2 border-waterBlue ml-5"></span>
      </div>
    </div>
    <div class="flex justify-center flex-wrap md:mt-[62px] mt-[30px]">
      <div
        v-for="item of whyChoice"
        class="3xl:w-[315px] md:w-[48%] w-[95%] border-2 border-waterBlue rounded-xl p-2 mr-2 mb-2"
      >
        <div
          class="flex flex-col items-center border border-waterBlue rounded-xl p-3"
        >
          <div
            class="2xl:w-[200px] xl:w-[150px] w-[120px] xl:h-[150px] h-[110px] xl:mb-6 mb-2"
          >
            <img :src="getAssetsFile(item.img)" class="w-full" />
          </div>
          <p class="text-waterBlue text-h5">{{ $t(item.title) }}</p>
          <p class="text-white text-subText">{{ $t(item.subTitle) }}</p>
          <p class="text-white text-subText">{{ $t(item.subTitle2) }}</p>
        </div>
      </div>
    </div>
    <!-- <div class="relative mt-[46px]">
      <a-carousel ref="carouselRef" effect="fade" arrows autoplay>
        <div v-for="item of carouselList">
          <NuxtLink :to="item.href" target="_blank">
            <img :src="item.img" />
          </NuxtLink>
        </div>
      </a-carousel>
      <LeftCircleOutlined class="absolute left-6 top-[50%] text-gray_800 hover:text-white 2xl:scale-[3] lg:scale-[2] scale-[1.5] cursor-pointer transition-all duration-300" @click="goLeft"></LeftCircleOutlined>
      <RightCircleOutlined class="absolute right-6 top-[50%] text-gray_800 hover:text-white 2xl:scale-[3] lg:scale-[2] scale-[1.5] cursor-pointer transition-all duration-300" @click="goRight"></RightCircleOutlined>
    </div> -->
  </section>
  <section class="2xl:px-[300px] xl:px-[200px] lg:px-[100px] px-[50px]">
    <div
      class="flex items-center mt-[45px] mb-[75px] md:justify-start justify-center"
    >
      <span class="block h-[44px] border border-waterBlue"></span>
      <h1 class="md:text-h2 text-h5 text-white pl-6">
        {{ $t("index.marketTrend") }}
      </h1>
    </div>
    <div>
      <div class="flex items-center justify-center">
        <p
          class="text-white md:text-h3 text-[28px] w-[50%] text-center cursor-pointer"
          :class="`${
            isBuy
              ? `border-t border-l border-r border-waterBlue rounded-tr-lg`
              : `border-b border-waterBlue`
          }`"
          @click="marketAction('buy')"
        >
          {{ $t("index.buyCurrency") }}
        </p>
        <p
          class="text-white md:text-h3 text-[28px] w-[50%] text-center cursor-pointer"
          :class="`${
            isBuy
              ? `border-b border-waterBlue`
              : `border-t border-l border-r border-waterBlue rounded-tl-lg`
          }`"
          @click="marketAction('sale')"
        >
          {{ $t("index.saleCurrency") }}
        </p>
      </div>
      <div class="border-l border-r border-b border-waterBlue">
        <template v-if="!isMobile">
          <Table
            :columns="columns"
            :data="tableData"
            :pagination="false"
            class="indexTable"
          >
            <template v-slot:symbol="{ record }">
              <template v-if="record.symbol === 'BTC'">
                <div class="flex items-center">
                  <img src="/assets/img/BTC.svg" class="w-10 h-10 mr-3" />
                  <div class="flex flex-col items-start">
                    <p class="font-bold text-subTitle text-white">
                      {{ record.symbol }}
                    </p>
                    <p class="text-subText text-grayText">{{ record.name }}</p>
                  </div>
                </div>
              </template>
              <template v-if="record.symbol === 'ETH'">
                <div class="flex items-center">
                  <img src="/assets/img/ETH.svg" class="w-10 h-10 mr-3" />
                  <div class="flex flex-col items-start">
                    <p class="font-bold text-subTitle text-white">
                      {{ record.symbol }}
                    </p>
                    <p class="text-subText text-grayText">{{ record.name }}</p>
                  </div>
                </div>
              </template>
              <template v-if="record.symbol === 'USDT'">
                <div class="flex items-center">
                  <img src="/assets/img/USDT.svg" class="w-10 h-10 mr-3" />
                  <div class="flex flex-col items-start">
                    <p class="font-bold text-subTitle text-white">
                      {{ record.symbol }}
                    </p>
                    <p class="text-subText text-grayText">{{ record.name }}</p>
                  </div>
                </div>
              </template>
              <template v-if="record.symbol === 'USDC'">
                <div class="flex items-center">
                  <img src="/assets/img/USDC.svg" class="w-10 h-10 mr-3" />
                  <div class="flex flex-col items-start">
                    <p class="font-bold text-subTitle text-white">
                      {{ record.symbol }}
                    </p>
                    <p class="text-subText text-grayText">{{ record.name }}</p>
                  </div>
                </div>
              </template>
            </template>
            <template v-slot:updown="{ text }">
              <span v-if="Number(text) > 0" class="text-orange"
                >{{ text }}%</span
              >
              <span v-else class="text-green">{{ text }}%</span>
            </template>
            <template v-slot:trend="{ record }">
              <Echarts
                :id="record.symbol"
                :series="record.trend"
                :max="record.max"
                :min="record.min"
                location="index"
                :show="false"
                symbol="none"
                :inverse="false"
              />
              <!-- <p>{{ record }}</p> -->
            </template>
          </Table>
        </template>
        <template v-else>
          <div
            class="flex p-5"
            v-for="(item, index) of tableData"
            :key="`${item.symbol} + ${index}`"
          >
            <div class="w-[50%] flex items-center">
              <img
                :src="getAssetsFile(`${item.symbol}.svg`)"
                class="xxs:w-[50px] w-[30px] xxs:h-[50px] h-[30px]"
              />
              <div class="flex flex-col items-start xxs:ml-4 ml-1">
                <p class="font-bold text-subTitle text-white">
                  {{ item.symbol }}
                </p>
                <p class="text-subText text-grayText">{{ item.name }}</p>
              </div>
            </div>
            <div class="w-[50%] flex flex-col items-end">
              <p
                v-if="isBuy"
                class="font-bold xxs:text-subTitle text-subText text-white"
              >
                ≈ $ {{ formatValue(item.buyPrice) }}
              </p>
              <p
                v-else
                class="font-bold xxs:text-subTitle text-subText text-white"
              >
                ≈ $ {{ formatValue(item.sellPrice) }}
              </p>
              <p
                v-if="Number(item.upsAndDowns) > 0"
                class="xxs:text-subText text-small text-orange"
              >
                {{ item.upsAndDowns }} %
              </p>
              <p v-else class="xxs:text-subText text-small text-green">
                {{ item.upsAndDowns }} %
              </p>
            </div>
          </div>
        </template>
      </div>
    </div>
  </section>
  <section
    class="relative 2xl:px-[300px] xl:px-[200px] lg:px-[100px] px-[50px]"
  >
    <div
      class="flex items-center mt-[85px] mb-[75px] md:justify-start justify-center"
    >
      <span class="block h-[44px] border border-waterBlue"></span>
      <h1 class="md:text-h2 text-h5 text-white pl-6">
        {{ isMobile ? $t("index.signUpIsEasy") : $t("index.whySignUpEasy") }}
      </h1>
    </div>
    <div class="flex items-center justify-between">
      <div class="w-[48%] md:block hidden">
        <img src="assets/img/indexEarth.png" />
      </div>
      <div class="md:w-[50%] w-full">
        <div v-for="item of signUpList" class="flex mb-4 ml-10">
          <img
            :src="getAssetsFile(item.img)"
            class="xxs:h-[130px] w-[100px]"
            style="margin-top: 20px"
          />
          <div class="flex flex-col md:ml-10 ml-4 mt-5">
            <div class="flex items-end">
              <p
                class="text-waterBlue xl:text-h2 lg:text-h4 md:text-[28px] text-[24px]"
              >
                {{ $t(`${item.title}`) }}
              </p>
              <span
                class="block text-white xl:text-h5 md:text-[25px] text-[20px] xl:ml-[30px] lg:ml-[20px] ml-[10px] align-bottom"
                >{{ item.type }}</span
              >
            </div>
            <div class="text-white xl:text-subTitle md:text-normal text-small">
              {{ $t(`${item.subTitle}`) }}
            </div>
            <div
              v-if="item.subTitle2"
              class="text-white xl:text-subTitle md:text-normal text-small"
            >
              {{ $t(`${item.subTitle2}`) }}
            </div>
          </div>
        </div>
      </div>
    </div>
    <img
      src="assets/img/Group16.png"
      class="absolute top-0 right-0 md:w-[700px] w-[350px]"
    />
    <img
      src="assets/img/Group37.png"
      class="absolute left-0 top-90 md:w-[550px] w-0"
    />
  </section>
  <section>
    <div class="relative bg-grayBg md:mt-[120px] mt-[60px] z-20">
      <div
        class="2xl:px-[300px] xl:px-[200px] lg:px-[100px] px-[50px] py-10 flex justify-center"
      >
        <div class="md:flex hidden">
          <img src="assets/img/takeHands.png" class="w-[50%]" />
          <div class="w-[50%] flex flex-col items-center justify-center">
            <p class="3xl:text-h1 lg:text-h3 text-h5 text-primary">
              {{ $t("index.waitingForWhat") }}
            </p>
            <NuxtLink :to="localePath('/signUp')">
              <Button
                variant="waterBlue"
                class="xl:w-[290px] w-[200px] xl:text-h4 text-[24px] rounded-4xl mt-8"
                >{{ $t("index.startSignUp") }}</Button
              >
            </NuxtLink>
          </div>
        </div>
        <div class="md:hidden flex flex-col items-center">
          <p class="xxs:text-h5 text-[26px] font-bold text-primary">
            {{ $t("index.waitingForWhat") }}
          </p>
          <img src="assets/img/takeHands.png" class="w-full mb-10" />
          <NuxtLink :to="localePath('/signUp')" class="flex justify-center">
            <Button
              variant="waterBlue"
              class="w-[200px] text-[16px] rounded-4xl absolute -bottom-[20px]"
              >{{ $t("index.startSignUp") }}</Button
            >
          </NuxtLink>
        </div>
      </div>
    </div>
  </section>
</template>

<script setup>
import useUserStore from "@/stores/user";
import { formatValue, getAssetsFile } from "@/config/config";

const localePath = useLocalePath();
const { t } = useI18n();
useHead({
  title: t("title.index"),
  meta: [{ name: "description", content: "" }],
});
const userStore = useUserStore();
const apiBase = useRuntimeConfig().public.apiBase;
const tableData = ref([]);
const bannerList = ref([]);
const bannerReady = ref(false);
const getData = async () => {
  const result = await userStore.getSpotTrend();
  result.data.value.forEach((item, index) => {
    item.trend = item.trend.map((str) => parseFloat(str));
    item.upsAndDowns = parseFloat(item.upsAndDowns).toFixed(2);
    item.max = Math.max(...item.trend);
    item.min = Math.min(...item.trend);
  });
  tableData.value = result.data.value;
};
const getBanner = async () => {
  const result = await userStore.getBanners();
  result.data.value.forEach((item) => {
    if (item.webImage !== "" || item.appImage !== "") {
      item.webImage = `${apiBase}/assets/banners/${item.webImage}`;
      item.appImage = `${apiBase}/assets/banners/${item.appImage}`;
    }
  });
  bannerList.value = result.data.value;
  bannerReady.value = true;
};

const whyChoice = [
  {
    img: "indexWhy01.png",
    title: "index.whyChoiceTitle1",
    subTitle: "index.whyChoiceText1",
    subTitle2: "index.whyChoiceText2",
  },
  {
    img: "indexWhy02.png",
    title: "index.whyChoiceTitle2",
    subTitle: "index.whyChoiceText3",
    subTitle2: "index.whyChoiceText4",
  },
  {
    img: "indexWhy03.png",
    title: "index.whyChoiceTitle3",
    subTitle: "index.whyChoiceText5",
    subTitle2: "index.whyChoiceText6",
  },
  {
    img: "indexWhy04.png",
    title: "index.whyChoiceTitle4",
    subTitle: "index.whyChoiceText7",
    subTitle2: "index.whyChoiceText8",
  },
];
// const carouselList = [
//   {
//     img: 'assets/img/indexWhy02.png',
//     mobileImg: 'assets/img/indexWhy01.png',
//     href: 'https://google.com.tw',
//   },
//   {
//     img: 'assets/img/carousel1.png',
//     mobileImg: 'assets/img/carousel1.png',
//     href: 'https://google.com.tw',
//   },

//   {
//     img: 'assets/img/indexWhy03.png',
//     mobileImg: 'assets/img/indexWhy02.png',
//     href: 'https://google.com.tw',
//   },
// ];
const signUpList = [
  {
    img: "Group237.png",
    title: "index.Email",
    subTitle: "index.EmailSubWord",
    type: "OTP",
  },
  {
    img: "Group238.png",
    title: "index.IDVerify",
    subTitle: "index.IDVerifySubWord",
    subTitle2: "index.IDVerifySubWord2",
    type: "KYC",
  },
  {
    img: "Group236.png",
    title: "index.enterMoneyTrade",
    subTitle: "index.enterMoneyTradeSubWord",
    type: "GO",
  },
];
const isBuy = ref(true);
const marketAction = (type) => {
  if (type === "buy") {
    isBuy.value = true;
    columns.value[1] = {
      title: "table.buyMoney",
      dataIndex: "buyPrice",
      key: "buyPrice",
      type: "number",
      width: 200,
    };
    getData();
  } else {
    isBuy.value = false;
    columns.value[1] = {
      title: "table.saleMoney",
      dataIndex: "sellPrice",
      key: "sellPrice",
      type: "number",
      width: 200,
    };
    getData();
  }
};
const columns = ref([
  {
    title: "table.coinType",
    dataIndex: "symbol",
    key: "symbol",
    width: 200,
  },
  {
    title: "table.buyMoney",
    dataIndex: "buyPrice",
    key: "buyPrice",
    type: "number",
    width: 200,
  },
  {
    title: "table.upsAndDown",
    dataIndex: "upsAndDowns",
    key: "upsAndDowns",
    width: 200,
  },
  {
    title: "table.priceTrend",
    dataIndex: "trend",
    key: "trend",
  },
]);
// const carouselRef = ref();
// const goLeft = () => {
//   carouselRef.value.prev();
// };
// const goRight = () => {
//   carouselRef.value.next();
// };
const isMobile = ref(false);
let timeOut = ref(null);
let timer = ref(null);
onMounted(() => {
  if (window.innerWidth < 768) {
    isMobile.value = true;
  }
  timeOut = setTimeout(() => {
    getData();
    getBanner();
  }, 300);
  timer = setInterval(() => {
    getData();
  }, 5000);
});
onUnmounted(() => {
  clearTimeout(timeOut);
  clearTimeout(timer);
});
</script>
<style>
.indexTable .ant-table table {
  background: #060a2f !important;
}
.indexTable .ant-table-thead > tr > th {
  padding: 10px 12px !important;
  background: #060a2f !important;
  text-align: center !important;
  color: #fff !important;
  border: none !important;
}
.indexTable .ant-table-tbody > tr > td {
  padding: 8px 12px !important;
  color: #181b22 !important;
  font-weight: 400 !important;
  text-align: center !important;
  color: #fff !important;
  border: none !important;
}
/* .ant-table-tbody > tr > td:last-child {
  display: flex;
  justify-content: center;
} */
.indexTable .ant-table-tbody .textRight {
  text-align: right !important;
}
.indexTable .ant-table-tbody > tr > td.ant-table-cell-row-hover {
  background-color: #f0f0f01d !important;
}
.indexTable
  :where(.css-dev-only-do-not-override-eq3tly).ant-table-wrapper
  .ant-table-tbody
  > tr.ant-table-row:hover
  > td,
.indexTable
  :where(.css-dev-only-do-not-override-eq3tly).ant-table-wrapper
  .ant-table-tbody
  > tr
  > td.ant-table-cell-row-hover {
  background: #f0f0f01d !important;
}
.indexTable .ant-empty-description {
  color: #b0b0b0;
}
</style>
