import { defineNuxtPlugin } from "#imports";
export default defineNuxtPlugin(() => {
  __DEBUG__ && console.log("load $i18n type definition plugin for legacy mode");
});
