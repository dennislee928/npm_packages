<script setup>
import useAPI from "@/hooks/useAPI";
import useNotification from "@/hooks/useNotification";
import usePagination from "@/hooks/usePagination";
import useParseDate from "@/hooks/useParseDate";
import useUserInfo from "@/hooks/useUserInfo";
import useUserOptions from "@/hooks/useUserOptions";
import { PlusOutlined } from "@ant-design/icons-vue";

const toast = useNotification();
const { getUsers, createJuridicalUser, getExportUsersUrl } = useAPI();
const { result: industrialList, loading: industrialLoading } =
  useUserOptions("ics");
const { result: countryList, loading: countryLoading } =
  useUserOptions("countries");
const parseDate = useParseDate();
const { userInfo } = useUserInfo();

const form = ref({
  status: "-1",
  type: "-1",
  startAt: "",
  endAt: "",
  search: "",
});
const exportDialog = ref(false);
const legalDialog = ref(false);
const loading = ref(true);

const { data, page, perPage, totalPages, totalRecord, updatePagination } =
  usePagination();

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
  let options = {};
  for (let [key, value] of Object.entries(form.value)) {
    if (value !== "-1" && value !== "") {
      options[key] = value;
    }
  }
  if (options?.startAt || options?.endAt) {
    if (!options.startAt) {
      options.startAt = "1990/01/01";
    }
    if (!options.endAt) {
      options.endAt = new Date()
        .toISOString()
        .slice(0, 10)
        .replaceAll("-", "/");
    }
    options.startAt = parseDate(options.startAt, { dateOnly: true });
    options.endAt = parseDate(options.endAt, { dateOnly: true });
  }
  let result = await getUsers(page.value, perPage.value, options);
  loading.value = false;
  if (result.error) return;

  updatePagination(result);
}

