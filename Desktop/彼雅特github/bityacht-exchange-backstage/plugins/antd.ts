import { defineNuxtPlugin } from "#app";
import Antd from "ant-design-vue";

export default defineNuxtPlugin((nuxtApp) => {
  nuxtApp.vueApp.use(Antd);
});
