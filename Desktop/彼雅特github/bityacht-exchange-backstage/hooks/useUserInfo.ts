import { useStorage } from "@vueuse/core";
import { decode } from "js-base64";
export default function useUserInfo() {
  const accessToken = useStorage("accessToken", "");
  if (!accessToken.value || accessToken.value === "") return { userInfo: null };
  let splitedJWT = accessToken.value?.split(".")[1]!;
  let userInfo = JSON.parse(decode(splitedJWT));
  userInfo.isExpired = userInfo.exp * 1000 < Date.now();
  return { userInfo };
}
