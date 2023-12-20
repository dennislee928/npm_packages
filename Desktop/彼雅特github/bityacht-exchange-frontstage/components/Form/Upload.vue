<script setup>
const { modelValue, label, type } = defineProps(['modelValue', 'label', 'type']);
const emit = defineEmits(['upload']);
const { t } = useI18n();

const handleChange = async (info) => {
  // console.log('info :>> ', info);
  const max = 16777215;
  if (info.file.size > max) {
    message.error(t('error.upload'));
    return false;
  } else {
    const base64 = await convertFileToBase64(info.file);
    const name = info.file.name;
    emit('upload', base64, type, name);
  }
};
const convertFileToBase64 = (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.readAsDataURL(file);
    reader.onload = () => resolve(reader.result);
    reader.onerror = (error) => reject(error);
  });
};
</script>
<template>
  <FormBase :label="label" class="w-full">
    <a-upload accept=".png,.jpg,.jpeg" :showUploadList="false" :customRequest="handleChange" class="relative">
      <slot />
    </a-upload>
  </FormBase>
</template>
<style>
.ant-upload {
  width: 100%;
}
</style>
