<script setup>
import useAPI from "@/hooks/useAPI";
import useNotification from "@/hooks/useNotification";
import useUserStatus from "@/hooks/useUserStatus";
import useUserInfo from "@/hooks/useUserInfo";
import useUserOptions from "@/hooks/useUserOptions";

const { status, badgeStatus } = useUserStatus();
const { getUserInfo, updateUserStatus, updateUserLevel } = useAPI();
const route = useRoute();
const toast = useNotification();
const { userInfo } = useUserInfo();
const { result: countryList, loading: countryLoading } =
  useUserOptions("countries");
const { result: industrialList, loading: industrialLoading } =
  useUserOptions("ics");

const loading = ref(false);
// Dialog
const statusDialog = ref(false);
const statuslogDialog = ref(false);
const levelDialog = ref(false);

// Data
const memberInfo = ref(null);
const { uid } = route.params;
onMounted(async () => {
  await getData();
});
async function getData() {
  loading.value = true;
  let result = await getUserInfo(uid);
  memberInfo.value = result;
  loading.value = false;
}
async function changeStatus() {
  if (document.getElementById("update-status-form").reportValidity()) {
    loading.value = true;
    let result = await updateUserStatus(
      uid,
      statusDialog.value.status,
      statusDialog.value.comment
    );
    if (!result.error) {
      statusDialog.value = false;
      await getData();
    }
    loading.value = false;
  }
}
async function changeLevel() {
  if (levelDialog.value.level) {
    loading.value = true;
    let result = await updateUserLevel(uid, levelDialog.value.level);
    if (!result.error) {
      levelDialog.value = false;
      await getData();
    }
    loading.value = false;
  } else {
    toast.error("請選擇等級");
  }
}
</script>
<template>
  <div v-if="!loading && memberInfo">
    <TableContainer>
      <TableHeader>
        <div>欄位</div>
        <div class="flex items-center justify-between">
          <span> 資料 </span>
          <Button size="mini" @click="statuslogDialog = true">
            異動紀錄
          </Button>
        </div>
      </TableHeader>
      <TableItem>
        <div>UID</div>
        <div>{{ memberInfo.id }}</div>
      </TableItem>
      <TableItem>
        <div>狀態</div>
        <div class="flex items-center justify-between max-w-[200px]">
          <StatusBadge :status="badgeStatus[memberInfo.status]">
            {{ status[memberInfo.status] }}
          </StatusBadge>
          <Button
            size="mini"
            @click="statusDialog = { status: memberInfo.status, comment: '' }"
            v-show="[4, 2, 1].includes(userInfo.managersRolesID)"
          >
            修改
          </Button>
        </div>
      </TableItem>
      <TableItem>
        <div>類型</div>
        <div>{{ memberInfo.type === 1 ? "自然人" : "法人" }}</div>
      </TableItem>
      <TableItem>
        <div>等級</div>
        <div class="flex items-center justify-between max-w-[200px]">
          <span> LV {{ memberInfo.level }} </span>
          <Button
            size="mini"
            @click="
              levelDialog = {
                level: null,
              }
            "
            v-show="[4, 2, 1].includes(userInfo.managersRolesID)"
          >
            修改
          </Button>
        </div>
      </TableItem>
      <TableItem>
        <div>{{ memberInfo.type === 1 ? "姓名" : "名稱" }}</div>
        <div>{{ memberInfo.lastName }}{{ memberInfo.firstName }}</div>
      </TableItem>
      <TableItem>
        <div>電話</div>
        <div>{{ memberInfo.phone }}</div>
      </TableItem>
      <TableItem>
        <div>E-mail</div>
        <div>{{ memberInfo.account }}</div>
      </TableItem>
      <template v-if="memberInfo.type === 2">
        <TableItem v-if="!countryLoading">
          <div>註冊地</div>
          <div>
            {{
              countryList.filter((x) => x.code === memberInfo.countriesCode)[0]
                ?.chinese
            }}
          </div>
        </TableItem>
        <TableItem>
          <div>註冊登記日</div>
          <div>{{ memberInfo.birthDate }}</div>
        </TableItem>
        <TableItem>
          <div>法人性質</div>
          <div>{{ memberInfo.juridicalPersonNature }}</div>
        </TableItem>
        <TableItem>
          <div>營業地址</div>
          <div>{{ memberInfo.address }}</div>
        </TableItem>
        <TableItem v-if="!industrialLoading">
          <div>行業別</div>
          <div>
            {{
              industrialList.filter(
                (x) => x.id === memberInfo.industrialClassificationsID
              )[0]?.chinese
            }}
          </div>
        </TableItem>
        <TableItem>
          <div>虛擬資產來源</div>
          <div>{{ memberInfo.juridicalPersonCryptocurrencySources }}</div>
        </TableItem>
        <TableItem>
          <div>法幣資金來源</div>
          <div>{{ memberInfo.fundsSources }}</div>
        </TableItem>
        <TableItem>
          <div>被授權人姓名</div>
          <div>{{ memberInfo.authorizedPersonName }}</div>
        </TableItem>
        <TableItem>
          <div>被授權人身份證字號</div>
          <div>{{ memberInfo.authorizedPersonNationalID }}</div>
        </TableItem>
        <TableItem>
          <div>被授權人聯絡電話</div>
          <div>{{ memberInfo.authorizedPersonPhone }}</div>
        </TableItem>
        <TableItem>
          <div>備註</div>
          <div>{{ memberInfo.comment }}</div>
        </TableItem>
      </template>
      <TableItem>
        <div>註冊時間</div>
        <div>{{ memberInfo.createdAt }}</div>
      </TableItem>
    </TableContainer>
    <TableContainer>
      <TableHeader>
        <div>持有資產</div>
        <div>數量</div>
      </TableHeader>
      <TableItem v-for="item of memberInfo.assets" :key="item.currenciesSymbol">
        <div>{{ item.currenciesSymbol }}</div>
        <div>{{ item.freeAmount }}</div>
      </TableItem>
      <TableItem v-if="!memberInfo.assets.length">
        <div class="text-[#C5C6C9]">尚無資料</div>
      </TableItem>
    </TableContainer>
    <TableContainer>
      <TableHeader>
        <div>發票載具</div>
      </TableHeader>
      <TableItem>
        <div v-if="memberInfo.mobileBarcode">
          {{ memberInfo.mobileBarcode }}
        </div>
        <div v-else class="text-[#C5C6C9]">尚未設定</div>
      </TableItem>
    </TableContainer>
    <UserWithdrawalWhitelist />
  </div>
  <Loader v-else />
  <Dialog v-model="statusDialog" :loading="loading" title="修改狀態">
    <form
      @submit.prevent
      class="grid grid-cols-[4rem_1fr] gap-4 items-center justify-end"
      id="update-status-form"
      v-if="statusDialog"
    >
      <FormLabel label="狀態" />
      <div>
        <a-radio-group
          v-model:value="statusDialog.status"
          :options="
            Array.from(Object.entries(status)).map(([key, value]) => ({
              label: value,
              value: Number(key),
            }))
          "
        />
      </div>
      <FormLabel label="緣由備註" />
      <FormInput v-model="statusDialog.comment" required />
    </form>
    <template #actions>
      <Button variant="outline" @click="statusDialog = false"> 取消 </Button>
      <Button @click="changeStatus"> 確定 </Button>
    </template>
  </Dialog>
  <Dialog v-model="levelDialog" :loading="loading" title="修改狀態">
    <form
      @submit.prevent
      class="grid grid-cols-[4rem_1fr] gap-4 items-center justify-end"
      id="update-level-form"
      v-if="levelDialog"
    >
      <FormLabel label="會員等級" />
      <div>
        <a-radio-group
          v-model:value="levelDialog.level"
          :options="
            (memberInfo.type === 1 ? [2, 3, 4, 5] : [1, 2]).map((x) => ({
              label: 'LV ' + x,
              value: x,
            }))
          "
          required
        />
      </div>
    </form>
    <template #actions>
      <Button variant="outline" @click="levelDialog = false"> 取消 </Button>
      <Button @click="changeLevel"> 確定 </Button>
    </template>
  </Dialog>
  <Dialog
    v-model="statuslogDialog"
    title="異動紀錄"
    max-width="1100px"
    variant="sidebar"
  >
    <UserStatuslog v-if="statuslogDialog" />
    <template #actions>
      <Button @click="statuslogDialog = false"> 關閉 </Button>
    </template>
  </Dialog>
</template>
