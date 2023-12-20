<script setup></script>
<template>
  <div class="m-6 flex gap-2 justify-between items-end">
    <div class="flex gap-2 items-end">
      <FormSelect label="狀態">
        <option value="all">全部</option>
        <option value="in-progress">處理中</option>
        <option value="completed">已完成</option>
        <option value="rejected">已失敗</option>
        <option value="rejected">已取消</option>
      </FormSelect>
      <FormSelect label="類型">
        <option value="all">全部</option>
        <option value="natural">出金</option>
        <option value="legal">入金</option>
      </FormSelect>
      <FormDate label="日期" />
      <FormInput label="關鍵字" placeholder="交易編號" />
      <Button> 重設 </Button>
    </div>
    <div class="flex gap-2 items-end">
      <Button variant="export"> 匯出 </Button>
    </div>
  </div>
  <TablePaginationContainerDemo
    cols="1.5fr 1fr 0.5fr 1fr 0.5fr 1fr"
    :data="Array.from({ length: 100 }).map((_, i) => i + 1)"
  >
    <template #default="{ data }">
      <TableHeader>
        <div>訂單編號</div>
        <div>狀態</div>
        <div>類型</div>
        <div class="text-right">金額</div>
        <div class="text-right">手續費</div>
        <div>交易時間</div>
      </TableHeader>
      <TableItem v-for="item in data" :key="item">
        <div>2023G111{{ item }}</div>
        <div><StatusBadge status="completed">已完成</StatusBadge></div>
        <div>入金</div>
        <ParsePrice :price="Math.floor(Math.random() * 10000)" unit="TWD" />
        <ParsePrice :price="Math.random() > 0.5 ? 0 : 15" unit="TWD" />
        <div>2023/10/10 12:00:00</div>
      </TableItem>
    </template>
  </TablePaginationContainerDemo>
</template>
