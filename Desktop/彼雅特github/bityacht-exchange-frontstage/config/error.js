import { message } from 'ant-design-vue';
import { clearAllCookie } from '@/config/config';

const error = (value, t, lastChange) => {
  switch (value.code) {
    case 4000:
      clearAllCookie();
      message.error(t('error.actionError'));
      navigateTo('/');
      setTimeout(() => {
        location.reload();
      }, 1000);
      return false;
    case 4001:
      console.log('4001 :>> ', 4001);
      message.error(t('error.accountOrPassword'));
      return false;
    case 4002:
    case 4003:
      console.log('4002 || 4003 :>> ', 4002 || 4003);
      message.error(t('error.token'));
      clearAllCookie();
      navigateTo('/');
      setTimeout(() => {
        location.reload();
      }, 1000);
      return false;
    case 4004:
    case 4005:
      console.log('4004 || 4005 :>> ', 4005);
      clearAllCookie();
      navigateTo('/');
      message.error(t('error.actionError'));
      setTimeout(() => {
        location.reload();
      }, 1000);
      return false;
    case 4006:
      message.error(t('error.permissions'));
      return false;
    case 4007:
      message.error(t('error.accountRepeat'));
      return false;
    case 4011:
      console.log('4011 :>> ', 4011);
      if (value.data) {
        return value;
      } else {
        message.error(t('error.notNull'));
      }
      return false;
    case 4012:
      message.error(t('error.accountSuspension'));
      navigateTo('/suspension');
      return false;
    case 4022:
      console.log('4022 :>> ', 4022);
      message.error(t('error.errorCode'));
      return false;
    case 4023:
      console.log('4023 :>> ', 4023);
      message.error(t('error.codeValid'));
      setTimeout(() => {
        location.reload();
      }, 1000);
      return false;
    case 4024:
      console.log('4024 :>> ', 4024);
      message.error(t('error.mobileBarcodeError'));
      return false;
    case 4025:
      console.log('4025 :>> ', 4025);
      message.error(t('error.inviteCodeError'));
      return false;
    case 4030:
      console.log('4030 :>> ', 4030);
      message.error(t('error.IDRepeat'));
      return false;
    case 4031:
      console.log('4031 :>> ', 4031);
      message.error(t('error.idvIng'));
      return false;
    case 4038:
      console.log('4038 :>> ', 4038);
      if (lastChange) {
        message.error(t('error.changePassword24h') + ', ' + t('error.lastChange') + `${lastChange}`);
      } else {
        message.error(t('error.passwordErrorLimit'));
      }
      return false;
    case 4040:
      console.log('4040 :>> ', 4040);
      message.error(t('error.mobileBarcodeError'));
      return false;
    case 4041:
      console.log('4041 :>> ', 4041);
      message.error(t('error.addressError'));
      return false;
    case 4042:
      console.log('4042 :>> ', 4042);
      message.error(t('error.phoneError'));
      return false;
    case 4043:
      console.log('4043 :>> ', 4043);
      message.error(t('error.addressOverAdd'));
      return false;
    case 5020:
      console.log('5020 :>> ', 5020);
      message.error(t('error.apiError'));
      return false;
    case 5026:
      console.log('5026 :>> ', 5026);
      message.error(t('error.MinistryOfFinanceError'));
      return false;
    case 5036:
      console.log('5036 :>> ', 5036);
      message.error(t('error.walletAlreadySetting'));
      return false;
    default:
      message.error(t('error.errorUnknow'));
      return false;
  }
};
export default error;
