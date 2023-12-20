<script setup>
import useAPI from "@/hooks/useAPI";
const { id } = defineProps({
  id: {
    type: String,
  },
});
const { getReceipt } = useAPI();
const loading = ref(false);
const receipt = ref(null);
const receiptMonth = computed(() => {
  let invoiceIssuedAt = receipt.value?.invoiceIssuedAt;
  if (!invoiceIssuedAt) return "";

  let year = invoiceIssuedAt.slice(0, 4) - 1911;
  let month = invoiceIssuedAt.slice(5, 7);
  let startMonth = month;
  let endMonth = month;
  if (month % 2 == 0) {
    startMonth = parseInt(month) - 1;
  } else {
    endMonth = parseInt(month) + 1;
  }
  const zeroPad = (num) => String(num).padStart(2, "0");
  return `${year} 年 ${zeroPad(startMonth)}-${zeroPad(endMonth)} 月`;
});
onMounted(async () => {
  getData();
});
async function getData() {
  loading.value = true;
  receipt.value = null;
  if (id) {
    const res = await getReceipt(id);
    receipt.value = res;
  }
  loading.value = false;
}
</script>
<template>
  <div
    class="grid grid-cols-[9rem_1fr] gap-4 items-center justify-end"
    v-if="receipt"
  >
    <FormLabel label="發票月份：" />
    <div v-if="receipt.invoiceID != ''">{{ receiptMonth }}</div>
    <div v-else class="opacity-50">尚未開立</div>
    <FormLabel label="發票號碼：" />
    <div v-if="receipt.invoiceID != ''">{{ receipt.invoiceID }}</div>
    <div v-else class="opacity-50">尚未開立</div>
    <FormLabel label="品名：" />
    <div>平台手續費</div>
    <FormLabel label="銷售額 (未稅價)：" />
    <div>{{ receipt.salesAmount }}</div>
    <FormLabel label="營業稅額：" />
    <div>{{ receipt.tax }}</div>
    <FormLabel label="收款金額：" />
    <div>{{ receipt.invoiceAmount }}</div>
    <FormLabel label="載具編號：" />
    <div>{{ receipt.barcode }}</div>
  </div>
  <div class="grid grid-cols-[9rem_1fr] gap-4 items-center justify-end" v-else>
    <FormLabel label="發票月份：" />
    <div class="opacity-50">-</div>
    <FormLabel label="發票號碼：" />
    <div class="opacity-50">-</div>
    <FormLabel label="品名：" />
    <div>平台手續費</div>
    <FormLabel label="銷售額 (未稅價)：" />
    <div class="opacity-50">-</div>
    <FormLabel label="營業稅額：" />
    <div class="opacity-50">-</div>
    <FormLabel label="收款金額：" />
    <div class="opacity-50">-</div>
    <FormLabel label="載具編號：" />
    <div class="opacity-50">-</div>
  </div>
</template>
