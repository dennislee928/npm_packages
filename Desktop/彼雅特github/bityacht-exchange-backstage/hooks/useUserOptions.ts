import useAPI from "./useAPI";
import { ref } from "vue";
export default function useUserOptions(option: "countries" | "ics") {
  const { getUserOptions } = useAPI();
  const result = ref<object>({});
  const loading = ref(false);
  const error = ref<any>(null);
  const fetch = async () => {
    loading.value = true;
    try {
      let res;
      if (sessionStorage.getItem("userOptions")) {
        res = JSON.parse(sessionStorage.getItem("userOptions") ?? "");
      } else {
        res = await getUserOptions();
      }
      sessionStorage.setItem("userOptions", JSON.stringify(res));
      result.value = res[option];
    } catch (err) {
      error.value = err;
    } finally {
      loading.value = false;
    }
  };
  fetch();
  return {
    result,
    loading,
    error,
  };
}
