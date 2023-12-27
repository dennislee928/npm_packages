<template>
  <section class="relative 3xl:px-[300px] xl:px-[200px] lg:px-[100px] px-[30px] bg-white pb-[100px]">
    <template v-if="step === 1">
      <div class="relative z-30">
        <h1 class="flex items-center justify-center md:text-h2 text-h4 text-gray_800 text-center pt-[88px] md:pb-[60px] pb-7"><left-circle-outlined class="mr-4 text-h5" @click="back" />{{ $t('AuthVerify.IDVerify') }}</h1>
        <div class="border border-gray_800 rounded-lg py-[80px] px-5 md:w-[450px] w-[90%] mx-auto">
          <p class="text-black md:text-h5 text-[26px] font-bold text-center">{{ $t('AuthVerify.yourCountry') }}</p>
          <div class="flex md:flex-row flex-col justify-between md:items-start items-center mt-5">
            <div class="cursor-pointer scale-100 hover:scale-105 transition-all duration-300 md:mb-0 mb-8" @click="setCountry('local')">
              <img src="/assets/img/localCountry.png" alt="local" />
            </div>
            <div class="cursor-pointer scale-100 hover:scale-105 transition-all duration-300" @click="setCountry('foreign')">
              <img src="/assets/img/foreignCountry.png" alt="foreign" />
            </div>
          </div>
        </div>
      </div>
    </template>
    <template v-if="step === 2 && country === 'local'">
      <div class="relative z-20">
        <h1 class="md:text-h2 text-h4 text-gray_800 text-center pt-[88px] md:pb-[60px] pb-7">{{ $t('AuthVerify.verifyData') }}</h1>
        <a-divider class="md:hidden block border border-gray_300"></a-divider>
        <a-form :model="localFormState" @finish="submitCheck" @validate="validateRules" :rules="rules">
          <div class="flex flex-col">
            <a-form-item name="nationalID" class="w-full">
              <FormLabel label="AuthVerify.nationalIDText" haveIcon required></FormLabel>
              <FormInput v-model="localFormState.nationalID" :placeholder="$t('AuthVerify.pleaseEnter') + $t('AuthVerify.nationalIDText')" type="search" />
            </a-form-item>
            <a-form-item name="countriesCode" class="w-full">
              <FormLabel label="AuthVerify.countriesCode" haveIcon required></FormLabel>
              <FormAntdSelect v-model="localFormState.countriesCode" disabled class="mt-2 disabled:text-waterBlue" :placeholder="$t('AuthVerify.pleaseSelect') + $t('AuthVerify.countriesCode')">
                <a-select-option v-for="i of copyCountryList" :key="i" :value="i.code">{{ i.chinese }} / {{ i.english }}</a-select-option>
              </FormAntdSelect>
            </a-form-item>
            <a-form-item name="dualNationalityCode" class="w-full">
              <FormLabel label="AuthVerify.secondCountry" haveIcon></FormLabel>
              <FormAntdSelect v-model="localFormState.dualNationalityCode" allowClear class="mt-2 disabled:text-waterBlue" :placeholder="$t('AuthVerify.pleaseSelect') + $t('AuthVerify.secondCountry')">
                <a-select-option v-for="i of countryList" :key="i" :value="i.code">{{ i.chinese }} / {{ i.english }}</a-select-option>
              </FormAntdSelect>
            </a-form-item>
          </div>
          <a-divider class="md:hidden block border border-gray_300"></a-divider>
          <div class="flex flex-col">
            <a-form-item name="lastName" class="w-full">
              <FormLabel label="AuthVerify.lastName" haveIcon required></FormLabel>
              <FormInput v-model="localFormState.lastName" :placeholder="$t('AuthVerify.pleaseEnter') + $t('AuthVerify.lastName')" />
            </a-form-item>
            <a-form-item name="firstName" class="w-full">
              <FormLabel label="AuthVerify.firstName" haveIcon required></FormLabel>
              <FormInput v-model="localFormState.firstName" :placeholder="$t('AuthVerify.pleaseEnter') + $t('AuthVerify.firstName')" />
            </a-form-item>
            <a-form-item name="birthDate" class="w-full">
              <FormLabel label="AuthVerify.birthDate" haveIcon required></FormLabel>
              <FormDate v-model:value="localFormState.birthDate" :placeholder="$t('AuthVerify.pleaseEnter') + $t('AuthVerify.birthDate')" @change-date="changeDate" />
            </a-form-item>
            <a-form-item name="phone" class="w-full">
              <FormLabel label="AuthVerify.phone" notice="AuthVerify.notice1" haveIcon required></FormLabel>
              <FormInput v-model="localFormState.phone" :placeholder="$t('AuthVerify.pleaseEnter') + $t('AuthVerify.phone') + $t('AuthVerify.phoneExample')" />
            </a-form-item>
          </div>
          <div class="flex flex-col">
            <a-form-item name="industrialClassificationsID" class="w-full">
              <FormLabel label="AuthVerify.jobStatus" haveIcon required></FormLabel>
              <FormAntdSelect v-model="localFormState.industrialClassificationsID" class="mt-2 disabled:text-waterBlue" :placeholder="$t('AuthVerify.pleaseSelect') + $t('AuthVerify.jobStatus')">
                <a-select-option v-for="i of jobList" :key="i" :value="i.id">{{ i.chinese }} / {{ i.english }}</a-select-option>
              </FormAntdSelect>
            </a-form-item>
          </div>
          <div class="flex flex-col">
            <div class="flex justify-between w-full md:flex-nowrap flex-wrap">
              <a-form-item name="address" class="md:w-[20%] w-[48%]">
                <FormLabel label="AuthVerify.address" haveIcon required></FormLabel>
                <FormAntdSelect v-model="localFormState.city" class="mt-2 disabled:text-waterBlue" :placeholder="$t('AuthVerify.pleaseSelect') + $t('AuthVerify.city')">
                  <a-select-option v-for="i of cityData" :key="i" :value="i">{{ i }}</a-select-option>
                </FormAntdSelect>
              </a-form-item>
              <a-form-item name="area" class="md:w-[20%] w-[48%]">
                <!-- <FormLabel label=""></FormLabel> -->
                <FormAntdSelect v-model="localFormState.area" class="mt-9 disabled:text-waterBlue" :class="{ 'border border-[#ff7875]': addressError }" :placeholder="$t('AuthVerify.pleaseSelect') + $t('AuthVerify.area')">
                  <a-select-option v-for="i of areaList" :key="i" :value="i">{{ i }}</a-select-option>
                </FormAntdSelect>
              </a-form-item>
              <a-form-item name="street" class="md:w-[55%] w-full">
                <!-- <FormLabel label=""></FormLabel> -->
                <FormInput v-model="localFormState.street" :address-error="addressError" :placeholder="$t('AuthVerify.pleaseEnter') + $t('AuthVerify.street')" class="md:mt-7 -mt-2" @input="setAddress($event.target.value)" />
              </a-form-item>
            </div>
            <a-divider class="md:hidden block border border-gray_300"></a-divider>
            <a-form-item name="annualIncome" class="w-full">
              <FormLabel label="AuthVerify.annualIncome" haveIcon required></FormLabel>
              <FormAntdSelect v-model="localFormState.annualIncome" class="mt-2 disabled:text-waterBlue" :placeholder="$t('AuthVerify.pleaseSelect') + $t('AuthVerify.annualIncome')">
                <a-select-option v-for="i of yearsMoneyList" :key="i" :value="$t(`${i}`)">
                  {{ $t(`${i}`) }}
                </a-select-option>
              </FormAntdSelect>
            </a-form-item>
          </div>
          <div class="flexflex-col">
            <a-form-item name="fundsSources" class="w-full">
              <FormLabel label="AuthVerify.fundsSources" haveIcon required></FormLabel>
              <FormAntdSelect v-model="localFormState.fundsSources" class="mt-2 disabled:text-waterBlue" :placeholder="$t('AuthVerify.pleaseSelect') + $t('AuthVerify.fundsSources')">
                <a-select-option v-for="i of getMoneyList" :key="i" :value="$t(`${i}`)">{{ $t(`${i}`) }}</a-select-option>
              </FormAntdSelect>
            </a-form-item>
            <a-form-item name="purposeOfUse" class="w-full">
              <FormLabel label="AuthVerify.purposeOfUse" haveIcon required></FormLabel>
              <FormAntdSelect v-model="localFormState.purposeOfUse" class="mt-2 disabled:text-waterBlue" :placeholder="$t('AuthVerify.pleaseSelect') + $t('AuthVerify.purposeOfUse')">
                <a-select-option v-for="i of usePurposeList" :key="i" :value="$t(`${i}`)">{{ $t(`${i}`) }}</a-select-option>
              </FormAntdSelect>
            </a-form-item>
            <a-form-item name="investmentExperience" class="w-full">
              <FormLabel label="AuthVerify.investmentExperience" haveIcon required></FormLabel>
              <FormAntdSelect v-model="localFormState.investmentExperience" class="mt-2 disabled:text-waterBlue" :placeholder="$t('AuthVerify.pleaseSelect') + $t('AuthVerify.investmentExperience')">
                <a-select-option v-for="i of investExpList" :key="i" :value="$t(`${i}`)">{{ $t(`${i}`) }}</a-select-option>
              </FormAntdSelect>
            </a-form-item>
          </div>
          <div class="bg-[#f1f3f6] p-8 rounded-lg mt-8">
            <p class="text-normal font-bold">{{ $t('AuthVerify.notice') }}</p>
            <p class="text-normal mt-2">{{ $t('AuthVerify.noticeText') }}</p>
          </div>
          <div class="md:ml-10 ml-5 my-5">
            <a-checkbox v-model:checked="checkSubmit">{{ $t('AuthVerify.agree') }}</a-checkbox>
          </div>
          <a-divider></a-divider>
          <div class="flex justify-center">
            <Button variant="gray" class="w-[100px] shadow-xl shadow-gray_700/20 mr-5" @click="backStep">{{ $t('AuthVerify.back') }}</Button>
            <Button variant="waterBlue" html-type="submit" class="w-[100px] shadow-xl shadow-gray_700/20" :disabled="!checkSubmit" :class="`${checkSubmit ? `` : `cursor-not-allowed`}`" :loading="isLoading">{{ $t('AuthVerify.submit') }}</Button>
          </div>
          <div class="flex justify-center mt-8 md:flex-row flex-col items-center">
            <p class="text-normal text-waterBlue">{{ $t('AuthVerify.noticeText2') }}</p>
            <p class="text-normal text-waterBlue">{{ $t('AuthVerify.noticeText3') }}</p>
          </div>
        </a-form>
      </div>
    </template>
    <template v-if="step === 2 && country === 'foreign'">
      <div class="relative z-20">
        <h1 class="md:text-h2 text-h4 text-gray_800 text-center pt-[88px] md:pb-[60px] pb-[30px]">{{ $t('AuthVerify.verifyData') }}</h1>
        <a-divider class="md:hidden block border border-gray_300"></a-divider>
        <a-form :model="localFormState" @finish="submitCheck" @validate="validateRules" :rules="rules">
          <div class="flex flex-col">
            <a-form-item name="nationalID" class="w-full">
              <FormLabel label="AuthVerify.residenceNumber" haveIcon required></FormLabel>
              <FormInput v-model="localFormState.nationalID" :placeholder="$t('AuthVerify.pleaseEnter') + $t('AuthVerify.residenceNumber')" />
            </a-form-item>
            <a-form-item name="idImage" class="w-full">
              <FormLabel label="AuthVerify.idImage" haveIcon required></FormLabel>
              <FormUpload type="idImage" @upload="handleUpload">
                <FormInput v-model="localFormState.idImageName" :placeholder="$t('AuthVerify.pleaseUpload') + $t('AuthVerify.idImage')" />
                <FolderOpenOutlined class="text-gray_800 absolute top-1/2 right-2 -translate-y-1/2" />
              </FormUpload>
            </a-form-item>
            <a-form-item name="idBackImage" class="w-full">
              <FormLabel label="AuthVerify.idBackImage" haveIcon required></FormLabel>
              <FormUpload type="idBackImage" @upload="handleUpload">
                <FormInput v-model="localFormState.idBackImageName" :placeholder="$t('AuthVerify.pleaseUpload') + $t('AuthVerify.idBackImage')" />
                <FolderOpenOutlined class="text-gray_800 absolute top-1/2 right-2 -translate-y-1/2" />
              </FormUpload>
            </a-form-item>
          </div>
          <div class="flex flex-col">
            <a-form-item name="passportNumber" class="w-full">
              <FormLabel label="AuthVerify.passportNumber" haveIcon required></FormLabel>
              <FormInput v-model="localFormState.passportNumber" :placeholder="$t('AuthVerify.pleaseEnter') + $t('AuthVerify.passportNumber')" />
            </a-form-item>
            <a-form-item name="passportImage" class="w-full">
              <FormLabel label="AuthVerify.passportImage" haveIcon required></FormLabel>
              <FormUpload type="passportImage" @upload="handleUpload">
                <FormInput v-model="localFormState.passportImageName" :placeholder="$t('AuthVerify.pleaseUpload') + $t('AuthVerify.passportImage')" />
                <FolderOpenOutlined class="text-gray_800 absolute top-1/2 right-2 -translate-y-1/2" />
              </FormUpload>
            </a-form-item>
            <a-form-item name="idAndFaceImage" class="w-full">
              <FormLabel label="AuthVerify.idAndFaceImage" haveIcon required></FormLabel>
              <FormUpload type="idAndFaceImage" @upload="handleUpload">
                <FormInput v-model="localFormState.idAndFaceImageName" :placeholder="$t('AuthVerify.pleaseUpload') + $t('AuthVerify.idAndFaceImage')" />
                <FolderOpenOutlined class="text-gray_800 absolute top-1/2 right-2 -translate-y-1/2" />
              </FormUpload>
            </a-form-item>
          </div>
          <a-divider class="md:hidden block border border-gray_300"></a-divider>
          <div class="flex flex-col">
            <a-form-item name="lastName" class="w-full">
              <FormLabel label="AuthVerify.lastName" haveIcon required></FormLabel>
              <FormInput v-model="localFormState.lastName" :placeholder="$t('AuthVerify.pleaseEnter') + $t('AuthVerify.lastName')" />
            </a-form-item>
            <a-form-item name="firstName" class="w-full">
              <FormLabel label="AuthVerify.firstName" haveIcon required></FormLabel>
              <FormInput v-model="localFormState.firstName" :placeholder="$t('AuthVerify.pleaseEnter') + $t('AuthVerify.firstName')" />
            </a-form-item>
            <!-- <a-form-item name="gender" class="w-full">
              <FormLabel label="AuthVerify.gender" haveIcon required></FormLabel>
              <FormAntdSelect v-model="localFormState.gender" class="mt-2 disabled:text-waterBlue" :placeholder="$t('AuthVerify.pleaseSelect') + $t('AuthVerify.gender')">
                <a-select-option :value="1">{{ $t('AuthVerify.genderOption1') }}</a-select-option>
                <a-select-option :value="2">{{ $t('AuthVerify.genderOption2') }}</a-select-option>
              </FormAntdSelect>
            </a-form-item> -->
          </div>
          <div class="flex flex-col">
            <a-form-item name="birthDate" class="w-full">
              <FormLabel label="AuthVerify.birthDate" haveIcon required></FormLabel>
              <FormDate v-model:value="localFormState.birthDate" :placeholder="$t('AuthVerify.pleaseEnter') + $t('AuthVerify.birthDate')" @change-date="changeDate" />
            </a-form-item>
            <a-form-item name="countriesCode" class="w-full">
              <FormLabel label="AuthVerify.countriesCode" haveIcon required></FormLabel>
              <FormAntdSelect v-model="localFormState.countriesCode" class="mt-2 disabled:text-waterBlue" :placeholder="$t('AuthVerify.pleaseSelect') + $t('AuthVerify.countriesCode')">
                <a-select-option v-for="i of copyCountryList" :key="i" :value="i.code">{{ i.chinese }} / {{ i.english }}</a-select-option>
              </FormAntdSelect>
            </a-form-item>
            <a-form-item name="dualNationalityCode" class="w-full">
              <FormLabel label="AuthVerify.secondCountry" haveIcon></FormLabel>
              <FormAntdSelect v-model="localFormState.dualNationalityCode" allowClear class="mt-2 disabled:text-waterBlue" :placeholder="$t('AuthVerify.pleaseSelect') + $t('AuthVerify.secondCountry')">
                <a-select-option v-for="i of countryList" :key="i" :value="i.code">{{ i.chinese }} / {{ i.english }}</a-select-option>
              </FormAntdSelect>
            </a-form-item>
          </div>
          <div class="flex flex-col">
            <div class="flex justify-between w-full md:flex-nowrap flex-wrap">
              <a-form-item name="address" class="md:w-[20%] w-[48%]">
                <FormLabel label="AuthVerify.address" haveIcon required></FormLabel>
                <FormAntdSelect v-model="localFormState.city" class="mt-2 disabled:text-waterBlue" :placeholder="$t('AuthVerify.pleaseSelect') + $t('AuthVerify.city')">
                  <a-select-option v-for="i of cityData" :key="i" :value="i">{{ i }}</a-select-option>
                </FormAntdSelect>
              </a-form-item>
              <a-form-item name="area" class="md:w-[20%] w-[48%]">
                <!-- <FormLabel label=""></FormLabel> -->
                <FormAntdSelect v-model="localFormState.area" class="mt-9 disabled:text-waterBlue" :class="{ 'border border-[#ff7875]': addressError }" :placeholder="$t('AuthVerify.pleaseSelect') + $t('AuthVerify.area')">
                  <a-select-option v-for="i of areaList" :key="i" :value="i">{{ i }}</a-select-option>
                </FormAntdSelect>
              </a-form-item>
              <a-form-item name="street" class="md:w-[55%] w-full">
                <!-- <FormLabel label=""></FormLabel> -->
                <FormInput v-model="localFormState.street" :address-error="addressError" :placeholder="$t('AuthVerify.pleaseEnter') + $t('AuthVerify.street')" class="md:mt-7 -mt-2" @input="setAddress($event.target.value)" />
              </a-form-item>
            </div>
            <a-divider class="md:hidden block border border-gray_300"></a-divider>
            <a-form-item name="phone" class="w-full">
              <FormLabel label="AuthVerify.phone" notice="AuthVerify.notice1" haveIcon required></FormLabel>
              <FormInput v-model="localFormState.phone" :placeholder="$t('AuthVerify.pleaseEnter') + $t('AuthVerify.phone') + $t('AuthVerify.phoneExample')" />
            </a-form-item>
          </div>
          <div class="flex flex-col">
            <a-form-item name="industrialClassificationsID" class="w-full">
              <FormLabel label="AuthVerify.jobStatus" haveIcon required></FormLabel>
              <FormAntdSelect v-model="localFormState.industrialClassificationsID" class="mt-2 disabled:text-waterBlue" :placeholder="$t('AuthVerify.pleaseSelect') + $t('AuthVerify.jobStatus')">
                <a-select-option v-for="i of jobList" :key="i" :value="i.id">{{ i.chinese }} / {{ i.english }}</a-select-option>
              </FormAntdSelect>
            </a-form-item>
            <a-form-item name="annualIncome" class="w-full">
              <FormLabel label="AuthVerify.annualIncome" haveIcon required></FormLabel>
              <FormAntdSelect v-model="localFormState.annualIncome" class="mt-2 disabled:text-waterBlue" :placeholder="$t('AuthVerify.pleaseSelect') + $t('AuthVerify.annualIncome')">
                <a-select-option v-for="i of yearsMoneyList" :key="i" :value="$t(`${i}`)">
                  {{ $t(`${i}`) }}
                </a-select-option>
              </FormAntdSelect>
            </a-form-item>
            <a-form-item name="fundsSources" class="w-full">
              <FormLabel label="AuthVerify.fundsSources" haveIcon required></FormLabel>
              <FormAntdSelect v-model="localFormState.fundsSources" class="mt-2 disabled:text-waterBlue" :placeholder="$t('AuthVerify.pleaseSelect') + $t('AuthVerify.fundsSources')">
                <a-select-option v-for="i of getMoneyList" :key="i" :value="$t(`${i}`)">{{ $t(`${i}`) }}</a-select-option>
              </FormAntdSelect>
            </a-form-item>
          </div>
          <div class="flexflex-col">
            <a-form-item name="purposeOfUse" class="w-full">
              <FormLabel label="AuthVerify.purposeOfUse" haveIcon required></FormLabel>
              <FormAntdSelect v-model="localFormState.purposeOfUse" class="mt-2 disabled:text-waterBlue" :placeholder="$t('AuthVerify.pleaseSelect') + $t('AuthVerify.purposeOfUse')">
                <a-select-option v-for="i of usePurposeList" :key="i" :value="$t(`${i}`)">{{ $t(`${i}`) }}</a-select-option>
              </FormAntdSelect>
            </a-form-item>
            <a-form-item name="investmentExperience" class="w-full">
              <FormLabel label="AuthVerify.investmentExperience" haveIcon required></FormLabel>
              <FormAntdSelect v-model="localFormState.investmentExperience" class="mt-2 disabled:text-waterBlue" :placeholder="$t('AuthVerify.pleaseSelect') + $t('AuthVerify.investmentExperience')">
                <a-select-option v-for="i of investExpList" :key="i" :value="$t(`${i}`)">{{ $t(`${i}`) }}</a-select-option>
              </FormAntdSelect>
            </a-form-item>
            <div class="md:w-[30%] w-full"></div>
          </div>
          <div class="bg-[#f1f3f6] p-8 rounded-lg mt-8">
            <p class="text-normal font-bold">{{ $t('AuthVerify.notice') }}</p>
            <p class="text-normal mt-2">{{ $t('AuthVerify.noticeText') }}</p>
          </div>
          <div class="md:ml-10 ml-5 my-5">
            <a-checkbox v-model:checked="checkSubmit">{{ $t('AuthVerify.agree') }}</a-checkbox>
          </div>
          <a-divider></a-divider>
          <div class="flex justify-center">
            <Button variant="gray" class="w-[100px] shadow-xl shadow-gray_700/20 mr-5" @click="backStep">{{ $t('AuthVerify.back') }}</Button>
            <Button variant="waterBlue" html-type="submit" class="w-[100px] shadow-xl shadow-gray_700/20" :disabled="!checkSubmit" :class="`${checkSubmit ? `` : `cursor-not-allowed`}`" :loading="isLoading">{{ $t('AuthVerify.submit') }}</Button>
          </div>
          <div class="flex justify-center mt-8 md:flex-row flex-col items-center">
            <p class="text-normal text-waterBlue">{{ $t('AuthVerify.noticeText2') }}</p>
            <p class="text-normal text-waterBlue">{{ $t('AuthVerify.noticeText3') }}</p>
          </div>
        </a-form>
      </div>
    </template>
    <img src="/assets/img/Group16.png" class="absolute right-0 top-0 md:w-[600px] w-[400px] z-10" />
    <img src="/assets/img/Group37.png" class="absolute left-0 top-0 md:w-[800px] w-[200px] z-10" />
  </section>
  <Dialog v-model="dialogOpen">
    <div class="flex flex-col items-center">
      <p class="text-subTitle font-bold mt-[70px] mb-5">{{ $t('AuthVerify.enterPhonePassCode') }}</p>
      <p class="text-normal text-grayText">{{ $t('signUp.passCodeAlreadySend') }}</p>
      <p class="text-normal text-grayText">{{ formatPhone(localFormState.phone) }}</p>
      <!-- <Button variant="white" :disabled="buttonDisabled" class="mt-[50px]" :class="{ 'bg-gray_300 cursor-not-allowed text-gray_400 hover:bg-opacity-100 active:bg-opacity-100': buttonDisabled }" @click="reSend"
        >{{ $t('signUp.reSend') }} <span v-if="buttonDisabled">{{ countDown }} s</span></Button
      > -->
      <FormEmailValid :codes="codes" />
      <Button v-if="continueDisabled" variant="disabled" :disabled="continueDisabled" class="w-[200px] mt-10 mb-4 text-[25px] rounded-xl" @click="verifyContinue">{{ $t('signUp.continue') }}</Button>
      <Button v-else variant="waterBlue" class="w-[200px] mt-10 mb-4 text-[25px] rounded-xl" :loading="isLoading" @click="verifyContinue">{{ $t('signUp.continue') }}</Button>
    </div>
  </Dialog>

  <Dialog v-model="submitCheckDialog">
    <div class="flex flex-col items-center">
      <p class="lg:text-[26px] text-[18px] mt-5">{{ $t('AuthVerify.submitCheck') }}</p>
      <div>
        <Button variant="gray" class="w-[120px] mt-10 mb-4 mr-4 lg:text-[18px] text-[16px] rounded-xl" @click="submitCheckDialog = !submitCheckDialog">{{ $t('memberCenter.cancel') }}</Button>
        <Button variant="waterBlue" class="w-[120px] mt-10 mb-4 lg:text-[18px] text-[16px] rounded-xl" :loading="isLoading" @click="submit">{{ $t('signUp.continue') }}</Button>
      </div>
    </div>
  </Dialog>
