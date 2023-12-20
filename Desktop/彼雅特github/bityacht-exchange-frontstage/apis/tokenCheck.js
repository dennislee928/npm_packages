import { parseJwt } from '@/config/config';
import { clearAllCookie } from '@/config/config';
const apiUrl = ({ selfAuth }) => {
  const urlText = {
    token: `${selfAuth}/token`,
  };
  return urlText;
};
const refreshToken = async (apiBase, selfAuth) => {
  const token = useCookie('token');
  const refreshToken = useCookie('refreshToken');
  const data = {
    id: Number(localStorage.getItem('id')),
    refreshToken: useCookie('refreshToken').value,
  };
  const result = await useFetch(apiUrl({ selfAuth }).token, {
    baseURL: apiBase,
    method: 'POST',
    body: data,
  });
  console.log('result :>> ', result);
  if (result.status.value === 'success') {
    const userData = parseJwt(result.data.value.accessToken);
    localStorage.setItem('userInfo', JSON.stringify(userData));
    token.value = result.data.value.accessToken;
    refreshToken.value = result.data.value.refreshToken;
  }
};

const tokenCheck = async (apiBase, selfAuth, refreshNow) => {
  const dayjs = useDayjs();
  const userInfo = JSON.parse(localStorage.getItem('userInfo'));
  const remainingTime = dayjs.unix(userInfo.exp);
  const now = dayjs();
  const duration = remainingTime.diff(now, 'minute');
  if (refreshNow) {
    await refreshToken(apiBase, selfAuth);
    return;
  }
  console.log('duration :>> ', duration);
  if (duration < 6 && duration >= 0) {
    await refreshToken(apiBase, selfAuth);
  } else if (duration < 0) {
    message.error('TOKEN過期，請重新登入');
    clearAllCookie();
    await navigateTo('/');
    setTimeout(() => {
      location.reload();
    }, 1000);
  }
};
export default tokenCheck;
