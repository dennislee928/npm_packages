import { getAssets, getHistories, getSpotOptions, getTrend, transactions, getTransactions, getAddress, produceAddress, withdraw, twofaWithdraw, withdrawOption, getWhiteList, postWhiteList, deleteWhiteList } from '@/apis/assets';

const selfAuth = 'user';
const apiBase = useRuntimeConfig().public.apiBase;

const useAssetsStore = defineStore('assets', {
  state: () => ({
    isLoadingTable: false,
    loadingButton: false,
  }),
  actions: {
    async getAssets(t) {
      const apiResponse = await getAssets(apiBase, selfAuth, t);
      return apiResponse;
    },
    async getHistories(action, currenciesSymbol, status, startAt = '""', endAt = '""', page, pageSize, t) {
      const apiResponse = await getHistories(apiBase, action, currenciesSymbol, status, startAt, endAt, page, pageSize, selfAuth, t);
      return apiResponse;
    },
    async getSpotOptions() {
      const apiResponse = await getSpotOptions(apiBase, selfAuth);
      return apiResponse;
    },
    async getTrend(t) {
      const apiResponse = await getTrend(apiBase, selfAuth, t);
      return apiResponse;
    },
    async transactions(data, t) {
      const apiResponse = await transactions(apiBase, data, selfAuth, t);
      return apiResponse;
    },
    async getTransactions(side, status, startAt = '""', endAt = '""', page, pageSize, t) {
      const apiResponse = await getTransactions(apiBase, side, status, startAt, endAt, page, pageSize, selfAuth, t);
      return apiResponse;
    },
    async getAddress(currencyType, mainnet) {
      const apiResponse = await getAddress(apiBase, currencyType, mainnet, selfAuth);
      return apiResponse;
    },
    async produceAddress(data, t) {
      const apiResponse = await produceAddress(apiBase, data, selfAuth, t);
      return apiResponse;
    },
    async withdraw(data, t, lastChange) {
      const apiResponse = await withdraw(apiBase, data, selfAuth, t, lastChange);
      return apiResponse;
    },
    async twofaWithdraw(data, t) {
      const apiResponse = await twofaWithdraw(apiBase, data, selfAuth, t);
      return apiResponse;
    },
    async withdrawOption(currencyType, mainnet, t) {
      const apiResponse = await withdrawOption(apiBase, currencyType, mainnet, selfAuth, t);
      return apiResponse;
    },
    async getWhiteList(mainnet, t) {
      const apiResponse = await getWhiteList(apiBase, mainnet, selfAuth, t);
      return apiResponse;
    },
    async postWhiteList(data, t) {
      const apiResponse = await postWhiteList(apiBase, data, selfAuth, t);
      return apiResponse;
    },
    async deleteWhiteList(id, t) {
      const apiResponse = await deleteWhiteList(apiBase, id, selfAuth, t);
      return apiResponse;
    },
  },
});
export default useAssetsStore;
