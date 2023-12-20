<script setup>
import { CopyOutlined } from "@ant-design/icons-vue";
import useNotification from "@/hooks/useNotification";

const toast = useNotification();

const { txid } = defineProps({
  txid: {
    type: String,
    required: true,
  },
});
async function copy() {
  try {
    await navigator.clipboard.writeText(txid);
    toast.success("已複製");
  } catch (e) {
    prompt("請複製以下交易編號", txid);
  }
}
</script>
<template>
  <div>
    <a-tooltip :mouseEnterDelay="0" :mouseLeaveDelay="0" color="white">
      <template #title>
        <span class="text-gray-950">
          {{ txid }}
        </span>
      </template>
      <div class="flex gap-1 items-center justify-center">
        {{ txid.slice(0, 6) }}...{{ txid.slice(-6) }}
        <button
          class="text-xs text-blue-600 hover:text-blue-800 hover:bg-gray-200 p-1 -m-1 active:text-blue-950 rounded-sm flex items-center justify-center"
          @click="copy"
        >
          <CopyOutlined />
        </button>
      </div>
    </a-tooltip>
  </div>
</template>
