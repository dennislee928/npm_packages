//@ts-ignore
import * as Toast from "vue-toastification/dist/index.mjs";
const { useToast } = Toast;
export default function useNotification() {
  const toast = useToast();
  return toast;
}
