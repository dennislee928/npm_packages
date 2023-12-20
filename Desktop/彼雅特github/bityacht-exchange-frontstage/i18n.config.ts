import zh_TW from "./locale/zh-TW.js";
import en_US from "./locale/en-US.js";
export default defineI18nConfig(() => ({
  legacy: false,
  locale: "zh-TW",
  globalInjection: true,
  messages: {
    "zh-TW": zh_TW,
    "en-US": en_US,
  },
}));
