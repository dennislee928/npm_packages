<template>
  <NuxtLayout name="members">
    <section
      class="relative bg-white 3xl:px-[300px] xl:px-[200px] lg:px-[100px] px-[30px] pb-10"
    >
      <div
        v-if="userStatus?.bankAccountStatus === 0"
        class="relative z-30 border border-gray_800 rounded-xl md:px-10 md:py-10 px-5 py-10"
      >
        <div class="flex md:flex-row flex-col justify-between">
          <div class="flex items-start">
            <img src="/icon/hexagon.svg" />
            <div class="flex flex-col ml-2">
              <p class="text-subTitle font-bold">
                {{ $t("memberCenter.bindBankAccount") }}
              </p>
              <p class="text-normal text-gray_800 mt-2">
                {{ $t("memberCenter.forTWDInOrOut") }}
              </p>
              <p class="text-small text-yellow mt-1">
                {{ $t("memberCenter.onlyUseHereToBind") }}
              </p>
            </div>
          </div>
          <div class="flex items-start md:mt-0 mt-8">
            <img src="/icon/light.svg" />
            <div class="flex flex-col ml-2">
              <p class="text-small text-yellow">
                {{ $t("memberCenter.precautions") }}
              </p>
              <p class="text-small text-gray_500 mt-2">
                {{ $t("memberCenter.precautionsText1") }}
              </p>
              <p class="text-small text-gray_500 mt-1">
                {{ $t("memberCenter.precautionsText2") }}
              </p>
            </div>
          </div>
        </div>
        <template v-if="!isMobile">
          <div class="flex md:flex-row flex-col mt-8 justify-between">
            <div class="md:w-[48%] w-full flex flex-col">
              <span class="flex items-center text-normal font-bold"
                ><img src="/icon/hexagon.svg" class="mr-2" />{{
                  $t("memberCenter.accountName")
                }}</span
              >
              <FormInput
                v-model="formState.name"
                :placeholder="
                  $t('memberCenter.pleaseEnter') +
                  $t('memberCenter.accountName')
                "
                disabled
              ></FormInput>
            </div>
            <div class="md:w-[48%] w-full flex flex-col md:mt-0 mt-8">
              <span class="flex items-center text-normal font-bold"
                ><img src="/icon/hexagon.svg" class="mr-2" />{{
                  $t("memberCenter.openAccountBank")
                }}</span
              >
              <FormAntdSelect
                v-model="formState.banksCode"
                class="mt-2"
                :placeholder="
                  $t('AuthVerify.pleaseSelect') +
                  $t('memberCenter.openAccountBank')
                "
                disabled
              >
                <a-select-option
                  v-for="item of bankData"
                  :key="item"
                  :value="item.code"
                >
                  ({{ item.code }}){{ item.chinese }}</a-select-option
                >
              </FormAntdSelect>
            </div>
          </div>
          <div class="flex md:flex-row flex-col mt-8 justify-between">
            <div class="md:w-[48%] w-full flex flex-col">
              <span class="flex items-center text-normal font-bold"
                ><img src="/icon/hexagon.svg" class="mr-2" />{{
                  $t("memberCenter.bankArea")
                }}</span
              >
              <FormAntdSelect
                v-model="formState.branchsCode"
                class="mt-2"
                :placeholder="
                  $t('AuthVerify.pleaseSelect') + $t('memberCenter.bankArea')
                "
                disabled
              >
                <a-select-option
                  v-for="item of branchList"
                  :key="item"
                  :value="item.code"
                >
                  ({{ item.code }}){{ item.chinese }}</a-select-option
                >
              </FormAntdSelect>
            </div>
            <div class="md:w-[48%] w-full flex flex-col md:mt-0 mt-8">
              <span class="flex items-center text-normal font-bold"
                ><img src="/icon/hexagon.svg" class="mr-2" />{{
                  $t("memberCenter.bankAccount")
                }}</span
              >
              <FormInput
                v-model="formState.account"
                type="number"
                :placeholder="
                  $t('memberCenter.pleaseEnter') +
                  $t('memberCenter.bankAccount')
                "
                disabled
              ></FormInput>
            </div>
          </div>
          <div
            class="bg-gray_800 text-white rounded-xl md:px-[70px] px-10 md:py-10 py-6 mt-8"
          >
            <ul class="list-disc md:text-normal text-subText">
              <li class="mb-1">{{ $t("memberCenter.bankACcountText1") }}</li>
              <li class="mb-1">{{ $t("memberCenter.bankACcountText2") }}</li>
              <li>{{ $t("memberCenter.bankACcountText3") }}</li>
            </ul>
          </div>
          <div class="flex mt-5 ml-0 justify-center">
            <a-button
              disabled
              class="w-[80px] bg-waterBlue rounded-2xl text-white hover:!text-white"
              @click="dialogOpen = !dialogOpen"
              >{{ $t("memberCenter.submit") }}</a-button
            >
          </div>
        </template>
        <div
          class="md:hidden block cursor-pointer text-subText text-waterBlue underline absolute top-[20px] sm:right-[50px] right-[25px]"
          @click="isMobile = !isMobile"
        >
          {{
            !isMobile ? $t("memberCenter.cancel") : $t("memberCenter.goVerify")
          }}
        </div>
      </div>
      <div
        v-else
        class="relative z-30 border border-gray_800 rounded-xl md:px-10 md:py-10 px-5 py-10"
      >
        <div class="flex md:flex-row flex-col justify-between">
          <div class="flex items-start relative">
            <img src="/icon/hexagon.svg" />
            <div class="flex flex-col ml-2">
              <p class="text-subTitle font-bold">
                {{ $t("memberCenter.bindBankAccount") }}
              </p>
              <p class="text-normal text-gray_800 mt-2">
                {{ $t("memberCenter.forTWDInOrOut") }}
              </p>
              <p class="text-small text-yellow mt-1">
                {{ $t("memberCenter.onlyUseHereToBind") }}
              </p>
              <div class="mt-5">
                <p>
                  {{ $t("memberCenter.bankAccount") }} :
                  {{ userStatus?.bankAccount }}
                </p>
                <p class="mt-2">
                  {{ $t("memberCenter.openAccountBank") }} : {{ bankName }}-{{
                    branchsName
                  }}
                </p>
                <p class="mt-2">
                  {{ $t("memberCenter.bindTime") }} :
                  {{ userStatus?.bankAccountCreatedAt }}
                </p>
              </div>
              <template v-if="userStatus?.bankAccountStatus === 1">
                <div class="mt-2 text-yellow flex items-center">
                  <span
                    class="inline-block w-4 h-4 rounded-full bg-yellow mr-2"
                  ></span>
                  <span>{{ $t("memberCenter.bankAccountStatus1") }}</span>
                </div>
              </template>
              <template v-else-if="userStatus?.bankAccountStatus === 2">
                <div class="mt-2 text-green flex items-center">
                  <span
                    class="inline-block w-4 h-4 rounded-full bg-green mr-2"
                  ></span>
                  <span>{{ $t("memberCenter.bankAccountStatus2") }}</span>
                </div>
              </template>
              <template v-else-if="userStatus?.bankAccountStatus === 3">
                <div class="mt-2 text-red flex items-center">
                  <span
                    class="inline-block w-4 h-4 rounded-full bg-red mr-2"
                  ></span>
                  <span>{{ $t("memberCenter.bankAccountStatus3") }}</span>
                </div>
              </template>
            </div>
            <!-- <div class="md:mt-5 mt-0 md:-ml-[70px] md:relative absolute -top-5 -right-2 flex md:flex-row flex-col">
              <Button variant="shallow_gray" class="md:w-[100px] w-[75px] md:text-normal text-subText rounded-xxl md:mb-0 mb-2" @click="userStatus.bankAccountStatus = 0">{{ $t('memberCenter.edit') }}</Button>
              <Button variant="waterBlue" class="md:w-[100px] w-[75px] md:ml-3 ml-0 md:text-normal text-subText rounded-xxl" @click="deleteCheckDialog = !deleteCheckDialog">{{ $t('memberCenter.delete') }}</Button>
            </div> -->
          </div>
        </div>
      </div>
    </section>
    <Dialog v-model="dialogOpen">
      <div class="flex flex-col mt-[60px] md:w-[340px] w-[280px]">
        <p class="md:text-[22px] text-subTitle text-center font-bold">
          {{ $t("memberCenter.pleaseConfirmInformation") }}
        </p>
        <p class="md:text-subText text-small text-center text-gray_500 mt-5">
          {{ $t("memberCenter.pleaseConfirmInformationText1") }}
        </p>
        <p class="md:text-subText text-small text-center text-gray_500">
          {{ $t("memberCenter.pleaseConfirmInformationText2") }}
        </p>
        <div class="bg-[#e5e5e5] md:p-8 p-5 rounded-lg mt-10">
          <p class="md:text-subText text-small">
            {{ $t("memberCenter.accountName") }} : {{ formState.name }}
          </p>
          <p class="mt-2 md:text-subText text-small">
            {{ $t("memberCenter.bankAccount") }} : {{ formState.account }}
          </p>
          <p class="mt-2 md:text-subText text-small">
            {{ $t("memberCenter.openAccountBank") }} : ({{
              formState.banksCode
            }}) {{ bankName }}
          </p>
          <p class="mt-2 md:text-subText text-small">
            {{ $t("memberCenter.bankArea") }} : ({{ formState.branchsCode }})
            {{ branchsName }}
          </p>
        </div>
        <div class="flex justify-center mt-10">
          <a-button
            :loading="isLoading"
            class="w-[160px] flex items-center justify-center py-5 bg-waterBlue rounded-lg text-white hover:!text-white"
            @click="submit"
            >{{ $t("memberCenter.submit") }}</a-button
          >
        </div>
        <close-circle-outlined
          class="md:text-h5 text-[22px] text-gray_500 absolute top-5 right-5"
          @click="dialogOpen = !dialogOpen"
        />
      </div>
    </Dialog>
    <Dialog v-model="deleteCheckDialog">
      <div class="flex flex-col items-center">
        <p class="lg:text-[24px] text-[18px] mt-5">
          {{ $t("memberCenter.deleteConfirm") }}
        </p>
        <div>
          <Button
            variant="gray"
            class="w-[120px] mt-10 mb-4 mr-4 lg:text-[16px] text-[14px] rounded-xl"
            @click="deleteCheckDialog = !deleteCheckDialog"
            >{{ $t("memberCenter.cancel") }}</Button
          >
          <Button
            variant="waterBlue"
            class="w-[120px] mt-10 mb-4 lg:text-[16px] text-[14px] rounded-xl"
            :loading="isLoading"
            @click="deleteBank"
            >{{ $t("signUp.continue") }}</Button
          >
        </div>
      </div>
    </Dialog>
  </NuxtLayout>
