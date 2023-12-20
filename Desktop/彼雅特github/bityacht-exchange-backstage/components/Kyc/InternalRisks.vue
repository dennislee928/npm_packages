<script setup>
import useAPI from "@/hooks/useAPI";
import useNotification from "@/hooks/useNotification";

const {
  getKycRisks,
  createKycRisk,
  updateKycRisk,
  deleteKycRisk,
  getUserKycRisks,
  updateUserKycRisks,
} = useAPI();
const route = useRoute();
const toast = useNotification();
const { uid } = route.params;
const updateModelValue = defineEmits(["update:modelValue"]);

const { modelValue } = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
});

const loading = ref(false);
const riskList = ref(null);
const editMode = ref(false);
const editDialog = ref(false);
const editData = ref({
  id: -1,
  factor: "",
  subFactor: "",
  detail: "",
  score: 0,
});

async function getRisks() {
  loading.value = true;
  const data = await getKycRisks();
  const userRisks = await getUserKycRisks(uid);
  // tidy up data
  let result = {};
  Object.values(data).map((x) => {
    if (x) {
      // factor subFactor
      if (!result[x.factor]) {
        result[x.factor] = {};
      }
      if (!result[x.factor][x.subFactor]) {
        result[x.factor][x.subFactor] = [];
      }
      result[x.factor][x.subFactor].push({
        ...x,
        selected: userRisks.includes(x.id),
      });
    }
  });
  riskList.value = result;
  loading.value = false;
}
function deselectAll() {
  for (let [factor, subFactorList] of Object.entries(riskList.value)) {
    for (let [subFactor, items] of Object.entries(subFactorList)) {
      for (let item of items) {
        item.selected = false;
      }
    }
  }
  toast.success("已取消選取所有指標");
}
async function updateUserRisks() {
  loading.value = true;
  let data = [];
  for (let [factor, subFactorList] of Object.entries(riskList.value)) {
    for (let [subFactor, items] of Object.entries(subFactorList)) {
      for (let item of items) {
        if (item.selected) {
          data.push(item.id);
        }
      }
    }
  }
  let result = await updateUserKycRisks(uid, data);
  if (!result.error) {
    toast.success("更新成功");
    updateModelValue("update:modelValue", false);
  } else {
    toast.error("更新失敗");
  }
  loading.value = false;
}
function openAddDialog(factor = "", subFactor = "") {
  editDialog.value = true;
  editData.value = {
    id: -1,
    factor,
    subFactor,
    detail: "",
    score: 0,
  };
}
function openEditDialog(item) {
  editDialog.value = true;
  editData.value = deepCopy(item);
}
async function deleteRisk() {
  let { id } = editData.value;
  if (id > 0) {
    loading.value = true;
    let result = await deleteKycRisk(id);
    if (!result.error) {
      toast.success("刪除成功");
      editDialog.value = false;
      await getRisks();
    } else {
      toast.error("刪除失敗");
    }
    loading.value = false;
  }
}
async function submitEditData() {
  let { id, factor, subFactor, detail, score } = editData.value;
  score = Number(score);
  if (factor === "" || subFactor === "" || detail === "" || score > 0) {
    toast.error("請填寫完整");
    return;
  }
  loading.value = true;
  let result;
  if (id > 0) {
    result = await updateKycRisk({ id, factor, subFactor, detail, score });
  } else {
    result = await createKycRisk({ factor, subFactor, detail, score });
  }
  if (!result.error) {
    toast.success("更新成功");
    editDialog.value = false;
    await getRisks();
  } else {
    toast.error("更新失敗");
  }
  loading.value = false;
}
onMounted(() => {
  getRisks();
});
const deepCopy = (obj) => JSON.parse(JSON.stringify(obj));
</script>
<template>
  <Dialog
    :modelValue="modelValue"
    @update:modelValue="$emit('update:modelValue', $event)"
    :title="editMode ? '編輯指標' : '風險評估'"
    variant="sidebar"
    max-width="600px"
  >
    <div class="flex flex-col">
      <details
        class="flex flex-col gap-2 group"
        v-for="[title, items] of Object.entries(riskList)"
        v-if="riskList"
      >
        <summary
          class="bg-[#F1F3F6] hover:bg-[#e9ebee] active:bg-[#dddfe2] py-2 px-4 rounded cursor-pointer text-[#7489A4] flex items-center justify-between"
        >
          <div>{{ title }}</div>
          <svg
            class="w-4 h-4 transform group-open:rotate-180 transition-transform"
            viewBox="0 0 20 20"
            fill="currentColor"
          >
            <path
              fill-rule="evenodd"
              clip-rule="evenodd"
              d="M10 13.586L3.707 7.293a1 1 0 011.414-1.414L10 10.172l4.879-4.879a1 1 0 111.414 1.414L10 13.586z"
            />
          </svg>
        </summary>
        <details
          class="flex flex-col group/sub"
          v-for="[subtitle, subitems] of Object.entries(items)"
          open
        >
          <summary
            class="cursor-pointer flex gap-2 text-black items-center py-2 px-4 hover:text-opacity-80 active:text-opacity-90"
          >
            {{ subtitle }}
            <svg
              class="w-4 h-4 transform translate-y-0.5 group-open/sub:rotate-180 group-open/sub:translate-y-0 transition-transform text-[#ABB1BB] group-open/sub:text-[#3663A7]"
              viewBox="0 0 20 20"
              fill="currentColor"
            >
              <path
                fill-rule="evenodd"
                clip-rule="evenodd"
                d="M10 13.586L3.707 7.293a1 1 0 011.414-1.414L10 10.172l4.879-4.879a1 1 0 111.414 1.414L10 13.586z"
              />
            </svg>
          </summary>
          <div
            class="grid grid-cols-2 gap-2 mb-4 border-t border-gray-100 py-4 px-4 text-[#6B6C6C]"
          >
            <div class="flex items-start gap-2" v-for="item of subitems">
              <div v-if="!editMode" class="flex items-start gap-2">
                <label
                  :for="item.id"
                  class="relative flex cursor-pointer items-center rounded-full p-0 mt-[2.5px]"
                >
                  <input
                    type="checkbox"
                    class="before:content[''] peer relative h-5 w-5 cursor-pointer appearance-none rounded-md border border-blue-gray-200 transition-all before:absolute before:top-2/4 before:left-2/4 before:block before:h-12 before:w-12 before:-translate-y-2/4 before:-translate-x-2/4 before:rounded-full before:bg-blue-gray-500 before:opacity-0 before:transition-opacity checked:border-[#18E1B1] checked:bg-[#18E1B1] checked:before:bg-[#18E1B1]"
                    :key="item.id"
                    :id="item.id"
                    v-model="item.selected"
                  />
                  <div
                    class="pointer-events-none absolute top-2/4 left-2/4 -translate-y-2/4 -translate-x-2/4 text-white opacity-0 transition-opacity peer-checked:opacity-100"
                  >
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      class="h-3.5 w-3.5"
                      viewBox="0 0 20 20"
                      fill="currentColor"
                      stroke="currentColor"
                      stroke-width="1"
                    >
                      <path
                        fill-rule="evenodd"
                        d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                        clip-rule="evenodd"
                      ></path>
                    </svg>
                  </div>
                </label>
                <label
                  :for="item.id"
                  class="relative flex cursor-pointer items-center rounded-full p-0 leading-6"
                >
                  {{ item.detail }} ({{ item.score }})
                </label>
              </div>
              <div
                v-if="editMode"
                class="flex items-center justify-between gap-2 w-full border border-gray-100 hover:border-gray-200 active:border-gray-300 active:bg-[#F4F7FA] cursor-pointer rounded-md p-1"
                @click="openEditDialog(item)"
              >
                <div class="flex-1">{{ item.detail }}</div>
                <div
                  class="bg-gray-100 rounded p-1 w-7 h-7 flex items-center justify-center"
                >
                  {{ item.score }}
                </div>
              </div>
            </div>
            <div
              v-if="editMode"
              class="border border-dotted border-gray-300 rounded-md p-1 bg-[#F4F7FA] hover:border-gray-200 active:border-gray-300 hover:bg-[#eaeef2] active:border-solid flex items-center justify-center gap-2 cursor-pointer"
              @click="openAddDialog(title, subtitle)"
            >
              新增指標
            </div>
          </div>
        </details>
      </details>
      <Loader v-if="loading" />
    </div>
    <template #actions>
      <Button @click="editMode = !editMode" variant="outline" v-if="!editMode">
        編輯指標
      </Button>
      <Button @click="deselectAll" variant="outline" v-if="!editMode">
        取消選取所有指標
      </Button>
      <div class="flex-1" />
      <Button
        @click="$emit('update:modelValue', false)"
        variant="outline"
        v-if="!editMode"
      >
        取消
      </Button>
      <Button @click="editMode = !editMode" v-if="editMode"> 完成編輯 </Button>
      <Button @click="updateUserRisks" v-if="!editMode"> 完成 </Button>
    </template>
  </Dialog>
  <Dialog title="編輯指標" v-model="editDialog" :loading="loading">
    <form
      @submit.prevent
      class="grid grid-cols-[4rem_1fr] gap-4 items-center justify-end"
      id="edit-risk-form"
      v-if="riskList"
    >
      <FormLabel label="分類" required />

      <div>
        <a-auto-complete
          v-model:value="editData.factor"
          style="width: 100%"
          size="large"
          :options="Object.keys(riskList).map((x) => ({ value: x }))"
        />
      </div>
      <FormLabel label="子分類" required />
      <div>
        <a-auto-complete
          v-model:value="editData.subFactor"
          style="width: 100%"
          size="large"
          :options="
            Object.keys(riskList[editData.factor] || {}).map((x) => ({
              value: x,
            }))
          "
        />
      </div>
      <FormLabel label="描述" required />
      <FormInput v-model="editData.detail" required="" />
      <FormLabel label="分數" required />
      <FormInput v-model="editData.score" type="number" required />
    </form>
    <template #actions>
      <Button @click="deleteRisk" variant="danger" v-if="editData.id > 0">
        刪除
      </Button>
      <div class="flex-1" />
      <Button @click="editDialog = false" variant="outline"> 取消 </Button>
      <Button @click="submitEditData"> 完成</Button>
    </template>
  </Dialog>
</template>
