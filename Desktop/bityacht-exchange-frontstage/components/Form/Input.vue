<script setup>
const { modelValue, label, placeholder, type, password, disabled, addressError, readonly } = defineProps(['modelValue', 'label', 'placeholder', 'type', 'password', 'disabled', 'addressError', 'readonly']);
const inputElement = ref(null);
const done = () => {
  inputElement.value.blur();
};
</script>
<template>
  <FormBase :label="label">
    <template v-if="password">
      <a-input-password :value="modelValue" class="bg-white border border-[#DDE1E5] rounded-md px-4 py-2 w-full outline-0 mt-2" @input="$emit('update:modelValue', $event.target.value)" :placeholder="placeholder" :id="label" :type="type" @change="$emit('change', $event.target.value)" :disabled="disabled" />
    </template>
    <template v-else>
      <a-input :value="modelValue" ref="inputElement" class="bg-white border border-[#DDE1E5] rounded-md px-4 py-2 w-full outline-0 mt-2 whitespace-nowrap text-ellipsis overflow-hidden pr-6" :class="{ 'border-[#ff7875]': addressError }" :placeholder="placeholder" :id="label" :type="type" :disabled="disabled" :readonly="readonly" @keydown.enter="done()" @input="$emit('update:modelValue', $event.target.value)" />
    </template>
    <slot />
  </FormBase>
</template>
