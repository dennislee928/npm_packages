<script setup>
import {
  EditOutlined,
  DeleteOutlined,
  PlusOutlined,
} from "@ant-design/icons-vue";
import draggable from "vuedraggable";
import useNotification from "@/hooks/useNotification";
import useAPI from "@/hooks/useAPI";
import usePagination from "@/hooks/usePagination";
import useParseDate from "@/hooks/useParseDate";

const toast = useNotification();
const dragOptions = {
  animation: 250,
  group: "banners",
  disabled: false,
  ghostClass: "opacity-50 bg-gray-500",
};
const parseDate = useParseDate();

const {
  getBanners,
  createBanner,
  updateBanner,
  deleteBanner,
  updateBannerPriority,
} = useAPI();
const addDialog = ref(null);
const editDialog = ref(null);
const deleteDialog = ref(false);

const loading = ref(true);
const form = ref({
  status: "-1",
});

const { data, page, perPage, totalPages, updatePagination } = usePagination();

onMounted(async () => {
  getData();
});
watch(
  [page, perPage, form],
  () => {
    getData();
  },
  { deep: true }
);
async function getData() {
  loading.value = true;
  let result = await getBanners(page.value, perPage.value, form.value);
  loading.value = false;
  if (result.error) return;

  updatePagination(result);
}
async function addBanner() {
  if (!document.getElementById("add-banner-form").reportValidity()) {
    return;
  }
  let formData = new FormData();
  let postData = deepCopy(addDialog.value);

  for (let key of [
    "status",
    "priority",
    "title",
    "subTitle",
    "button",
    "buttonUrl",
    "webImage",
    "appImage",
  ]) {
    formData.append(key, addDialog.value[key]);
  }
  if (!postData[`isPermanent`]) {
    formData.append("startAt", `"${parseDate(postData.startAt)}"`);
    formData.append("endAt", `"${parseDate(postData.endAt)}"`);
  }

  loading.value = true;
  let result = await createBanner(formData);
  loading.value = false;
  if (result.error) {
    toast.error("新增失敗");
    return;
  }
  addDialog.value = false;
  await getData();
}
async function changeBanner() {
  if (!document.getElementById("edit-banner-form").reportValidity()) {
    return;
  }
  let formData = new FormData();
  let postData = deepCopy(editDialog.value);

  for (let key of [
    "id",
    "status",
    "priority",
    "title",
    "subTitle",
    "button",
    "buttonUrl",
  ]) {
    formData.append(key, postData[key]);
  }
  if (!postData[`isPermanent`]) {
    formData.append("startAt", `"${parseDate(postData.startAt)}"`);
    formData.append("endAt", `"${parseDate(postData.endAt)}"`);
  }
  if (postData["_webImage"])
    formData.append("webImage", editDialog.value["_webImage"]);
  if (postData["_appImage"])
    formData.append("appImage", editDialog.value["_appImage"]);

  loading.value = true;
  let result = await updateBanner(postData.id, formData);
  loading.value = false;
  if (result.error) {
    if (result.code === 4010) {
      // CodeRecordNoChange
      toast.info("資料無變更");
    } else {
      toast.error("修改失敗");
    }
    return;
  }
  editDialog.value = false;
  await getData();
}