</template>

<script setup>
import useUserStore from '@/stores/user';
import { yearsMoneyList, getMoneyList, usePurposeList, investExpList, cityData, cityAreaData } from '@/config/list.js';
import { formatPhone } from '@/config/config';

const localePath = useLocalePath();
const back = async () => {
  await navigateTo(localePath(`/Members/${userData.value.id}/level`));
};
const { t } = useI18n();
useHead({
  title: t('title.memberCenter_IDVerify'),
  meta: [{ name: 'description', content: '' }],
});
const userStore = useUserStore();
const isLoading = computed(() => userStore.loadingButton);
const step = ref(0);
const country = ref('');
watch(
  () => country.value,
  (val) => {
    if (val === 'local') {
      localFormState.countriesCode = 'TWN';
    }
  }
);
const localFormState = reactive({
  city: '',
  area: '',
  street: '',
  address: '',
  annualIncome: null,
  birthDate: '',
  countriesCode: null,
  dualNationalityCode: null,
  firstName: '',
  lastName: '',
  fundsSources: null,
  industrialClassificationsID: null,
  investmentExperience: null,
  nationalID: '',
  passportNumber: '',
  phone: '',
  phoneToken: '',
  phoneVerificationCode: '',
  purposeOfUse: null,
  idImage: '', //file
  idImageName: '',
  idBackImage: '', //file
  idBackImageName: '',
  passportImage: '', //file
  passportImageName: '',
  idAndFaceImage: '', //file
  idAndFaceImageName: '',
});
const addressError = ref(false);
const rules = {
  idImage: {
    require: true,
    validator: async (rule, value) => {
      if (value === '') {
        return Promise.reject(t('rules.dontEmpty'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'submit',
  },
  idBackImage: {
    require: true,
    validator: async (rule, value) => {
      if (value === '') {
        return Promise.reject(t('rules.dontEmpty'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'submit',
  },
  passportImage: {
    require: true,
    validator: async (rule, value) => {
      if (value === '') {
        return Promise.reject(t('rules.dontEmpty'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'submit',
  },
  passportNumber: {
    require: true,
    validator: async (rule, value) => {
      if (value === '') {
        return Promise.reject(t('rules.dontEmpty'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'submit',
  },
  nationalID: {
    require: true,
    validator: async (rule, value) => {
      if (country.value === 'local') {
        const regex = /^[A-Z][12]\d{8}$/;
        if (!regex.test(value)) {
          return Promise.reject(t('rules.dontEmpty') + t('rules.or') + t('rules.wrongFormat'));
        }
        return Promise.resolve(t('rules.actionSuccess'));
      } else {
        const regex = /^[A-Z][A-Z\d]\d{8}$/;
        if (!regex.test(value)) {
          return Promise.reject(t('rules.dontEmpty') + t('rules.or') + t('rules.wrongFormat'));
        }
        return Promise.resolve(t('rules.actionSuccess'));
      }
    },
    trigger: 'blur',
  },
  firstName: {
    require: true,
    validator: async (rule, value) => {
      if (value === '') {
        return Promise.reject(t('rules.dontEmpty'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'submit',
  },
  lastName: {
    require: true,
    validator: async (rule, value) => {
      if (value === '') {
        return Promise.reject(t('rules.dontEmpty'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'submit',
  },
  birthDate: {
    require: true,
    validator: async (rule, value) => {
      if (value === '') {
        return Promise.reject(t('rules.dontEmpty'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'submit',
  },
  // gender: {
  //   require: true,
  //   validator: async (rule, value) => {
  //     if (!value) {
  //       return Promise.reject(t('rules.dontEmpty'));
  //     }
  //     return Promise.resolve(t('rules.actionSuccess'));
  //   },
  //   trigger: 'submit',
  // },
  phone: {
    require: true,
    validator: async (rule, value) => {
      const regex = /^09\d{8}$/;
      if (!regex.test(value)) {
        return Promise.reject(t('rules.dontEmpty') + t('rules.or') + t('rules.wrongFormat'));
      }
      const result = await userStore.checkPhone({ phone: value }, t);
      if (!result) {
        return Promise.reject(t('rules.phoneError'));
      } else {
        return Promise.resolve(t('rules.actionSuccess'));
      }
    },
    trigger: 'blur',
  },
  countriesCode: {
    require: true,
    validator: async (rule, value) => {
      if (!value) {
        return Promise.reject(t('rules.dontEmpty'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'submit',
  },
  industrialClassificationsID: {
    require: true,
    validator: async (rule, value) => {
      if (!value) {
        return Promise.reject(t('rules.dontEmpty'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'submit',
  },
  address: {
    require: true,
    validator: async (rule, value) => {
      if (value === '') {
        addressError.value = true;
        return Promise.reject(t('rules.dontEmpty'));
      }
      addressError.value = false;
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'submit',
  },
  annualIncome: {
    require: true,
    validator: async (rule, value) => {
      if (!value) {
        return Promise.reject(t('rules.dontEmpty'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'submit',
  },
  fundsSources: {
    require: true,
    validator: async (rule, value) => {
      if (!value) {
        return Promise.reject(t('rules.dontEmpty'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'submit',
  },
  purposeOfUse: {
    require: true,
    validator: async (rule, value) => {
      if (!value) {
        return Promise.reject(t('rules.dontEmpty'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'submit',
  },
  investmentExperience: {
    require: true,
    validator: async (rule, value) => {
      if (!value) {
        return Promise.reject(t('rules.dontEmpty'));
      }
      return Promise.resolve(t('rules.actionSuccess'));
    },
    trigger: 'submit',
  },
};
const validateRules = (name, status, errorMsgs) => {
  if (!status) {
    message.error(t(`AuthVerify.${name}`) + errorMsgs[0]);
  }
};
const setAddress = (value) => {
  if (localFormState.city === '' || localFormState.area === '') {
    localFormState.street = '';
    return message.warn(t('rules.pleaseEnterCityFirst'));
  }
  localFormState.address = localFormState.city + localFormState.area + localFormState.street;
};

const countryList = ref([]);
const copyCountryList = ref([]);
const jobList = ref([]);
const checkSubmit = ref(false);
const dialogOpen = ref(false);
const submitCheckDialog = ref(false);
const codes = ref(['', '', '', '', '', '']);
const continueDisabled = ref(true);
watch(
  () => [codes.value],
  ([newVal]) => {
    const noEmpty = newVal.every((item) => item !== '');
    if (noEmpty) continueDisabled.value = false;
    else continueDisabled.value = true;
  },
  { deep: true }
);
const changeDate = (date) => {
  localFormState.birthDate = date;
};
const setCountry = (type) => {
  step.value++;
  country.value = type;
};
const backStep = () => {
  step.value--;
  country.value = '';
  addressError.value = false;
  checkSubmit.value = false;
  Object.keys(localFormState).forEach((key) => {
    localFormState[key] = '';
  });
};
const submitCheck = async () => {
  submitCheckDialog.value = true;
};
const submit = async () => {
  submitCheckDialog.value = false;
  userStore.loadingButton = true;
  const result = await userStore.sendPhoneVerify(localFormState, t);
  userStore.loadingButton = false;
  if (result?.code === 4011) {
    const ErrorKey = Object.keys(result.data);
    for (let i = 0; i < ErrorKey.length; i++) {
      message.error(t(`AuthVerify.${ErrorKey[i]}`) + t(`AuthVerify.errorText`));
    }
  } else if (result.status?.value === 'success') {
    message.success(t('AuthVerify.alreadySendPhonePassCode'));
    dialogOpen.value = true;
  } else {
  }
};
const verifyContinue = async () => {
  userStore.loadingButton = true;
  const data = {
    phone: localFormState.phone,
    verificationCode: '',
  };
  data.verificationCode = codes.value.join('');
  const result = await userStore.phoneVerify(data, t);
  userStore.loadingButton = false;
  if (result.status.value === 'success') {
    localFormState.phoneToken = result.data.value.phoneToken;
    const idvResult = await userStore.idv(localFormState, t);
    if (idvResult.status.value === 'success') {
      if (country.value === 'foreign') {
        message.success(t('AuthVerify.IDVerifySent'));
        await navigateTo({
          path: localePath('/status'),
          query: {
            title: 'AuthVerify.IDVerifySent',
            imgSource: 'dataSuccess',
          },
        });
      } else {
        const url = idvResult.data.value.idVerificationUrl;
        message.success(t('AuthVerify.pleaseGoContinue'));
        await navigateTo({
          path: localePath('/status'),
          query: {
            title: 'AuthVerify.pleaseGoContinue',
            url: url,
          },
        });
      }
    }
  }
};

watch(
  () => localFormState.countriesCode,
  (newVal) => {
    countryList.value = copyCountryList.value;
    countryList.value = countryList.value.filter((item) => item.code !== newVal);
    countryList.value.unshift({
      chinese: '無',
      code: '',
      english: 'none',
      locale: '',
    });
  }
);

const areaList = ref([]);
watch(
  () => localFormState.city,
  (newVal) => {
    areaList.value = cityAreaData[newVal];
  }
);

const handleUpload = (file, type, fileName) => {
  if (type === 'idImage') {
    localFormState.idImage = file;
    localFormState.idImageName = fileName;
  } else if (type === 'idBackImage') {
    localFormState.idBackImage = file;
    localFormState.idBackImageName = fileName;
  } else if (type === 'passportImage') {
    localFormState.passportImage = file;
    localFormState.passportImageName = fileName;
  } else if (type === 'idAndFaceImage') {
    localFormState.idAndFaceImage = file;
    localFormState.idAndFaceImageName = fileName;
  }
};
const userData = ref();
onMounted(async () => {
  const info = await userStore.getUserInfo(t, true);
  if(info.data.value.idVerificationStatus === 2){
    await navigateTo(localePath('/MyAssets'));
  } else if(info.data.value.idVerificationStatus === 0){
    step.value = 1
    userData.value = JSON.parse(localStorage.getItem('userInfo'));
    const result = await userStore.getIdvOption(t);
    const data = result.data.value.countries;
    copyCountryList.value = [...data];
    countryList.value = [...data];
    countryList.value.unshift({
      chinese: '無',
      code: '',
      english: 'none',
      locale: '',
    });
    jobList.value = result.data.value.ics;
  }else{
    const id = Number(localStorage.getItem('id'));
    await navigateTo(localePath(`/Members/${id}/level`));
  }
});
</script>
<style>
.code::-webkit-outer-spin-button,
.code::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}
</style>
