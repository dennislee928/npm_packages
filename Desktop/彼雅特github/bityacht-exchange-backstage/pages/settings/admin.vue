<script setup>
import {
  EditOutlined,
  DeleteOutlined,
  KeyOutlined,
  PlusOutlined,
} from "@ant-design/icons-vue";
import useNotification from "@/hooks/useNotification";
import useAPI from "@/hooks/useAPI";
import useManagersRoles from "@/hooks/useManagersRoles";
import usePagination from "@/hooks/usePagination";

const { getAdmins, createAdmin, updateAdmin, deleteAdmin } = useAPI();
const { getRoleNameFromId } = useManagersRoles();
const { data, page, perPage, totalPages, updatePagination } = usePagination();
const toast = useNotification();

const editDialog = ref(null);
const deleteDialog = ref(null);
const passwordDialog = ref(null);
const addDialog = ref(false);
const loading = ref(true);

onMounted(async () => {
  getData();
});
watch([page, perPage], () => {
  getData();
});
async function getData() {
  loading.value = true;
  let result = await getAdmins(page.value, perPage.value);
  loading.value = false;
  if (result.error) return;

  updatePagination(result);
}
async function addAdmin() {
  // validate
  if (!document.getElementById("add-admin-form").reportValidity()) {
    return;
  }

  loading.value = true;
  let result = await createAdmin(addDialog.value);
  if (!result.error) {
    toast.success("新增成功");
    addDialog.value = false;
    await getData();
  } else {
    toast.error("新增失敗");
  }
  loading.value = false;
}
async function editAdmin() {
  let result = await updateAdmin(editDialog.value);

  loading.value = true;
  if (!result.error) {
    toast.success("修改成功");
    editDialog.value = false;
    await getData();
  } else {
    if (result.code === 4010) {
      // CodeRecordNoChange
      toast.info("資料無變更");
    } else {
      toast.error("修改失敗");
    }
  }
  loading.value = false;
}
async function delAdmin() {
  let result = await deleteAdmin(deleteDialog.value.id);
  loading.value = true;
  if (!result.error) {
    toast.success("刪除成功");
    deleteDialog.value = false;
    await getData();
  } else {
    toast.error("刪除失敗");
  }
  loading.value = false;
}
async function changePassowrd() {
  if (!document.getElementById("change-password-form").reportValidity()) {
    return;
  }
  if (passwordDialog.value.password !== passwordDialog.value.confirmPassword) {
    toast.error("密碼不一致");
    return;
  }
  if (
    !passwordDialog.value.password.match(
      /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$/
    )
  ) {
    toast.error("密碼格式錯誤");
    return;
  }

  loading.value = true;
  let result = await updateAdmin({
    id: passwordDialog.value.id,
    password: passwordDialog.value.password,
  });
  if (!result.error) {
    toast.success("修改成功");
    passwordDialog.value = false;
    await getData();
  } else {
    toast.error("修改失敗");
  }
  loading.value = false;
}
const deepCopy = (x) => JSON.parse(JSON.stringify(x));
</script>
<template>
  <div>
    <div class="flex justify-between items-center">
      <PageTitle> 管理者帳號設定 </PageTitle>
      <div class="mx-4 flex gap-2 justify-end items-end">
        <Button
          @click="
            addDialog = {
              account: ``,
              managersRolesID: 1,
              name: ``,
            }
          "
        >
          <PlusOutlined /> 管理員
        </Button>
      </div>
    </div>
    <TablePaginationContainer
      cols="repeat(4, 1fr) 4em"
      v-model:perPage="perPage"
      v-model:page="page"
      :totalPages="totalPages"
      v-if="!loading"
    >
      <TableHeader>
        <div>帳號</div>
        <div>狀態</div>
        <div>角色</div>
        <div>姓名</div>
        <div></div>
      </TableHeader>
      <TableItem v-for="item in data" :key="item.id">
        <div>{{ item.account }}</div>
        <div>
          <StatusBadge status="completed" v-if="item.status">啟用</StatusBadge>
          <StatusBadge status="cancelled" v-else>停用</StatusBadge>
        </div>
        <div>{{ getRoleNameFromId(item.managersRolesID).zh }}</div>
        <div>{{ item.name }}</div>

        <div class="flex justify-end gap-2">
          <KeyOutlined
            @click="
              passwordDialog = deepCopy({
                ...item,
                password: '',
                confirmPassword: '',
              })
            "
          />
          <EditOutlined @click="editDialog = deepCopy(item)" />
          <DeleteOutlined @click="deleteDialog = deepCopy(item)" />
        </div>
      </TableItem>
    </TablePaginationContainer>
    <Loader v-else />
    <Dialog v-model="editDialog" :loading="loading" title="管理者設定">
      <div class="flex flex-col gap-4" v-if="editDialog">
        <div class="grid grid-cols-[2rem_1fr] gap-4 items-center justify-end">
          <FormLabel label="帳號" />
          <div>{{ editDialog.account }}</div>
          <FormLabel label="姓名" />
          <div>{{ editDialog.name }}</div>
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
          <FormLabel label="身份" />
          <div>
            <a-radio-group
              v-model:value="editDialog.managersRolesID"
              :options="
                Array.from({ length: 5 }, (_, i) => ({
                  label: getRoleNameFromId(i + 1).zh,
                  value: i + 1,
                }))
              "
            />
          </div>
        </div>
      </div>
      <template #actions>
        <Button variant="outline" @click="editDialog = false"> 取消 </Button>
        <Button @click="editAdmin"> 確定 </Button>
      </template>
    </Dialog>
    <Dialog v-model="addDialog" :loading="loading" title="新增管理者">
      <form
        @submit.prevent
        class="flex flex-col gap-4"
        v-if="addDialog"
        id="add-admin-form"
      >
        <div class="grid grid-cols-[2rem_1fr] gap-4 items-center justify-end">
          <FormLabel label="帳號" />
          <FormInput v-model="addDialog.account" type="email" />
          <FormLabel label="姓名" />
          <FormInput v-model="addDialog.name" />
          <FormLabel label="身份" />
          <div>
            <a-radio-group
              v-model:value="addDialog.managersRolesID"
              :options="
                Array.from({ length: 5 }, (_, i) => ({
                  label: getRoleNameFromId(i + 1).zh,
                  value: i + 1,
                }))
              "
            />
          </div>
        </div>
      </form>
      <template #actions>
        <Button variant="outline" @click="addDialog = false"> 取消 </Button>
        <Button @click="addAdmin"> 確定 </Button>
      </template>
    </Dialog>
    <Dialog v-model="deleteDialog" :loading="loading" title="剛除管理者">
      <p v-if="deleteDialog">
        確定要刪除 {{ deleteDialog.name }} 管理者帳號嗎？
      </p>
      <template #actions>
        <Button variant="outline" @click="deleteDialog = false"> 取消 </Button>
        <Button variant="danger" @click="delAdmin"> 確定 </Button>
      </template>
    </Dialog>
    <Dialog v-model="passwordDialog" :loading="loading" title="修改密碼">
      <form
        @submit.prevent
        class="grid grid-cols-[4rem_1fr] gap-4 items-center justify-end"
        id="change-password-form"
        v-if="passwordDialog"
      >
        <FormLabel label="設定密碼" />
        <FormInput v-model="passwordDialog.password" required />
        <FormLabel label="確認密碼" />
        <FormInput v-model="passwordDialog.confirmPassword" required />
        <div></div>
        <div class="text-sm opacity-75 -mt-2">
          密碼至少 8
          碼，必須包含至少一個英文大寫字母、一個英文小寫字母以及一個數字。
        </div>
      </form>
      <template #actions>
        <Button variant="outline" @click="passwordDialog = false">
          取消
        </Button>
        <Button @click="changePassowrd"> 確定 </Button>
      </template>
    </Dialog>
  </div>
</template>