async function removeBanner() {
  loading.value = true;
  let result = await deleteBanner(deleteDialog.value.id);
  loading.value = false;
  if (result.error) {
    toast.error("刪除失敗");
    return;
  }
  deleteDialog.value = false;
  await getData();
}
async function updatePriority() {
  let result = await updateBannerPriority({
    rows: data.value.map((item, i) => ({ id: item.id, priority: i })),
  });
  if (result.error) {
    toast.error("修改失敗");
    loading.value = true;
    await getData();
    loading.value = false;
    return;
  }
}
const deepCopy = (x) => JSON.parse(JSON.stringify(x));
const convertToDateTimeLocalString = (date) => {
  const year = date.getFullYear();
  const month = (date.getMonth() + 1).toString().padStart(2, "0");
  const day = date.getDate().toString().padStart(2, "0");
  const hours = date.getHours().toString().padStart(2, "0");
  const minutes = date.getMinutes().toString().padStart(2, "0");

  return `${year}-${month}-${day}T${hours}:${minutes}`;
};
</script>
<template>
  <div>
    <PageTitle> Banner 管理</PageTitle>

    <div class="m-6 flex gap-2 justify-between items-end">
      <div class="flex gap-2 items-end">
        <FormSelect label="狀態" v-model="form.status">
          <option value="-1">全部</option>
          <option value="0">關閉</option>
          <option value="1">開啟</option>
        </FormSelect>
      </div>
      <div class="flex gap-2 items-end">
        <Button
          variant="secondary"
          @click="
            addDialog = {
              title: '',
              subTitle: '',
              status: 1,
              startAt: null,
              endAt: null,
              webImage: null,
              appImage: null,
              isPermanent: true,
              button: '',
              buttonUrl: '',
              priority: '',
            }
          "
        >
          <PlusOutlined /> 新增
        </Button>
      </div>
    </div>
    <TablePaginationContainer
      cols="2em 5em repeat(5, 1fr) 3em"
      v-model:perPage="perPage"
      v-model:page="page"
      :totalPages="totalPages"
      v-if="!loading"
    >
      <TableHeader>
        <div></div>
        <div>狀態</div>
        <div>WEB 圖片</div>
        <div>APP 圖片</div>
        <div>標題</div>
        <div>連結</div>
        <div>輪播時間</div>
        <div></div>
      </TableHeader>
      <draggable
        v-model="data"
        item-key="id"
        :dragOptions="dragOptions"
        handle=".handle"
        @end="updatePriority"
      >
        <template #item="{ element: item }">
          <TableItem>
            <div class="handle cursor-row-resize text-xl opacity-60">⠿</div>
            <div>
              <StatusBadge
                status="completed"
                variant="filled"
                v-if="item.status === 1"
              >
                開啟
              </StatusBadge>
              <StatusBadge
                status="cancelled"
                variant="filled"
                v-else-if="item.status === 0"
              >
                關閉
              </StatusBadge>
            </div>
            <div>
              <img
                :src="item.webImage"
                class="h-20 max-w-full object-contain"
              />
            </div>
            <div>
              <img
                :src="item.appImage"
                class="h-20 max-w-full object-contain"
              />
            </div>
            <div>
              {{ item.title }}
              <br />
              {{ item.subTitle }}
            </div>
            <div>
              {{ item.button }}
              <br />
              {{ item.buttonUrl }}
            </div>
            <div v-if="!item.isPermanent">
              <div>{{ item.startAt }}</div>
              -
              <br />
              <div>{{ item.endAt }}</div>
            </div>
            <div v-else>永久</div>
            <div class="flex justify-end gap-2">
              <EditOutlined
                @click="
                  editDialog = deepCopy({
                    ...item,
                    startAt: convertToDateTimeLocalString(
                      new Date(item.startAt)
                    ),
                    endAt: convertToDateTimeLocalString(new Date(item.endAt)),
                  })
                "
              />
              <DeleteOutlined @click="deleteDialog = deepCopy(item)" />
            </div>
          </TableItem>
        </template>
      </draggable>
    </TablePaginationContainer>
    <Loader v-if="loading" />
  </div>
  <Dialog v-model="editDialog" title="編輯 Banner" :loading="loading">
    <form
      @submit.prevent
      class="flex flex-col gap-4"
      v-if="editDialog"
      id="edit-banner-form"
    >
      <div class="grid grid-cols-[5rem_1fr] gap-4 items-center justify-end">
        <FormLabel label="狀態" />
        <div>
          <a-radio-group
            v-model:value="editDialog.status"
            :options="[
              { label: '開啟', value: 1 },
              { label: '關閉', value: 0 },
            ]"
          />
        </div>
        <FormLabel label="WEB 圖檔" />
        <FormInput
          type="file"
          id="webImage"
          @change="
            (e) => {
              editDialog._webImage = e.target.files[0];
            }
          "
        />
        <FormLabel label="手機版圖檔" />
        <FormInput
          type="file"
          id="appImage"
          @change="
            (e) => {
              editDialog._appImage = e.target.files[0];
            }
          "
        />
        <FormLabel label="標題" />
        <FormInput v-model="editDialog.title" required />
        <FormLabel label="子標題" />
        <FormInput v-model="editDialog.subTitle" required />
        <FormLabel label="按鈕文字" />
        <FormInput v-model="editDialog.button" required />
        <FormLabel label="按鈕連結" />
        <FormInput v-model="editDialog.buttonUrl" required />
        <FormLabel label="輪播時間" />
        <a-radio-group
          v-model:value="editDialog.isPermanent"
          :options="[
            { label: '永久', value: true },
            { label: '自訂', value: false },
          ]"
        />

        <FormLabel label="自訂時間" />
        <FormInput
          type="datetime-local"
          v-model="editDialog.startAt"
          :disabled="editDialog.isPermanent"
          required
        />
        <FormLabel label="至" />
        <FormInput
          type="datetime-local"
          v-model="editDialog.endAt"
          :disabled="editDialog.isPermanent"
          required
        />
      </div>
    </form>
    <template #actions>
      <Button variant="outline" @click="editDialog = false"> 取消 </Button>
      <Button @click="changeBanner"> 確定 </Button>
    </template>
  </Dialog>
  <Dialog v-model="addDialog" title="新增 Banner" :loading="loading">
    <form
      @submit.prevent
      class="flex flex-col gap-4"
      v-if="addDialog"
      id="add-banner-form"
    >
      <div class="grid grid-cols-[5rem_1fr] gap-4 items-center justify-end">
        <FormLabel label="狀態" />
        <div>
          <a-radio-group
            v-model:value="addDialog.status"
            :options="[
              { label: '開啟', value: 1 },
              { label: '關閉', value: 0 },
            ]"
          />
        </div>
        <FormLabel label="WEB 圖檔" />
        <FormInput
          type="file"
          id="webImage"
          @change="
            (e) => {
              addDialog.webImage = e.target.files[0];
            }
          "
          required
        />
        <FormLabel label="手機版圖檔" />
        <FormInput
          type="file"
          id="appImage"
          @change="
            (e) => {
              addDialog.appImage = e.target.files[0];
            }
          "
          required
        />
        <FormLabel label="標題" />
        <FormInput v-model="addDialog.title" required />
        <FormLabel label="子標題" />
        <FormInput v-model="addDialog.subTitle" required />
        <FormLabel label="按鈕文字" />
        <FormInput v-model="addDialog.button" required />
        <FormLabel label="按鈕連結" />
        <FormInput v-model="addDialog.buttonUrl" required />
        <FormLabel label="輪播時間" />
        <a-radio-group
          v-model:value="addDialog.isPermanent"
          :options="[
            { label: '永久', value: true },
            { label: '自訂', value: false },
          ]"
        />
        <FormLabel label="自訂時間" />
        <FormInput
          type="datetime-local"
          v-model="addDialog.startAt"
          :disabled="addDialog.isPermanent"
          required
        />
        <FormLabel label="至" />
        <FormInput
          type="datetime-local"
          v-model="addDialog.endAt"
          :disabled="addDialog.isPermanent"
          required
        />
      </div>
    </form>
    <template #actions>
      <Button variant="outline" @click="addDialog = false"> 取消 </Button>
      <Button @click="addBanner"> 確定 </Button>
    </template>
  </Dialog>
  <Dialog v-model="deleteDialog" title="剛除 Banner">
    <p>確定要刪除這個 Banner？</p>
    <template #actions>
      <Button variant="outline" @click="deleteDialog = false"> 取消 </Button>
      <Button @click="removeBanner"> 確定 </Button>
    </template>
  </Dialog>
</template>
