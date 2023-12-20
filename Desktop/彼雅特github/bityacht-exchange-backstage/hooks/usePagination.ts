import { ref } from "vue";
export default function usePagination() {
  /**
   * API 回傳的資料
   */
  const data = ref<any[]>([]);
  /**
   * 目前頁數
   * @default 1
   */
  const page = ref(1);
  /**
   * 每頁筆數
   * @default 25
   */
  const perPage = ref(25);
  /**
   * 總頁數
   * @default 1
   */
  const totalPages = ref(1);
  /**
   * 總筆數
   * @default 0
   */
  const totalRecord = ref(0);
  /**
   * API 回傳的資料格式
   */
  type Result = {
    data: any[];
    page: number;
    pageSize: number;
    totalRecord: number;
  };
  /**
   * 使用 API 回傳的資料更新分頁資訊
   * @param result
   */
  function updatePagination(result: Result) {
    data.value = result.data;
    page.value = result.page;
    perPage.value = result.pageSize;
    totalPages.value = Math.ceil(result.totalRecord / result.pageSize);
    totalRecord.value = result.totalRecord;
    // 避免頁數超過總頁數
    if (page.value > totalPages.value) {
      page.value = totalPages.value;
    }
  }
  return {
    data,
    page,
    perPage,
    totalPages,
    totalRecord,
    updatePagination,
  };
}
