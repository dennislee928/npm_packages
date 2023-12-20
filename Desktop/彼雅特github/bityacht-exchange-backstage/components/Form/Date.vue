<script setup>
import useNotification from "@/hooks/useNotification";
const toast = useNotification();
const { startDate, endDate, label, required } = defineProps([
  "startDate",
  "endDate",
  "label",
  "required",
]);
//emit startDate and endDate to parent component
const start = ref(startDate);
const end = ref(endDate);
const $emit = defineEmits(["update:startDate", "update:endDate"]);
function updateDate(position, value) {
  if (position === "start") {
    if (end.value && value > end.value) {
      toast.error("開始日期不可晚於結束日期");
      start.value = "";
      $emit("update:startDate", "");
      return;
    }
    start.value = value;
    $emit("update:startDate", value);
  } else {
    if (value && value < start.value) {
      toast.error("結束日期不可早於開始日期");
      end.value = "";
      $emit("update:endDate", "");
      return;
    }
    end.value = value;
    $emit("update:endDate", value);
  }
}
</script>

<template>
  <FormBase :label="label">
    <div
      class="flex bg-white border border-[#DDE1E5] rounded-[8px] transition-colors ring-[#00b96b] ring-opacity-10 [&:has(:focus)]:border-[#00b96b] [&:has(:focus)]:ring-2"
    >
      <input
        type="date"
        class="w-full px-4 py-2 bg-transparent outline-0 rounded-md"
        :id="label"
        :value="startDate"
        @input="updateDate('start', $event.target.value)"
        :required="required"
      />
      <div class="py-2 opacity-50">-</div>
      <input
        type="date"
        class="w-full px-4 py-2 bg-transparent outline-0 rounded-md"
        :id="label"
        :value="endDate"
        @input="updateDate('end', $event.target.value)"
        :required="required"
      />
    </div>
  </FormBase>
</template>
