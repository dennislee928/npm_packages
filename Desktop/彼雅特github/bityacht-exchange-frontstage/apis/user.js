import error from '@/config/error';
import tokenCheck from './tokenCheck';

const apiUrl = ({ selfAuth, verificationCode, page, pageSize }) => {
  const urlText = {
    register: `${selfAuth}/register`,
    emailVerify: `${selfAuth}/verify`,
    resendEmailVerify: `${selfAuth}/resend-verify`,
    login: `${selfAuth}/login`,
    token: `${selfAuth}/token`,
    logOut: `${selfAuth}/logout`,
    forgotPassword: `${selfAuth}/forgot-password`,
    resetPassword: `${selfAuth}/reset-password`,
    verifyResetPassword: `${selfAuth}/verify-reset-password`,
    getSpotTrend: `${selfAuth}/spot-trend`,
    getBanners: `${selfAuth}/banners`,
    getUserInfo: `${selfAuth}/info`,
    getLoginRecord: `${selfAuth}/login-logs?page=${page}&pageSize=${pageSize}`,
    editPassword: `${selfAuth}/settings/password`,
    reSet2fa: `${selfAuth}/settings/2fa`,
    login2fa: `${selfAuth}/2fa-login`,
    mobileBarcode: `${selfAuth}/settings/mobile-barcode`,
    getIdvOption: `${selfAuth}/idv/options`,
    checkPhone: `${selfAuth}/idv/check-phone`,
    sendPhoneVerify: `${selfAuth}/idv/issue-phone-verify`,
    phoneVerify: `${selfAuth}/idv/verify-phone`,
    idv: `${selfAuth}/idv`,
    commissions: `${selfAuth}/commissions?page=${page}&pageSize=${pageSize}`,
    commissionsWithdraw: `${selfAuth}/commissions/withdraw`,
    patchImage: `${selfAuth}/idv/image`,
    kryptoGoUrl: `${selfAuth}/idv/krypto-go-url`,
    post2fa: `${selfAuth}/settings/issue-withdraw-2fa-verify`,
    patch2fa: `${selfAuth}/settings/withdraw-2fa`,
    get2fa: `${selfAuth}/settings/withdraw-2fa-info?verificationCode=${verificationCode}`,
    getBankOptions: `${selfAuth}/banks/options`,
    setBank: `${selfAuth}/banks/account`,
    deleteBank: `${selfAuth}/banks/account`,
  };
  return urlText;
};

