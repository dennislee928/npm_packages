<script setup>
const { cols, data } = defineProps({
  cols: {
    type: String,
    default: "250px 1fr",
  },
  data: {
    type: Array,
    required: true,
  },
});
const page = ref(1);
const perPage = ref(10);
const totalPages = computed(() => Math.ceil(data.length / perPage.value));
watch([page, perPage], () => {
  if (page.value > totalPages.value) {
    page.value = totalPages.value;
  }
});
</script>
<template>
  <div :style="`--cols: ${cols}`" class="m-6">
    <div class="overflow-x-auto">
      <div class="min-w-[1000px]">
        <slot :data="data.slice((page - 1) * perPage, page * perPage)" />
      </div>
    </div>
    <div class="flex gap-2 items-center mt-2">
      <button
        class="bg-[#DDE1E5] bg-opacity-0 hover:bg-opacity-50 active:bg-opacity-100 rounded-sm focus:outline-none w-8 h-8 text-center disabled:opacity-50 disabled:pointer-events-none"
        :disabled="page === 1"
        @click="page--"
      >
        ❮
      </button>
      <select
        v-model="page"
        class="border border-[#DDE1E5] rounded-sm focus:outline-none h-8"
      >
        <option v-for="i in totalPages" :key="i" :value="i">
          {{ i }}
        </option>
      </select>
      <button
        class="bg-[#DDE1E5] bg-opacity-0 hover:bg-opacity-50 active:bg-opacity-100 rounded-sm focus:outline-none w-8 h-8 text-center disabled:opacity-50 disabled:pointer-events-none"
        :disabled="page === totalPages"
        @click="page++"
      >
        ❯
      </button>

      <select
        v-model="perPage"
        class="border border-[#DDE1E5] rounded-sm focus:outline-none h-8"
      >
        <option v-for="i in [10, 20, 50, 100]" :key="i" :value="i">
          {{ i }} 筆 / 頁
        </option>
      </select>
    </div>
  </div>
</template>
