<script setup>
const { modelValue, label, placeholder } = defineProps(['modelValue', 'label', 'placeholder']);
const emit = defineEmits(['changeDate']);

const dayjs = useDayjs();
const change = (value) => {
  const time = dayjs(value);
  const result = time.year() + '/' + Number(time.month() + 1) + '/' + time.date();
  const formatResult = dayjs(result).format('YYYY/MM/DD');
  emit('changeDate', formatResult);
};
</script>
<template>
  <FormBase :label="label">
    <a-date-picker :value="modelValue" class="bg-white border border-[#DDE1E5] rounded-md px-4 py-2 w-full outline-0 mt-2" :placeholder="placeholder" :id="label" format="YYYY/MM/DD" @change="change" />
  </FormBase>
</template>
