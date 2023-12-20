<script setup>
const { modelValue, label, disabled, placeholder, allowClear } = defineProps(['modelValue', 'label', 'disabled', 'placeholder', 'allowClear']);
// const emit = defineEmits(['update:modelValue']);

const selectValue = ref(modelValue);
// const changeSelect = (value) => {
//   console.log('value :>> ', value);
//   emit('update:modelValue', value);
// };
</script>
<template>
  <!-- <FormBase :label="label"> -->
  <a-select v-model:value="selectValue" :disabled="disabled" :placeholder="placeholder" :allowClear="allowClear" class="bg-white rounded-lg w-full outline-0 transition-colors" @change="$emit('update:modelValue', selectValue)" :id="label">
    <slot />
  </a-select>

  <!-- </FormBase> -->
</template>
<style>
.ant-select-selector {
  height: 40px !important;
  padding: 5px 10px !important;
}
</style>