async function addJuridicalUser() {
  // validate
  if (!document.getElementById("add-juridical-user-form").reportValidity()) {
    return;
  }
  loading.value = true;
  let postData = JSON.parse(JSON.stringify(legalDialog.value));
  postData.birthDate = postData.birthDate.replaceAll("-", "/");
  let result = await createJuridicalUser(postData);
  loading.value = false;
  if (!result.error) {
    toast.success("新增成功");
    legalDialog.value = false;
    await getData();
  } else {
    toast.error("新增失敗");
  }
}
function exportUsersData() {
  if (!document.getElementById("export-users-form").reportValidity()) {
    return;
  }
  window.open(
    getExportUsersUrl({
      statusList: exportDialog.value.status,
      startAt: parseDate(exportDialog.value.startAt, { dateOnly: true }),
      endAt: parseDate(exportDialog.value.endAt, { dateOnly: true }),
    })
  );
  exportDialog.value = false;
}
</script>
<template>
  <div>
    <PageTitle>
      會員管理
      <template #action> 會員數：{{ totalRecord.toLocaleString() }}</template>
    </PageTitle>

    <div class="m-6 flex gap-2 justify-between items-end">
      <div class="flex gap-2 items-end">
        <FormSelect label="狀態" v-model="form.status">
          <option value="-1">全部</option>
          <option value="0">未啟用</option>
          <option value="1">已啟用</option>
          <option value="2">已停用</option>
          <option value="3">凍結中</option>
        </FormSelect>
        <FormSelect label="類型" v-model="form.type">
          <option value="-1">全部</option>
          <option value="1">自然人</option>
          <option value="2">法人</option>
        </FormSelect>
        <FormDate
          label="日期"
          v-model:start-date="form.startAt"
          v-model:end-date="form.endAt"
        />
        <FormInput
          label="搜尋"
          v-model="form.search"
          placeholder="UID / E-MAIL / 姓名 / 手機"
        />
        <Button
          @click="
            form = {
              status: '-1',
              type: '-1',
              startAt: '',
              endAt: '',
              search: '',
            }
          "
        >
          重設
        </Button>
      </div>
      <div class="flex gap-2 items-end">
        <Button
          variant="export"
          @click="
            exportDialog = {
              startAt: new Date(new Date() - 30 * 24 * 60 * 60 * 1000)
                .toISOString()
                .slice(0, 10),
              endAt: new Date().toISOString().slice(0, 10),
              status: [0, 1, 2, 3],
            }
          "
          v-show="[4, 2, 1].includes(userInfo.managersRolesID)"
        >
          匯出
        </Button>
        <Button
          variant="secondary"
          @click="
            legalDialog = {
              account: '',
              address: '',
              authorizedPersonName: '',
              authorizedPersonNationalID: '',
              authorizedPersonPhone: '',
              birthDate: '',
              comment: '',
              country: '',
              industrialClassificationsID: 1,
              juridicalPersonCryptocurrencySources: '',
              juridicalPersonFundsSources: '',
              juridicalPersonNature: '',
              name: '',
              nationalID: '',
              phone: '',
            }
          "
          v-show="[4, 2, 1].includes(userInfo.managersRolesID)"
        >
          <PlusOutlined /> 法人
        </Button>
      </div>
    </div>
    <TablePaginationContainer
      cols="6em 5em 50px 1.5fr 1.2fr 1fr 3em 11em"
      v-model:perPage="perPage"
      v-model:page="page"
      :totalPages="totalPages"
      v-if="!loading"
    >
      <TableHeader>
        <div>UID</div>
        <div>狀態</div>
        <div>等級</div>
        <div>E-MAIL</div>
        <div>手機</div>
        <div>姓名</div>
        <div>身份</div>
        <div>註冊日期</div>
      </TableHeader>
      <TableItem v-for="item in data" :key="item">
        <nuxt-link :to="`/members/${item.id}/info`" class="link">
          {{ item.id }}
        </nuxt-link>
        <div>
          <StatusBadge status="in-progress" v-if="item.status === 0">
            未啟用
          </StatusBadge>
          <StatusBadge status="completed" v-if="item.status === 1">
            已啟用
          </StatusBadge>
          <StatusBadge status="cancelled" v-if="item.status === 2">
            已停用
          </StatusBadge>
          <StatusBadge status="frozen" v-if="item.status === 3">
            凍結中
          </StatusBadge>
        </div>
        <div>LV {{ item.level }}</div>
        <div>{{ item.account }}</div>
        <div>{{ item.phone }}</div>
        <div>{{ item.lastName }} {{ item.firstName }}</div>
        <div>{{ item.type === 1 ? "自然人" : "法人" }}</div>
        <div>{{ item.createdAt }}</div>
      </TableItem>
    </TablePaginationContainer>
    <Loader v-else />
    <Dialog v-model="exportDialog" title="匯出報表">
      <form
        @submit.prevent
        class="grid grid-cols-[4rem_1fr] gap-4 items-center justify-end"
        id="export-users-form"
        v-if="exportDialog"
      >
        <FormLabel label="註冊日期" />
        <FormDate
          v-model:start-date="exportDialog.startAt"
          v-model:end-date="exportDialog.endAt"
          required
        />
        <FormLabel label="狀態" />
        <a-checkbox-group
          v-model:value="exportDialog.status"
          :options="[
            { label: '未啟用', value: 0 },
            { label: '已啟用', value: 1 },
            { label: '已停用', value: 2 },
            { label: '凍結中', value: 3 },
          ]"
          required
        />
      </form>
      <template #actions>
        <Button variant="outline" @click="exportDialog = false"> 取消 </Button>
        <Button @click="exportUsersData"> 送出 </Button>
      </template>
    </Dialog>
    <Dialog v-model="legalDialog" title="新增法人" :loading="loading">
      <form
        @submit.prevent
        class="flex flex-col gap-4"
        v-if="legalDialog"
        id="add-juridical-user-form"
      >
        <div class="grid grid-cols-[9rem_1fr] gap-2 items-center justify-end">
          <FormLabel label="E-mail" required />
          <FormInput type="email" v-model="legalDialog.account" required />
          <FormLabel label="名稱" required />
          <FormInput v-model="legalDialog.name" required />
          <FormLabel label="統一編號" />
          <FormInput v-model="legalDialog.nationalID" />
          <FormLabel label="註冊地" />
          <a-select
            v-model:value="legalDialog.country"
            size="large"
            v-if="!countryLoading"
          >
            <a-select-option
              :value="item.code"
              v-for="item of countryList"
              :key="item.code"
              >{{ item.chinese }}</a-select-option
            >
          </a-select>
          <FormLabel label="註冊登記日" />
          <FormInput type="date" v-model="legalDialog.birthDate" />
          <FormLabel label="法人性質" />
          <FormInput v-model="legalDialog.juridicalPersonNature" />
          <FormLabel label="營業地址" />
          <FormInput v-model="legalDialog.address" />
          <FormLabel label="聯繫電話" />
          <FormInput v-model="legalDialog.phone" />
          <FormLabel label="行業別" />
          <a-select
            v-model:value="legalDialog.industrialClassificationsID"
            size="large"
            v-if="!industrialLoading"
          >
            <a-select-option
              :value="item.id"
              v-for="item of industrialList"
              :key="item.id"
              >{{ item.chinese }}</a-select-option
            >
          </a-select>
          <FormLabel label="法幣資金來源" />
          <FormInput v-model="legalDialog.juridicalPersonFundsSources" />
          <FormLabel label="虛擬資產來源" />
          <FormInput
            v-model="legalDialog.juridicalPersonCryptocurrencySources"
          />
          <FormLabel label="被授權人姓名" />
          <FormInput v-model="legalDialog.authorizedPersonName" />
          <FormLabel label="被授權人身份證字號" />
          <FormInput v-model="legalDialog.authorizedPersonNationalID" />
          <FormLabel label="被授權人聯絡電話" />
          <FormInput v-model="legalDialog.authorizedPersonPhone" />
          <FormLabel label="備註" />
          <FormInput v-model="legalDialog.comment" />
        </div>
      </form>
      <template #actions>
        <Button variant="outline" @click="legalDialog = false"> 取消 </Button>
        <Button @click="addJuridicalUser"> 送出 </Button>
      </template>
    </Dialog>
  </div>
</template>
