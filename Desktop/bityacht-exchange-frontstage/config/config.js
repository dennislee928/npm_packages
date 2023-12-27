const formatAccount = (value) => {
  const data = value.split('@');
  const first = data[0];
  const last = '@' + data[1];
  const formatValue = first.slice(0, -3);
  const replacedString = formatValue + 'xxx';
  return replacedString + last;
};
export default formatAccount;

export const formatPhone = (value) => {
  const first = value.slice(0, 2);
  const last = value.slice(-3);
  const result = first + '****' + last;
  return result;
};

export const parseJwt = (token) => {
  const base64Url = token.split('.')[1]; // 取得 JWT 的 payload 部分
  const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/'); // 替換 URL 安全的字符
  const jsonPayload = decodeURIComponent(
    atob(base64)
      .split('')
      .map((c) => {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2); // 解碼 base64 字符串
      })
      .join('')
  );
  return JSON.parse(jsonPayload); // 解析 JSON 字符串並返回
};
export const getCookie = (name, cookieName) => {
  const data = name.split(';');
  let result = '';
  data.forEach((item) => {
    let comparison = item.split('=');
    if (comparison[0] === cookieName) {
      result = comparison[1];
    }
  });
  return result;
};
export const clearAllCookie = () => {
  let date = new Date();
  date.setTime(date.getTime() - 10000);
  let keys = document.cookie.match(/[^ =;]+(?=\=)/g);
  if (keys) {
    for (var i = keys.length; i--; ) document.cookie = keys[i] + '=0; expire=' + date.toGMTString() + '; path=/';
  }
};
export const formatValue = (value) => {
  if (value === '-') return value;
  const number = parseFloat(value);
  // 格式化為千分位，保留兩位小數
  const formattedNumber = number.toLocaleString(undefined, {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });
  return formattedNumber;
};
export const formatValueByDigits = (value, digits) => {
  if (value === '-') return value;
  const number = parseFloat(value);
  // 格式化為千分位，保留兩位小數
  const formattedNumber = number.toLocaleString(undefined, {
    minimumFractionDigits: digits,
    maximumFractionDigits: digits,
  });
  return formattedNumber;
};
export const getAssetsFile = (url) => {
  return new URL(`../assets/img/${url}`, import.meta.url).href;
};
export const copyText = (str, t) => {
  const copy = str;
  navigator.clipboard
    .writeText(copy)
    .then(() => {
      message.success(t('memberCenter.copySuccess'));
    })
    .catch(() => {
      message.error(t('memberCenter.copyError'));
    });
};
export const shortString = (str) => {
  if (str !== undefined) {
    if (str.length > 20) {
      const start = str.slice(0, 3);
      const end = str.slice(-3);
      return start + '***' + end;
    } else {
      return str;
    }
  }
};
export const shortHref = (str) => {
  if (str !== undefined && str !== '') {
    const url = new URL(str);
    const domain = `${url.protocol}//${url.host}`;
    const inviteCode = str.slice(-7);
    return domain + '...' + inviteCode;
  }
};
