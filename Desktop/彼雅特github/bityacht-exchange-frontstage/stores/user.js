import { register, token, emailVerify, resendEmailVerify, login, logOut, forgotPassword, resetPassword, verifyResetPassword, getSpotTrend, getBanners, getUserInfo, getLoginRecord, editPassword, reSet2fa, login2fa, mobileBarcode, getIdvOption, sendPhoneVerify, phoneVerify, idv, commissions, checkPhone, commissionsWithdraw, patchImage, kryptoGoUrl, post2fa, patch2fa, get2fa, getBankOptions, setBank, deleteBank } from '@/apis/user';
import { parseJwt, getCookie } from '@/config/config';

const selfAuth = 'user';
const apiBase = useRuntimeConfig().public.apiBase;

const useUserStore = defineStore('user', {
  state: () => ({
    loadingButton: false,
  }),
  actions: {
    async getSpotTrend() {
      const apiResponse = await getSpotTrend(apiBase, selfAuth);
      return apiResponse;
    },
    async getBanners() {
      const apiResponse = await getBanners(apiBase, selfAuth);
      return apiResponse;
    },
    async getIdvOption(t) {
      const apiResponse = await getIdvOption(apiBase, selfAuth, t);
      return apiResponse;
    },
    async token(refreshNow) {
      await token(apiBase, selfAuth, refreshNow);
    },
    async getUserInfo(t, noRefresh = false) {
      const apiResponse = await getUserInfo(apiBase, selfAuth, t, noRefresh);
      return apiResponse;
    },
    async getLoginRecord(page, pageSize, t) {
      const apiResponse = await getLoginRecord(apiBase, page, pageSize, selfAuth, t);
      return apiResponse;
    },
    async sendPhoneVerify(accountData, t) {
      const apiResponse = await sendPhoneVerify(apiBase, accountData, selfAuth, t);
      return apiResponse;
    },
    async phoneVerify(accountData, t) {
      const apiResponse = await phoneVerify(apiBase, accountData, selfAuth, t);
      return apiResponse;
    },
    async idv(accountData, t) {
      const apiResponse = await idv(apiBase, accountData, selfAuth, t);
      return apiResponse;
    },
    async checkPhone(phone, t) {
      const apiResponse = await checkPhone(apiBase, phone, selfAuth, t);
      return apiResponse;
    },
    async register(accountData, t) {
      const apiResponse = await register(apiBase, accountData, selfAuth, t);
      return apiResponse;
    },
    async emailVerify(data, t) {
      const apiResponse = await emailVerify(apiBase, data, selfAuth, t);
      return apiResponse;
    },
    async resendEmailVerify(data, t) {
      const apiResponse = await resendEmailVerify(apiBase, data, selfAuth, t);
      return apiResponse;
    },
    async login(data, t) {
      const apiResponse = await login(apiBase, data, selfAuth, t);
      const isLogin = useCookie('isLogin');
      const token = useCookie('token');
      const refreshToken = useCookie('refreshToken');
      if (apiResponse?.data.value.accessToken) {
        const tokenValue = apiResponse.data.value.accessToken;
        const refreshTokenValue = apiResponse.data.value.refreshToken;
        token.value = `${tokenValue}`;
        refreshToken.value = `${refreshTokenValue}`;
        isLogin.value = true;
        const userData = parseJwt(apiResponse.data.value.accessToken);
        localStorage.setItem('userInfo', JSON.stringify(userData));
      }
      return apiResponse;
    },
    async reSet2fa(data, t) {
      const apiResponse = await reSet2fa(apiBase, data, selfAuth, t);
      return apiResponse;
    },
    async login2fa(data, t) {
      const apiResponse = await login2fa(apiBase, data, selfAuth, t);
      const isLogin = useCookie('isLogin');
      const token = useCookie('token');
      const refreshToken = useCookie('refreshToken');
      if (apiResponse?.data.value.accessToken) {
        const tokenValue = apiResponse.data.value.accessToken;
        const refreshTokenValue = apiResponse.data.value.refreshToken;
        token.value = `${tokenValue}`;
        refreshToken.value = `${refreshTokenValue}`;
        isLogin.value = true;
        const userData = parseJwt(apiResponse.data.value.accessToken);
        localStorage.setItem('userInfo', JSON.stringify(userData));
      }
      return apiResponse;
    },
    async logOut(t) {
      const apiResponse = await logOut(apiBase, selfAuth, t);
      return apiResponse;
    },
    async forgotPassword(data, t) {
      const apiResponse = await forgotPassword(apiBase, data, selfAuth, t);
      return apiResponse;
    },
    async verifyResetPassword(data) {
      const apiResponse = await verifyResetPassword(apiBase, data, selfAuth);
      return apiResponse;
    },
    async resetPassword(data, t) {
      const apiResponse = await resetPassword(apiBase, data, selfAuth, t);
      return apiResponse;
    },
    async editPassword(data, t) {
      const apiResponse = await editPassword(apiBase, data, selfAuth, t);
      return apiResponse;
    },
    async mobileBarcode(data, t) {
      const apiResponse = await mobileBarcode(apiBase, data, selfAuth, t);
      return apiResponse;
    },
    async commissions(page, pageSize, t, noRefresh = false) {
      const apiResponse = await commissions(apiBase, selfAuth, page, pageSize, t, noRefresh);
      return apiResponse;
    },
    async commissionsWithdraw(t) {
      const apiResponse = await commissionsWithdraw(apiBase, selfAuth, t);
      return apiResponse;
    },
    async patchImage(data, t) {
      const apiResponse = await patchImage(apiBase, data, selfAuth, t);
      return apiResponse;
    },
    async kryptoGoUrl(t) {
      const apiResponse = await kryptoGoUrl(apiBase, selfAuth, t);
      return apiResponse;
    },
    async post2fa(t) {
      const apiResponse = await post2fa(apiBase, selfAuth, t);
      return apiResponse;
    },
    async patch2fa(data, t) {
      const apiResponse = await patch2fa(apiBase, data, selfAuth, t);
      return apiResponse;
    },
    async get2fa(verificationCode, t) {
      const apiResponse = await get2fa(apiBase, verificationCode, selfAuth, t);
      return apiResponse;
    },
    async getBankOptions(t, noRefresh = false) {
      const apiResponse = await getBankOptions(apiBase, selfAuth, t, noRefresh);
      return apiResponse;
    },
    async setBank(data, t) {
      const apiResponse = await setBank(apiBase, data, selfAuth, t);
      return apiResponse;
    },
    async deleteBank(t) {
      const apiResponse = await deleteBank(apiBase, selfAuth, t);
      return apiResponse;
    },
  },
});
export default useUserStore;
