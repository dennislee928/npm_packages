import error from '@/config/error';
import tokenCheck from './tokenCheck';

const apiUrl = ({ selfAuth, page, pageSize, action, currenciesSymbol, status, side, startAt, endAt, currencyType, mainnet, id }) => {
  const urlText = {
    getAssets: `${selfAuth}/assets`,
    histories: `${selfAuth}/assets/spot/histories?action=${action}&currenciesSymbol=${currenciesSymbol}&status=${status}&startAt=${startAt}&endAt=${endAt}&page=${page}&pageSize=${pageSize}`,
    options: `${selfAuth}/assets/spot/options`,
    trend: `${selfAuth}/assets/spot/trend`,
    transactions: `${selfAuth}/assets/spot/transactions`,
    getTransactionsList: `${selfAuth}/assets/spot/transactions?side=${side}&status=${status}&startAt=${startAt}&endAt=${endAt}&page=${page}&pageSize=${pageSize}`,
    produceAddress: `${selfAuth}/wallets/deposit/address`,
    getAddress: `${selfAuth}/wallets/deposit/address?currencyType=${currencyType}&mainnet=${mainnet}`,
    withdraw: `${selfAuth}/wallets/withdraw`,
    twofaWithdraw: `${selfAuth}/wallets/2fa-withdraw`,
    withdrawOption: `${selfAuth}/wallets/withdraw-info?currencyType=${currencyType}&mainnet=${mainnet}`,
    getWhiteList: `${selfAuth}/wallets/withdrawal-whitelist?mainnet=${mainnet}`,
    postWhiteList: `${selfAuth}/wallets/withdrawal-whitelist`,
    deleteWhiteList: `${selfAuth}/wallets/withdrawal-whitelist/${id}`,
  };
  return urlText;
};
export const getAssets = async (apiBase, selfAuth, t) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth }).getAssets, {
    baseURL: apiBase,
    method: 'GET',
    headers: {
      Authorization: 'Bearer ' + useCookie('token').value,
      'Content-Type': 'application/json',
    },
  });
  if (response.error.value) {
    error(response.error.value.data, t);
  } else {
    return response;
  }
};
export const getHistories = async (apiBase, action, currenciesSymbol, status, startAt, endAt, page, pageSize, selfAuth, t) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth, action, currenciesSymbol, startAt, endAt, status, page, pageSize }).histories, {
    baseURL: apiBase,
    method: 'GET',
    headers: {
      Authorization: 'Bearer ' + useCookie('token').value,
      'Content-Type': 'application/json',
    },
  });
  if (response.error.value) {
    error(response.error.value.data, t);
  } else {
    return response;
  }
};
export const getSpotOptions = async (apiBase, selfAuth) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth }).options, {
    baseURL: apiBase,
    method: 'GET',
    headers: {
      Authorization: 'Bearer ' + useCookie('token').value,
      'Content-Type': 'application/json',
    },
  });
  if (response.error.value) {
    error(response.error.value.data);
  } else {
    return response;
  }
};
export const getTrend = async (apiBase, selfAuth, t) => {
  // await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth }).trend, {
    baseURL: apiBase,
    method: 'GET',
    headers: {
      Authorization: 'Bearer ' + useCookie('token').value,
      'Content-Type': 'application/json',
    },
  });
  if (response.error.value) {
    error(response.error.value.data, t);
  } else {
    return response;
  }
};
export const transactions = async (apiBase, data, selfAuth, t) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth }).transactions, {
    baseURL: apiBase,
    body: data,
    method: 'POST',
    headers: {
      Authorization: 'Bearer ' + useCookie('token').value,
      'Content-Type': 'application/json',
    },
    watch: false,
  });
  if (response.error.value) {
    error(response.error.value.data, t);
  } else {
    return response;
  }
};
export const getTransactions = async (apiBase, side, status, startAt, endAt, page, pageSize, selfAuth, t) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth, side, startAt, endAt, status, page, pageSize }).getTransactionsList, {
    baseURL: apiBase,
    method: 'GET',
    headers: {
      Authorization: 'Bearer ' + useCookie('token').value,
      'Content-Type': 'application/json',
    },
  });
  if (response.error.value) {
    error(response.error.value.data, t);
  } else {
    return response;
  }
};
export const getAddress = async (apiBase, currencyType, mainnet, selfAuth) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth, currencyType, mainnet }).getAddress, {
    baseURL: apiBase,
    method: 'GET',
    headers: {
      Authorization: 'Bearer ' + useCookie('token').value,
      'Content-Type': 'application/json',
    },
  });
  if (response.error.value) {
    return response;
  } else {
    return response;
  }
};
export const produceAddress = async (apiBase, data, selfAuth, t) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth }).produceAddress, {
    baseURL: apiBase,
    method: 'POST',
    body: data,
    headers: {
      Authorization: 'Bearer ' + useCookie('token').value,
      'Content-Type': 'application/json',
    },
    watch: false,
  });
  if (response.error.value) {
    error(response.error.value.data, t);
  } else {
    return response;
  }
};
export const withdraw = async (apiBase, data, selfAuth, t, lastChange) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth }).withdraw, {
    baseURL: apiBase,
    method: 'POST',
    body: data,
    headers: {
      Authorization: 'Bearer ' + useCookie('token').value,
      'Content-Type': 'application/json',
    },
    watch: false,
  });
  if (response.error.value) {
    error(response.error.value.data, t, lastChange);
  } else {
    return response;
  }
};
export const twofaWithdraw = async (apiBase, data, selfAuth, t) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth }).twofaWithdraw, {
    baseURL: apiBase,
    method: 'POST',
    body: data,
    headers: {
      Authorization: 'Bearer ' + useCookie('token').value,
      'Content-Type': 'application/json',
    },
    watch: false,
  });
  if (response.error.value) {
    error(response.error.value.data, t);
  } else {
    return response;
  }
};
export const withdrawOption = async (apiBase, currencyType, mainnet, selfAuth, t) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth, currencyType, mainnet }).withdrawOption, {
    baseURL: apiBase,
    method: 'GET',
    headers: {
      Authorization: 'Bearer ' + useCookie('token').value,
      'Content-Type': 'application/json',
    },
  });
  if (response.error.value) {
    error(response.error.value.data, t);
  } else {
    return response;
  }
};
export const getWhiteList = async (apiBase, mainnet, selfAuth, t) => {
  const response = await useFetch(apiUrl({ selfAuth, mainnet }).getWhiteList, {
    baseURL: apiBase,
    method: 'GET',
    headers: {
      Authorization: 'Bearer ' + useCookie('token').value,
      'Content-Type': 'application/json',
    },
  });
  if (response.error.value) {
    error(response.error.value.data, t);
  } else {
    return response;
  }
};
export const postWhiteList = async (apiBase, data, selfAuth, t) => {
  const response = await useFetch(apiUrl({ selfAuth }).postWhiteList, {
    baseURL: apiBase,
    method: 'POST',
    body: data,
    headers: {
      Authorization: 'Bearer ' + useCookie('token').value,
      'Content-Type': 'application/json',
    },
    watch: false,
  });
  if (response.error.value) {
    error(response.error.value.data, t);
  } else {
    return response;
  }
};
export const deleteWhiteList = async (apiBase, id, selfAuth, t) => {
  const response = await useFetch(apiUrl({ selfAuth, id }).deleteWhiteList, {
    baseURL: apiBase,
    method: 'DELETE',
    headers: {
      Authorization: 'Bearer ' + useCookie('token').value,
      'Content-Type': 'application/json',
    },
    watch: false,
  });
  if (response.error.value) {
    error(response.error.value.data, t);
  } else {
    return response;
  }
};
