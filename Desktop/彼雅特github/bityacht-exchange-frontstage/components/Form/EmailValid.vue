<template>
  <div class="flex items-center justify-center mt-3">
    <input v-for="(code, index) in codes" :key="index" type="number" inputmode="numeric" min="0" max="9" required ref="codeInput" v-model="codes[index]" @input="handleInput(index, $event)" @paste="handlePaste(index, $event)" @keydown="handleKeyDown(index, $event)" @click="handleClick(index,$event)" class="code bg-grayBg text-black text-h5 xxs:mx-2 mx-1 xxs:w-[48px] w-[30px] xxs:h-[92px] h-[60px] rounded-xl text-center" />
  </div>
</template>
<script setup>
const { codes } = defineProps({
  codes: {
    type: Array,
    default: [],
  },
});
const codeInput = ref();
const timer = ref(null);
const indexNow = ref(0);
const handlePaste = (index, event) => {
  const clipboardData = event.clipboardData || window.Clipboard;
  const pastedText = clipboardData.getData('text');
  if (pastedText.length === 6) {
    const text = pastedText.split('');
    for (let i = 0; i < pastedText.length; i++) {
      codes[i] = text[i];
    }
    event.preventDefault();
  }
};
const handleClick = (index,event) => {
  indexNow.value = index
}
const handleKeyDown = (index, event) => {
  indexNow.value = index
  if (event.key === 'Backspace') {
    timer.value = setTimeout(() => {
      if (indexNow.value > 0) {
        const prevInput = event.target.previousElementSibling;
        if (prevInput) {
          prevInput.focus();
        }
      }
    }, 100);
  }
};
const handleInput = (index, event) => {
  // const inputChar = event.target.value;
  // const regex = /^[0-9]$/;
  // if (event.key >= 0 && event.key <= 9) {
  //   codes.splice(index, 1, '');
  //   timer.value = setTimeout(() => {
  //     if (index < codes.length - 1) {
  //       codeInput.value[index + 1].focus();
  //     }
  //   }, 10);
  // } else if (event.key === 'Backspace') {
  //   timer.value = setTimeout(() => {
  //     if (index > 0) {
  //       codeInput.value[index - 1].focus();
  //     }
  //   }, 10);
  // }
  // console.log('index :>> ', index);
  // indexNow.value = index
  const inputValue = event.target.value;
  event.target.value = inputValue.slice(-1);
  if (!/^\d+$/.test(inputValue)) {
    event.target.value = '';
    return;
  }
  if ( indexNow.value < codeInput.value.length - 1 && inputValue !== '') {
    codeInput.value[indexNow.value + 1] = inputValue;
    codeInput.value[ indexNow.value] = '';
    const nextInput = event.target.nextElementSibling;
    if (nextInput) {
      nextInput.focus();
    }
  }
};
</script>
