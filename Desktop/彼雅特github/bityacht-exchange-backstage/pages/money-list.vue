<script setup>
import { FileTextOutlined } from "@ant-design/icons-vue";
import useUserInfo from "@/hooks/useUserInfo";
const { userInfo } = useUserInfo();
const form = ref({
  status: "all",
  type: "all",
  date: {
    startDate: "",
    endDate: "",
  },
});
const exportDialog = ref(false);
</script>
<template>
  <div>
    <PageTitle> 出/入金列表 </PageTitle>
    <div class="m-6 flex gap-2 justify-between items-end">
      <div class="flex gap-2 items-end">
        <FormSelect label="狀態" v-model="form.status">
          <option value="all">全部</option>
          <option value="in-progress">審核中</option>
          <option value="completed">已通過</option>
          <option value="rejected">已拒絕</option>
        </FormSelect>
        <FormSelect label="類型" v-model="form.type">
          <option value="all">全部</option>
          <option value="natural">自然人</option>
          <option value="legal">法人</option>
        </FormSelect>
        <FormDate
          label="日期"
          v-model:start-date="form.date.startDate"
          v-model:end-date="form.date.endDate"
        />
        <FormInput label="搜尋" />
        <Button> 重設 </Button>
      </div>
      <div class="flex gap-2 items-end">
        <Button
          variant="export"
          @click="exportDialog = true"
          v-show="[1, 2, 4, 5].includes(userInfo.managersRolesID)"
        >
          匯出
        </Button>
      </div>
    </div>
    <TablePaginationContainerDemo
      cols="1fr 0.5fr 0.5fr 1fr 1fr 1fr 1fr 2em"
      :data="Array.from({ length: 100 }).map((_, i) => i + 1)"
    >
      <template #default="{ data }">
        <TableHeader>
          <div>訂單編號</div>
          <div>狀態</div>
          <div>類型</div>
          <div class="text-right">金額</div>
          <div class="text-right">手續費</div>
          <div>UID</div>
          <div>交易時間</div>
          <div>發票</div>
        </TableHeader>
        <TableItem v-for="item in data" :key="item">
          <div>{{ `1023G105493` + item }}</div>
          <div>
            <StatusBadge status="in-progress" variant="filled">
              審核中
            </StatusBadge>
          </div>
          <div>入金</div>
          <ParsePrice :price="1000" unit="TWD" />
          <ParsePrice :price="0" unit="TWD" />
          <div>23394972</div>
          <div>2023-01-01 12:00:00</div>
          <div>
            <FileTextOutlined />
          </div>
        </TableItem>
      </template>
    </TablePaginationContainerDemo>
    <Dialog v-model="exportDialog" title="匯出報表">
      <div class="flex flex-col gap-4">
        <FormDate
          label="註冊日期"
          v-model:start-date="form.date.startDate"
          v-model:end-date="form.date.endDate"
        />
        <FormSelect label="狀態" v-model="form.status">
          <option value="all">全部</option>
          <option value="in-progress">審核中</option>
          <option value="completed">已通過</option>
          <option value="rejected">已拒絕</option>
        </FormSelect>
      </div>
      <template #actions>
        <Button variant="outline" @click="exportDialog = false"> 取消 </Button>
        <Button @click="exportDialog = false"> 確認 </Button>
      </template>
    </Dialog>
  </div>
</template>