export const getSpotTrend = async (apiBase, selfAuth) => {
  const response = await useFetch(apiUrl({ selfAuth }).getSpotTrend, {
    baseURL: apiBase,
    method: 'GET',
  });
  return response;
};
export const getBanners = async (apiBase, selfAuth) => {
  const response = await useFetch(apiUrl({ selfAuth }).getBanners, {
    baseURL: apiBase,
    method: 'GET',
  });
  return response;
};
export const token = async (apiBase, selfAuth, refreshNow) => {
  await tokenCheck(apiBase, selfAuth, refreshNow);
};
export const getIdvOption = async (apiBase, selfAuth, t) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth }).getIdvOption, {
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
export const sendPhoneVerify = async (apiBase, accountData, selfAuth, t) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth }).sendPhoneVerify, {
    baseURL: apiBase,
    method: 'POST',
    body: accountData,
    headers: {
      Authorization: 'Bearer ' + useCookie('token').value,
      'Content-Type': 'application/json',
    },
    watch: false,
  });
  if (response.error.value) {
    const result = error(response.error.value.data, t);
    return result;
  } else {
    return response;
  }
};
export const phoneVerify = async (apiBase, accountData, selfAuth, t) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth }).phoneVerify, {
    baseURL: apiBase,
    method: 'POST',
    body: accountData,
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
export const idv = async (apiBase, accountData, selfAuth, t) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth }).idv, {
    baseURL: apiBase,
    method: 'POST',
    body: accountData,
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
export const checkPhone = async (apiBase, phone, selfAuth, t) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth }).checkPhone, {
    baseURL: apiBase,
    method: 'POST',
    body: phone,
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
export const getUserInfo = async (apiBase, selfAuth, t, noRefresh) => {
  if (!noRefresh) {
    await tokenCheck(apiBase, selfAuth);
  }
  const response = await useFetch(apiUrl({ selfAuth }).getUserInfo, {
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
export const getLoginRecord = async (apiBase, page, pageSize, selfAuth, t) => {
  // await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ page, pageSize, selfAuth }).getLoginRecord, {
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
export const register = async (apiBase, accountData, selfAuth, t) => {
  const response = await useFetch(apiUrl({ selfAuth }).register, {
    baseURL: apiBase,
    method: 'POST',
    body: accountData,
    watch: false,
  });
  if (response.error.value) {
    error(response.error.value.data, t);
  } else {
    return response;
  }
};
export const emailVerify = async (apiBase, accountData, selfAuth, t) => {
  const response = await useFetch(apiUrl({ selfAuth }).emailVerify, {
    baseURL: apiBase,
    method: 'POST',
    body: accountData,
    watch: false,
  });
  if (response.error.value) {
    // error(response.error.value.data);
    return response.error.value.data;
  } else {
    return response;
  }
};
export const resendEmailVerify = async (apiBase, accountData, selfAuth, t) => {
  const response = await useFetch(apiUrl({ selfAuth }).resendEmailVerify, {
    baseURL: apiBase,
    method: 'POST',
    body: accountData,
    watch: false,
  });
  if (response.error.value) {
    error(response.error.value.data, t);
  } else {
    return response;
  }
};
export const login = async (apiBase, accountData, selfAuth, t) => {
  const response = await useFetch(apiUrl({ selfAuth }).login, {
    baseURL: apiBase,
    method: 'POST',
    body: accountData,
    watch: false,
  });
  if (response.error.value) {
    error(response.error.value.data, t);
  } else {
    return response;
  }
};
export const logOut = async (apiBase, selfAuth, t) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth }).logOut, {
    baseURL: apiBase,
    method: 'POST',
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
export const verifyResetPassword = async (apiBase, accountData, selfAuth) => {
  const response = await useFetch(apiUrl({ selfAuth }).verifyResetPassword, {
    baseURL: apiBase,
    method: 'POST',
    body: accountData,
    watch: false,
  });
  if (response.error.value) {
    // error(response.error.value.data);
    return response.error.value.data;
  } else {
    return response;
  }
};
export const forgotPassword = async (apiBase, accountData, selfAuth, t) => {
  const response = await useFetch(apiUrl({ selfAuth }).forgotPassword, {
    baseURL: apiBase,
    method: 'POST',
    body: accountData,
    watch: false,
  });
  if (response.error.value) {
    error(response.error.value.data, t);
  } else {
    return response;
  }
};
export const resetPassword = async (apiBase, accountData, selfAuth, t) => {
  const response = await useFetch(apiUrl({ selfAuth }).resetPassword, {
    baseURL: apiBase,
    method: 'POST',
    body: accountData,
    watch: false,
  });
  if (response.error.value) {
    error(response.error.value.data, t);
  } else {
    return response;
  }
};
export const editPassword = async (apiBase, accountData, selfAuth, t) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth }).editPassword, {
    baseURL: apiBase,
    method: 'PATCH',
    body: accountData,
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
export const reSet2fa = async (apiBase, accountData, selfAuth, t) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth }).reSet2fa, {
    baseURL: apiBase,
    method: 'PATCH',
    body: accountData,
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
export const login2fa = async (apiBase, accountData, selfAuth, t) => {
  const response = await useFetch(apiUrl({ selfAuth }).login2fa, {
    baseURL: apiBase,
    method: 'POST',
    body: accountData,
    watch: false,
  });
  if (response.error.value) {
    error(response.error.value.data, t);
  } else {
    return response;
  }
};
export const mobileBarcode = async (apiBase, accountData, selfAuth, t) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth }).mobileBarcode, {
    baseURL: apiBase,
    method: 'PATCH',
    body: accountData,
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
export const commissions = async (apiBase, selfAuth, page, pageSize, t, noRefresh) => {
  if (!noRefresh) {
    await tokenCheck(apiBase, selfAuth);
  }
  const response = await useFetch(apiUrl({ page, pageSize, selfAuth }).commissions, {
    baseURL: apiBase,
    method: 'GET',
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
export const commissionsWithdraw = async (apiBase, selfAuth, t) => {
  const response = await useFetch(apiUrl({ selfAuth }).commissionsWithdraw, {
    baseURL: apiBase,
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
export const patchImage = async (apiBase, accountData, selfAuth, t) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth }).patchImage, {
    baseURL: apiBase,
    method: 'PATCH',
    body: accountData,
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
export const kryptoGoUrl = async (apiBase, selfAuth, t) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth }).kryptoGoUrl, {
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
export const post2fa = async (apiBase, selfAuth, t) => {
  const response = await useFetch(apiUrl({ selfAuth }).post2fa, {
    baseURL: apiBase,
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
export const patch2fa = async (apiBase, data, selfAuth, t) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth }).patch2fa, {
    baseURL: apiBase,
    method: 'PATCH',
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
export const get2fa = async (apiBase, verificationCode, selfAuth, t) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth, verificationCode }).get2fa, {
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
export const getBankOptions = async (apiBase, selfAuth, t, noRefresh) => {
  if (!noRefresh) {
    await tokenCheck(apiBase, selfAuth);
  }
  const response = await useFetch(apiUrl({ selfAuth }).getBankOptions, {
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
export const setBank = async (apiBase, data, selfAuth, t) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth }).setBank, {
    baseURL: apiBase,
    method: 'PUT',
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
export const deleteBank = async (apiBase, selfAuth, t) => {
  await tokenCheck(apiBase, selfAuth);
  const response = await useFetch(apiUrl({ selfAuth }).deleteBank, {
    baseURL: apiBase,
    method: 'DELETE',
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