</template>

<script setup>
import useUserStore from "@/stores/user";
import { useLocalStorage } from "@vueuse/core";

const userStore = useUserStore();
const isLoading = computed(() => userStore.loadingButton);
const { t } = useI18n();
const formState = reactive({
  name: null,
  account: null,
  banksCode: null,
  branchsCode: null,
});
const dialogOpen = ref(false);
const deleteCheckDialog = ref(false);
const isMobile = ref(false);
const bankData = ref(null);
const branchList = ref(null);
const bankName = ref("");
const branchsName = ref("");
watch(
  () => formState.banksCode,
  (newVal) => {
    const list = bankData.value.filter((bank) => bank.code === newVal);
    bankName.value = list[0].chinese;
    branchList.value = list[0].branchs;
    formState.branchsCode = "";
  }
);
watch(
  () => formState.branchsCode,
  (newVal) => {
    if (newVal !== "") {
      const branchsBankName = branchList.value.filter(
        (bank) => bank.code === newVal
      );
      branchsName.value = branchsBankName[0].chinese;
    }
  }
);
const branchDisabled = computed(() => {
  return formState.banksCode === null;
});
const buttonDisabled = computed(() => {
  // if (!formState.name || !formState.account || !formState.banksCode || !formState.branchsCode) return true;
  // else return false;
  return (
    !formState.name ||
    !formState.account ||
    !formState.banksCode ||
    !formState.branchsCode
  );
});
const submit = async () => {
  userStore.loadingButton = true;
  const result = await userStore.setBank(formState, t);
  if (result.status.value === "success") {
    dialogOpen.value = false;
    message.success(t("signUp.success"));
    setTimeout(() => {
      location.reload();
    }, 1000);
  }
};
const getBank = async () => {
  const result = await userStore.getBankOptions(t, true);
  bankData.value = result.data.value.banks;
};
const deleteBank = async () => {
  const result = await userStore.deleteBank(t);
  if (result.status.value === "success") {
    deleteCheckDialog.value = false;
    message.success(t("signUp.success"));
    setTimeout(() => {
      location.reload();
    }, 1000);
  }
};
const userStatus = ref(null);
const getUserInfo = async () => {
  const result = await userStore.getUserInfo(t, true);
  userStatus.value = result.data.value;
  if (userStatus.value.bankAccountStatus > 0) {
    const bank = bankData.value.filter(
      (bank) => bank.code === userStatus.value.banksCode
    );
    bankName.value = bank[0].chinese;
    const branchsBankName = bank[0].branchs.filter(
      (branch) => branch.code === userStatus.value.branchsCode
    );
    branchsName.value = branchsBankName[0].chinese;
  }
};
onMounted(async () => {
  await getBank();
  setTimeout(() => {
    getUserInfo();
  }, 100);
  const screen = window.innerWidth;
  if (screen < 768) {
    isMobile.value = true;
  }
});
</script>
<style>
input[type="number"]::-webkit-outer-spin-button,
input[type="number"]::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}
</style>
