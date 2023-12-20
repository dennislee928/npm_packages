<script setup>
import useAPI from "@/hooks/useAPI";

const { getUserInviteInfo } = useAPI();
const route = useRoute();

const { uid } = route.params;
const loading = ref(false);
const data = ref(null);
const detailDialog = ref(false);

onMounted(async () => {
  getData();
});
async function getData() {
  loading.value = true;
  let result = await getUserInviteInfo(uid);
  loading.value = false;

  if (result.error) return;
  data.value = result;
}
</script>
<template>
  <TableContainer
    cols="repeat(auto-fit, minmax(4rem, 1fr))"
    v-if="!loading && data"
  >
    <TableHeader>
      <div>邀請人數</div>
      <div>達成人數</div>
      <div>累計獎勵</div>
      <div>尚未提領</div>
      <div>邀請碼</div>
    </TableHeader>
    <TableItem>
      <div>{{ data.totalInvited }}</div>
      <div class="link" @click="detailDialog = true">
        {{ data.totalSucceed }}
      </div>
      <ParsePrice :price="data.totalReward" unit="USDT" :align-right="false" />
      <ParsePrice :price="data.notWithdrew" unit="USDT" :align-right="false" />
      <div>{{ data.inviteCode }}</div>
    </TableItem>
  </TableContainer>
  <Dialog
    v-model="detailDialog"
    :loading="loading"
    title="邀請好友"
    variant="sidebar"
    max-width="1100px"
  >
    <UserInviteDetail v-if="detailDialog" />
    <template #actions>
      <div class="flex-1"></div>
      <Button @click="detailDialog = false"> 關閉 </Button>
    </template>
  </Dialog>
</template>
