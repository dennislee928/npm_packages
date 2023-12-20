export default function useAPI() {
  // API Setting
  const config = useRuntimeConfig();
  const baseURL = config.public.apiURL;
  /**
   * @param isFormData 是否為 FormData
   * @param parseResultAsArray 是否解析為 Array
   * @param parseResultAsBuffer 是否解析為 Buffer
   */
  type RequestOptions = {
    isFormData?: Boolean;
    parseResultAsArray?: Boolean;
    parseResultAsBuffer?: Boolean;
  };

  const request = async (
    url: string,
    method: string,
    data: any,
    options: RequestOptions = {
      isFormData: false,
      parseResultAsArray: false,
      parseResultAsBuffer: false,
    }
  ) => {
    let query = "";
    let isError = false;
    if (method === "GET" && Object.keys(data).length) {
      query =
        "?" +
        Object.keys(data)
          .map((key) => {
            if (Array.isArray(data[key])) {
              let queryString = "";
              for (const item of data[key]) {
                queryString += `${encodeURIComponent(key)}=${encodeURIComponent(
                  item
                )}`;
              }
              return queryString;
            }
            if (["startAt", "endAt"].includes(key)) {
              return `${encodeURIComponent(key)}=${encodeURIComponent(
                `"${data[key]}"`
              )}`;
            }
            return `${encodeURIComponent(key)}=${encodeURIComponent(
              data[key]
            )}`;
          })
          .join("&");
    }
    let headers: Record<string, string> = {};
    if (!options.isFormData) headers["Content-Type"] = "application/json";
    headers["Authorization"] = `Bearer ${
      localStorage.getItem("accessToken") || ""
    }`;
    let result = await fetch(`${baseURL}${url}${query}`, {
      method,
      headers,
      body:
        method === "GET"
          ? null
          : options.isFormData
          ? data
          : JSON.stringify(data),
    }).catch((e) => {
      isError = true;
      return e;
    });

    // parse as json
    try {
      if (result.status >= 400) isError = true;
      if (options.parseResultAsBuffer) {
        result = await result.blob();
        return result;
      }
      result = await result.json();
      if (result.code) isError = true;
      // if []
      if (options.parseResultAsArray) return result;
      return { ...result, error: isError };
    } catch (e) {}

    // return {}
    if (isError) {
      return { ...result, error: true };
    }
    return {};
  };
  const api = {
    /**
     * 建立 GET 請求
     * @param url API URL
     * @param data 查詢參數
     * @param options 請求選項
     * @returns
     */
    get: async (url: string, data?: any, options?: RequestOptions) =>
      request(url, "GET", data, options),
    /**
     * 建立 POST 請求
     * @param url API URL
     * @param data body
     * @param options 請求選項
     * @returns
     */
    post: async (url: string, data: any, options?: RequestOptions) =>
      request(url, "POST", data, options),
    /**
     * 建立 DELETE 請求
     * @param url API URL
     * @returns
     */
    delete: async (url: string) => request(url, "DELETE", {}),
    /**
     * 建立 PUT 請求
     * @param url API URL
     * @param data body
     * @returns
     */
    put: async (url: string, data: any, options?: RequestOptions) =>
      request(url, "PUT", data, options),
    /**
     * 建立 PATCH 請求
     * @param url API URL
     * @param data body
     * @param options 請求選項
     * @returns
     */
    patch: async (url: string, data: any, options?: RequestOptions) =>
      request(url, "PATCH", data, options),
  };
  /**
   * 產生用於匯出的 URL
   * @param url 網址
   * @param query 查詢參數（Object）
   * @returns URL
   */
  function getUrl(url: string, query?: any) {
    let queryString = `?accessToken=${localStorage.getItem("accessToken")}`;
    for (const key in query) {
      if (Array.isArray(query[key])) {
        for (const item of query[key]) {
          queryString += `&${encodeURIComponent(key)}=${encodeURIComponent(
            item
          )}`;
        }
        continue;
      }
      if (["startAt", "endAt"].includes(key)) {
        queryString += `&${encodeURIComponent(key)}=${encodeURIComponent(
          `"${query[key]}"`
        )}`;
        continue;
      }
      queryString += `&${encodeURIComponent(key)}=${encodeURIComponent(
        query[key]
      )}`;
    }

    return `${baseURL}${url}${queryString}`;
  }
  // API
  /**
   * 登入
   * @param account 帳號
   * @param password 密碼
   */
  const login = async (body: { account: string; password: string }) =>
    api.post("/admin/login", body);
  /**
   * 忘記密碼
   * @param account 帳號
   */
  const forgotPassword = async (account: string) =>
    api.post("/admin/forgot-password", { account });
  /**
   * 登出
   */
  const logout = async () => api.post("/admin/logout", {});

  /**
   * 修改密碼
   * @param password 新的密碼
   */
  const changePassword = async (password: string) =>
    api.patch("/admin/password", { password });

  // Admin management
  /**
   * 取得管理員列表
   * @param page 頁數
   * @param pageSize 每頁筆數
   */
  const getAdmins = async (
    page: string | number = 1,
    pageSize: string | number = 25
  ) => api.get("/admin/settings/managers", { page, pageSize });
  /**
   * 新增管理員
   * @param account 帳號
   * @param managersRolesID 角色 ID
   * - `1` 超級帳號
   * - `2` 法遵
   * - `3` 客服
   * - `4` 風控
   * - `5` 財務
   * @param name 名稱
   */
  const createAdmin = async (body: {
    account: string;
    managersRolesID: number;
    name: string;
  }) => api.post("/admin/settings/managers", body);
  /**
   * 更新管理員
   * @param id 管理員 ID
   * @param account 帳號
   * @param managersRolesID 角色 ID
   * - `1` 超級帳號
   * - `2` 法遵
   * - `3` 客服
   * - `4` 風控
   * - `5` 財務
   * @param name 名稱
   * @param password 密碼
   * @param status 狀態
   * - `0` 停用
   * - `1` 啟用
   */
  const updateAdmin = async ({
    id,
    account,
    managersRolesID,
    name,
    password,
    status,
  }: {
    id: number;
    account?: string;
    managersRolesID?: number;
    name?: string;
    password?: string;
    status?: number;
  }) =>
    api.patch(`/admin/settings/managers/${id}`, {
      account,
      managersRolesID,
      name,
      password,
      status,
    });
  /**
   * 刪除管理員
   * @param id 管理員 ID
   */
  const deleteAdmin = async (id: number) =>
    api.delete(`/admin/settings/managers/${id}`);

  // transactionpairs api
  /**
   * 取得交易對列表
   * @param page 頁數
   * @param pageSize 每頁筆數
   */
  const getTransactionpairs = async (
    page: string | number = 1,
    pageSize: string | number = 25
  ) => api.get("/admin/settings/transactionpairs", { page, pageSize });
  /**
   * 更新交易對
   * @param baseCurrenciesSymbol 基礎幣別 `USDT`/BTC
   * @param quoteCurrenciesSymbol 報價幣別 USDT/`BTC`
   * @param handlingChargeRate 手續費率
   * @param spreadsOfBuy 買價差
   * @param spreadsOfsell 賣價差
   * @param status 狀態
   * - `0` 停用
   * - `1` 啟用
   */
  const updateTransactionpairs = async (body: {
    baseCurrenciesSymbol: string;
    handlingChargeRate: string;
    quoteCurrenciesSymbol: string;
    spreadsOfBuy: string;
    spreadsOfsell: string;
    status: number;
  }) => api.patch(`/admin/settings/transactionpairs`, body);

  // Banner management
  /**
   * 取得橫幅列表
   * @param page 頁數
   * @param pageSize 每頁筆數
   * @param options.status 狀態
   * - `-1` 全部
   * - `0` 停用
   * - `1` 啟用
   */
  const getBanners = async (
    page: string | number = 1,
    pageSize: string | number = 25,
    options: {
      status: "-1" | "0" | "1"; // -1: all, 0: disabled, 1: enabled
    }
  ) =>
    api
      .get("/admin/settings/banners", { page, pageSize, ...options })
      .then((x) => {
        x.data = x.data
          .map((x: any) => {
            x.appImage = `${baseURL}/assets/banners/${x.appImage}`;
            x.webImage = `${baseURL}/assets/banners/${x.webImage}`;
            x.isPermanent = x.startAt === "" && x.endAt === "";
            return x;
          })
          .sort((a: any, b: any) => a.priority - b.priority);
        return x;
      });
  /**
   * 新增橫幅
   * @param formData 橫幅資料
   * @param formData.appImage app 圖片
   * @param formData.webImage web 圖片
   * @param formData.button 按鈕文字
   * @param formData.buttonUrl 按鈕連結
   * @param formData.priority 順序
   * @param formData.startAt 開始時間
   * @param formData.endAt 結束時間
   * @param formData.status 狀態
   * - `0` 停用
   * - `1` 啟用
   * @param formData.title 標題
   * @param formData.subTitle 副標題
   */
  const createBanner = async (formData: FormData) =>
    api.post("/admin/settings/banners", formData, { isFormData: true });
  /**
   * 更新橫幅
   * @param id 橫幅 ID
   * @param formData 橫幅資料
   * @param formData.appImage app 圖片
   * @param formData.webImage web 圖片
   * @param formData.button 按鈕文字
   * @param formData.buttonUrl 按鈕連結
   * @param formData.priority 順序
   * @param formData.startAt 開始時間
   * @param formData.endAt 結束時間
   * @param formData.status 狀態
   * - `0` 停用
   * - `1` 啟用
   * @param formData.title 標題
   * @param formData.subTitle 副標題
   */
  const updateBanner = async (id: string, formData: FormData) =>
    api.patch(`/admin/settings/banners/${id}`, formData, { isFormData: true });
  /**
   * 刪除橫幅
   * @param id 橫幅 ID
   */
  const deleteBanner = async (id: string) =>
    api.delete(`/admin/settings/banners/${id}`);
  /**
   * 更新橫幅順序
   * @param data 橫幅順序資料
   * ```json
   * {
   *    "rows": [
   *      { "id": 0, "priority": 0  }
   *    ]
   * }
   * ```
   */
  const updateBannerPriority = async (data: any) =>
    api.post(`/admin/settings/banners/priority`, data);
  // Mainnet management
  /**
   * 取得主網列表
   * @returns
   */
  const getMainnets = async () => api.get(`/admin/settings/mainnets`, {});
  /**
   * 編輯主網
   * @param mainnet 主網 ID
   * @param currency 幣種 ID
   * @param withdrawFee 提領手續費 （單位 = Currency）
   * @param withdrawMin 提領最小數量 （單位 = Currency）
   * @returns
   */
  const updateMainnets = async ({
    mainnet,
    currency,
    withdrawFee,
    withdrawMin,
  }: {
    mainnet: number;
    currency: number;
    withdrawFee: number;
    withdrawMin: number;
  }) =>
    api.patch(`/admin/settings/mainnets/${currency}/${mainnet}`, {
      withdrawFee,
      withdrawMin,
    });
  // users management
  /**
   * 取得使用者列表
   * @param page 頁數
   * @param pageSize 每頁筆數
   * @param options
   * @param options.status 狀態
   * - `0` 未啟用
   * - `1` 已啟用
   * - `2` 已停用
   * - `3` 凍結中
   * @param options.type 類型
   * - `1` 自然人
   * - `2` 法人
   * @param options.startAt 開始時間
   * @param options.endAt 結束時間
   * @param options.search 搜尋
   */
  const getUsers = async (
    page: string | number = 1,
    pageSize: string | number = 25,
    options: {
      status: string | number | null;
      type: string | number | null;
      startAt: string | null;
      endAt: string | null;
      search: string | null;
    }
  ) =>
    api.get("/admin/users", {
      page,
      pageSize,
      ...options,
    });
  /**
   * 建立法人帳號
   * @param account	E-mail
   * @param address	營業地址
   * @param authorizedPersonName	被授權人姓名
   * @param authorizedPersonNationalID	被授權人身分證字號
   * @param authorizedPersonPhone	被授權人聯絡電話
   * @param birthDate	註冊登記日
   * @param comment	其他資訊（備註）
   * @param country	註冊地
   * @param industrialClassificationsID	行業別
   * @param juridicalPersonCryptocurrencySources	虛擬資產來源
   * @param juridicalPersonFundsSources	法幣資金來源
   * @param juridicalPersonNature	法人性質
   * @param name	名稱
   * @param nationalID	統一編號
   * @param phone	聯絡電話
   */
  const createJuridicalUser = async (body: {
    account: string;
    address?: string;
    authorizedPersonName?: string;
    authorizedPersonNationalID?: string;
    authorizedPersonPhone?: string;
    birthDate?: string;
    comment?: string;
    country?: string;
    industrialClassificationsID?: number;
    juridicalPersonCryptocurrencySources?: string;
    juridicalPersonFundsSources?: string;
    juridicalPersonNature?: string;
    name: string;
    nationalID?: string;
    phone?: string;
  }) => api.post("/admin/users", body);
  /**
   * 匯出使用者
   * @param startAt 開始時間
   * @param endAt 結束時間
   * @param statusList 狀態
   * - `0` 未啟用
   * - `1` 已啟用
   * - `2` 已停用
   * - `3` 凍結中
   */
  const getExportUsersUrl = (body: {
    endAt: string;
    startAt: string;
    statusList: number[];
  }) => getUrl("/admin/users/export", body);
  /**
   * 取得使用者資訊
   * @param id 使用者 ID
   */
  const getUserInfo = async (id: string) => api.get(`/admin/users/${id}`, {});
  /**
   * 取得使用者 KYC 資訊
   * @param id 使用者 ID
   */
  const getUserKycInfo = async (id: string) =>
    api.get(`/admin/users/${id}/kycs`, {});
  /**
   * 更新使用者等級
   * @param id 使用者 ID
   * @param level 等級
   * - `3` LV 3
   * - `4` LV 4
   * - `5` LV 5
   */
  const updateUserLevel = async (id: string, level: number) =>
    api.patch(`/admin/users/${id}/level`, { level });
  /**
   * 取得該用戶資料狀態異動紀錄
   * @param id 使用者 ID
   * @param page 頁碼
   * @param pageSize 每頁筆數
   */
  const getUserStatuslog = async (id: string, page: number, pageSize: number) =>
    api.get(`/admin/users/${id}/statuslogs`, { page, pageSize });
  /**
   * 更新使用者狀態
   * @param id 使用者 ID
   * @param status 狀態
   * - `0` 未啟用
   * - `1` 已啟用
   * - `2` 已停用
   * - `3` 凍結中
   * @param comment 備註
   */
  const updateUserStatus = async (
    id: string,
    status: number,
    comment: string
  ) => api.patch(`/admin/users/${id}/status`, { status, comment });
  /**
   * 匯出使用者狀態異動紀錄
   * @param id 使用者 ID
   */
  const getExportUserStatuslogUrl = (id: string) =>
    getUrl(`/admin/users/${id}/statuslogs/export`);
  /**
   * 取得用戶提幣白名單
   * @param id 使用者 ID
   * @param page 頁碼
   * @param pageSize 每頁筆數
   * @returns
   */
  const getWithdrawalWhitelist = async (
    id: string,
    page: number,
    pageSize: number
  ) => api.get(`/admin/users/${id}/withdrawal-whitelist`, { page, pageSize });
  /**
   * 匯出用戶提幣白名單
   * @param id 使用者 ID
   * @returns
   */
  const getExportWithdrawalWhitelistUrl = (id: string) =>
    getUrl(`/admin/users/${id}/withdrawal-whitelist/export`);
  /**
   * 刪除用戶提幣白名單
   * @param id 使用者 ID
   * @param whitelistID 提幣白名單 ID
   * @returns
   */
  const deleteWithdrawalWhitelist = async (id: string, whitelistID: string) =>
    api.delete(`/admin/users/${id}/withdrawal-whitelist/${whitelistID}`);
  // user login history
  /**
   * 取得該用戶登入紀錄
   * @param id 使用者 ID
   * @param page 頁碼
   * @param pageSize 每頁筆數
   */
  const getUserLoginLog = async (id: string, page: number, pageSize: number) =>
    api.get(`/admin/users/${id}/login-logs`, { page, pageSize });
  /**
   * 匯出使用者登入紀錄
   * @param id 使用者 ID
   */
  const getExportUserLoginLogUrl = (id: string) =>
    getUrl(`/admin/users/${id}/login-logs/export`);
  /**
   * 取得該用戶邀請紀錄
   * @param id 使用者 ID
   * @param page 頁碼
   * @param pageSize 每頁筆數
   */
  const getUserInviteRecords = async (
    id: string,
    page: number,
    pageSize: number
  ) => api.get(`/admin/users/${id}/invite`, { page, pageSize });
  /**
   * 取得該用戶邀請資訊
   * @param id 使用者 ID
   */
  const getUserInviteInfo = async (id: string) =>
    api.get(`/admin/users/${id}/invite-info`, {});
  /**
   * 取得該用戶獎勵清單
   * @param id 使用者 ID
   * @param page 頁碼
   * @param pageSize 每頁筆數
   * @param options
   * @param options.action 類型
   * - `0` 全部
   * - `1` 返佣
   * - `2` 提領
   */
  const getUserInviteRewards = async (
    id: string,
    page: number,
    pageSize: number,
    options: {
      action: 0 | 1 | 2;
    }
  ) =>
    api.get(`/admin/users/${id}/invite-rewards`, {
      page,
      pageSize,
      ...options,
    });
  /**
   * 匯出該用戶獎勵清單
   * @param id 使用者 ID
   * @param action 類型
   * - `0` 全部
   * - `1` 返佣
   * - `2` 提領
   * @returns URL
   */
  const getExportUserInviteRewardsUrl = (id: string, action: 0 | 1 | 2) =>
    getUrl(`/admin/users/${id}/invite-rewards/export`, { action });
  /**
   * 檢查登入狀態
   * @returns
   */
  const checkAuth = async () => api.get("/admin/auth", {});
  // KYC
  /**
   * 取得 KYC 列表
   * @param page 頁碼
   * @param pageSize 每頁筆數
   * @param options 查詢條件
   * @param options.finalReview 最終審核
   * - `0` 全部
   * - `1` 待審核
   * - `2` 已通過
   * - `3` 已拒絕
   * @param options.complianceReview 法遵審核
   * - `0` 全部
   * - `1` 待審核
   * - `2` 已通過
   * - `3` 已拒絕
   * @param options.startAt 開始時間
   * @param options.endAt 結束時間
   */
  const getKycList = async (
    page: number,
    pageSize: number,
    options: {
      finalReview: number;
      complianceReview: number;
      startAt: string;
      endAt: string;
      search: string;
    }
  ) =>
    api.get("/admin/kycs", {
      page,
      pageSize,
      ...options,
    });
  /**
   * 匯出 KYC 資訊
   * @param startAt 開始時間
   * @param endAt 結束時間
   * @param statusList 狀態
   * - `0` 全部
   * - `1` 待審核
   * - `2` 已通過
   * - `3` 已拒絕
   */
  const getExportKycListUrl = (query: {
    startAt: string;
    endAt: string;
    statusList: number[];
  }) => getUrl("/admin/kycs/export", query);
  /**
   * 更新法遵審核
   * @param id user id
   * @param result 審核結果
   * - `2` 已通過
   * - `3` 已拒絕
   * @param comment 備註
   */
  const updateUserKycComplianceReview = async (
    id: string,
    result: number,
    comment: string
  ) =>
    api.patch(`/admin/users/${id}/kycs/compliance-review`, { result, comment });
  /**
   * 更新最終審核
   * @param id user id
   * @param result 審核結果
   * - `2` 已通過
   * - `3` 已拒絕
   * @param comment 備註
   */
  const updateUserKycFinalReview = async (
    id: string,
    result: number,
    comment: string
  ) =>
    api.patch(`/admin/users/${id}/kycs/final-review`, {
      result,
      comment,
    });
  /**
   * 更新認證結果(外籍人士)
   * @param id user id
   * @param auditStatus 審核結果
   * - `2` 通過
   * - `3` 不通過
   * @param comment 備註
   */
  const updateUserKycIdvAuditStatus = async (
    id: string,
    auditStatus: number,
    comment: string
  ) =>
    api.patch(`/admin/users/${id}/kycs/idv-audit-status`, {
      auditStatus,
      comment,
    });
  /**
   * 更新 Krypto 審核
   * @param id user id
   * @param result 審核結果
   * - `2` 已通過
   * - `3` 已拒絕
   * @param comment 備註
   */
  const updateUserKycKryptoReview = async (
    id: string,
    result: number,
    comment: string
  ) => api.patch(`/admin/users/${id}/kycs/krypto-review`, { result, comment });
  /**
   * 重送 Krypto 審核
   * @param id user id
   */
  const resendUserKycKryptoReview = async (id: string) =>
    api.post(`/admin/users/${id}/kycs/resent-krypto`, {});
  /**
   * 編輯 Krypto 單號
   * @param id user id
   * @param taskID Krypto 單號
   */
  const updateUserKycKryptoTaskID = async (id: string, taskID: string) =>
    api.patch(`/admin/users/${id}/kycs/task-id`, { taskID });
  /**
   * 取得 KYC 審核紀錄
   * @param id 使用者 ID
   * @param page 頁面
   * @param pageSize 每頁筆數
   */
  const getUserKycReviewslogs = async (
    id: string,
    page: number,
    pageSize: number
  ) => api.get(`/admin/users/${id}/kycs/reviewslogs`, { page, pageSize });
  /**
   * 匯出 KYC 審核紀錄
   * @param id 使用者 ID
   */
  const getExportUserKycReviewslogsUrl = (id: string) =>
    getUrl(`/admin/users/${id}/kycs/reviewslogs/export`);
  // KYC Name check
  /**
   * 更新姓名審核
   * @param id
   * @param formData.file PDF
   * @param formData.result 審核結果
   * - `2` 已通過
   * - `3` 已拒絕
   */
  const updateKycNameCheck = async (id: string, formData: FormData) =>
    api.patch(`/admin/users/${id}/kycs/name-check`, formData, {
      isFormData: true,
    });
  /**
   * 取得姓名審核 PDF
   * @param id
   * @returns PDF Buffer
   */
  const downloadKycNameCheck = async (id: string) =>
    api.get(
      `/admin/users/${id}/kycs/name-check/pdf`,
      {},
      { parseResultAsBuffer: true }
    );
  // KYC Risks
  /**
   * 取得 KYC 風險列表
   * @returns KYC 風險列表
   */
  const getKycRisks = async () => api.get(`/admin/kycs/risks`, {});
  /**
   * 建立 KYC 風險項目
   * @param detail 風險描述
   * @param score 風險分數
   * @param factor 風險分類
   * @param subFactor 風險子分類
   */
  const createKycRisk = async (body: {
    detail: string;
    factor: string;
    score: number;
    subFactor: string;
  }) => api.post(`/admin/kycs/risks`, body);
  /**
   * 更新 KYC 風險項目
   * @param id 風險項目 ID
   * @param detail 風險描述
   * @param score 風險分數
   * @param factor 風險分類
   * @param subFactor 風險子分類
   */
  const updateKycRisk = async ({
    id,
    detail,
    factor,
    score,
    subFactor,
  }: {
    id: number;
    detail: string;
    factor: string;
    score: number;
    subFactor: string;
  }) =>
    api.patch(`/admin/kycs/risks/${id}`, {
      detail,
      factor,
      score,
      subFactor,
    });
  /**
   * 刪除 KYC 風險項目
   * @param id 風險項目 ID
   */
  const deleteKycRisk = async (id: number) =>
    api.delete(`/admin/kycs/risks/${id}`);
  // User KYC Risks
  /**
   * 取得使用者 KYC 風險列表
   * @param id 使用者 ID
   */
  const getUserKycRisks = async (id: string) =>
    api.get(`/admin/users/${id}/kycs/risks`, {}, { parseResultAsArray: true });
  /**
   * 更新使用者 KYC 風險列表
   * @param id 使用者 ID
   * @param risksIDs 風險項目 ID 列表
   */
  const updateUserKycRisks = async (id: string, risksIDs: number[]) =>
    api.post(`/admin/users/${id}/kycs/risks`, { risksIDs });
  /**
   * 取得使用者銀行帳戶資訊
   * @param id 使用者 ID
   * @returns
   * {
   *   "account": "string",
   *   "auditTime": "string",
   *   "bankInfo": {
   *     "chinese": "string",
   *     "code": "string",
   *     "english": "string"
   *   },
   *   "banksCode": "string",
   *   "branchInfo": {
   *     "chinese": "string",
   *     "code": "string",
   *     "english": "string"
   *   },
   *   "branchsCode": "string",
   *   "coverImage": "string",
   *   "createdAt": "string",
   *   "id": 0,
   *   "name": "string",
   *   "status": 0
   * }
   */
  const getUserBank = async (id: string) =>
    api.get(`/admin/users/${id}/bank`, {});
  /**
   * 更新使用者銀行帳戶資訊
   * @param id 使用者 ID
   * @param body
   * @param body.comment 備註
   * @param body.id 帳戶 ID
   * @param body.status 狀態
   * - `1`: 審核中
   * - `2`: 已綁定
   * - `3`: 未通過
   * @returns
   */
  const updateUserBank = async (
    id: string,
    body: {
      comment: string;
      id: Number;
      status: Number;
    }
  ) => api.patch(`/admin/users/${id}/bank`, body);
  /**
   * 取得使用者銀行帳戶異動紀錄
   * @param id 使用者 ID
   * @param page 頁碼
   * @param pageSize 每頁筆數
   * @returns
   */
  const getUserBankLogs = async (id: string, page: number, pageSize: number) =>
    api.get(`/admin/users/${id}/bank/logs`, { page, pageSize });
  /**
   * 匯出使用者銀行帳戶資訊
   * @param id 使用者 ID
   * @returns URL
   */
  const getExportUserBankLogsUrl = (id: string) =>
    getUrl(`/admin/users/${id}/bank/logs/export`);
  /**
   * 取得年度審查列表
   * @param page 頁數
   * @param pageSize 每頁筆數
   * @param options
   * @param options.complianceReview 法遵審核
   * `0` 全部
   * `1` 待審核
   * `2` 已通過
   * `3` 已拒絕
   * @param options.startAt 開始時間
   * @param options.endAt 結束時間
   * @param options.search 關鍵字
   */
  const getAnnualKYCReviews = async (
    page: string | number = 1,
    pageSize: string | number = 25,
    options: {
      complianceReview: string | number | null;
      startAt: string | null;
      endAt: string | null;
      search: string | null;
    }
  ) =>
    api.get("/admin/kycs/annual", {
      page,
      pageSize,
      ...options,
    });
  /**
   * 匯出年度審查列表
   * @param startAt 開始時間
   * @param endAt 結束時間
   */
  const getExportAnnualKYCReviewsUrl = ({
    startAt,
    endAt,
  }: {
    startAt: string;
    endAt: string;
  }) => getUrl("/admin/kycs/annual/export", { startAt, endAt });
  // Spots API
  /**
   * 取得提入幣列表
   * @param id `0` 為全部
   * @param page 頁數
   * @param pageSize 每頁筆數
   * @param options.status 狀態
   * - `0` 全部
   * - `1` 待審核
   * - `2` 已通過
   * - `3` 已拒絕
   * @param options.coin 幣別
   * @param options.startAt 開始時間
   * @param options.endAt 結束時間
   * @param options.search 搜尋
   * @returns
   */
  const getSpots = async (
    id: string,
    page: string | number = 1,
    pageSize: string | number = 25,
    options: {
      status: string | number | null;
      coin: string | number | null;
      startAt: string | null;
      endAt: string | null;
      search: string | null;
    }
  ) =>
    api.get("/admin/spots", {
      usersID: id,
      page,
      pageSize,
      ...options,
    });
  /**
   * 匯出提入幣列表
   * @param startAt 開始時間
   * @param endAt 結束時間
   * @param statusList 狀態列表
   * - `0` 全部
   * - `1` 待審核
   * - `2` 已通過
   * - `3` 已拒絕
   * @param id `0` 為全部
   */
  const getExportSpotsUrl = (query: {
    startAt: string;
    endAt: string;
    statusList: number[];
    id: string;
  }) => getUrl("/admin/spots/export", query);
  /**
   * 提入幣 Aegis 匯出
   * @param startAt 開始時間
   * @param endAt 結束時間
   * @param mainnet 主網
   * - `1` BTC
   * - `2` ETH
   * - `3` ERC20
   * - `4` TRC20
   * @param usersID `0` 為全部
   */
  const getExportSpotsAegisUrl = (query: {
    startAt: string;
    endAt: string;
    mainnet: number;
    usersID: string;
  }) => getUrl("/admin/spots/aegis-export", query);

  /**
   * 提入幣 Aegis 匯入
   * @param formData
   * @param formData.file CSV
   */
  const importSpotsAegis = async (formData: FormData) =>
    api.post("/admin/spots/aegis-import", formData, { isFormData: true });
  // transactions
  /**
   * 取得交易訂單列表
   * @param id user id, `0` for all
   * @param page page number
   * @param pageSize page size
   * @param options.status `0` for all, `1` for filled, `2` for killed
   * @param options.side `0` for all, `1` for buy, `2` for sell
   * @param options.startAt `0001/01/01 00:00:00`
   * @param options.endAt `0001/01/01 00:00:00`
   * @param options.search order id or uid
   */
  const getTransactions = async (
    id: string,
    page: string | number = 1,
    pageSize: string | number = 25,
    options: {
      status: string | number | null;
      side: string | number | null;
      startAt: string | null;
      endAt: string | null;
      search: string | null;
    }
  ) =>
    api.get("/admin/transactions", {
      usersID: id,
      page,
      pageSize,
      ...options,
    });
  /**
   * 匯出交易訂單列表
   * @param startAt `0001/01/01 00:00:00`
   * @param endAt `0001/01/01 00:00:00`
   * @param statusList `0` for all, `1` for filled, `2` for killed
   * @param id user id, `0` for all
   */
  const getExportTransactionsUrl = (query: {
    startAt: string;
    endAt: string;
    statusList: number[];
    usersID: string;
  }) => getUrl("/admin/transactions/export", query);
  // Receipts API
  /**
   * 取得所有發票列表
   * @param page page number
   * @param pageSize page size
   * @param options.status `0` for all, `1` for pending, `2` for issuing, `3` for issued, `4` for failed
   * @param options.startAt `0001/01/01 00:00:00`
   * @param options.endAt `0001/01/01 00:00:00`
   * @param options.search receipt id or uid
   */
  const getReceipts = async (
    page: string | number = 1,
    pageSize: string | number = 25,
    options: {
      status: string | number | null;
      startAt: string | null;
      endAt: string | null;
      search: string | null;
    }
  ) =>
    api.get("/admin/receipts", {
      page,
      pageSize,
      ...options,
    });
  /**
   * 取得單一發票資訊
   * @param id receipt id
   */
  const getReceipt = async (id: string) => api.get(`/admin/receipts/${id}`, {});
  // post /admin/receipts/export
  /**
   * 匯出發票列表
   * @param startAt `0001/01/01 00:00:00`
   * @param endAt `0001/01/01 00:00:00`
   * @param statuses `0` for all, `1` for pending, `2` for issuing, `3` for issued, `4` for failed
   */
  const getExportReceiptsUrl = (query: {
    startAt: string;
    endAt: string;
    statuses: number[];
  }) => getUrl("/admin/receipts/export", query);
  /**
   * 發票開立
   * @param ids receipt ids
   */
  const issueReceipts = async (ids: string[]) =>
    api.post("/admin/receipts/issue", { ids });

  // User options
  /**
   * 取得使用者相關選項
   */
  const getUserOptions = async () => api.get("/admin/users/options", {});
  /**
   * 可疑交易管理
   */
  /**
   * 取得可疑交易列表
   * @param page
   * @param pageSize
   * @param options
   * @param options.startAt 掃描開始時間
   * @param options.endAt 掃描結束時間
   * @param options.search
   * @param options.dedicatedReviewResult 狀態 `0` 全部 | `1` 待審核 | `2` 通過 | `3` 駁回
   * @param options.type 類型 `0` 全部 | `1` 多筆提領態樣 | `2` 多筆同額態樣 | `3` 同一出金地址態樣 | `4` 迅速轉出態樣 | `5` 迅速買賣態樣 | `6` 小額接收大額轉出態樣
   * @returns
   */
  const getSuspiciousTransactions = async (
    page: string | number = 1,
    pageSize: string | number = 25,
    options: {
      startAt: string | null;
      endAt: string | null;
      search: string | null;
      dedicatedReviewResult: 0 | 1 | 2 | 3 | null;
      type: 0 | 1 | 2 | 3 | 4 | 5 | 6 | null;
    }
  ) =>
    api.get("/admin/suspicious-txs", {
      page,
      pageSize,
      ...options,
    });
  /**
   * 匯出可疑交易列表
   * @param startAt `0001/01/01 00:00:00`
   * @param endAt `0001/01/01 00:00:00`
   * @param dedicatedReviewResult 狀態 `0` 全部 | `1` 待審核 | `2` 通過 | `3` 駁回
   * @param type 類型 `0` 全部 | `1` 多筆提領態樣 | `2` 多筆同額態樣 | `3` 同一出金地址態樣 | `4` 迅速轉出態樣 | `5` 迅速買賣態樣 | `6` 小額接收大額轉出態樣
   */
  const getExportSuspiciousTransactionsUrl = (query: {
    startAt: string;
    endAt: string;
    dedicatedReviewResult: number[];
    type: number[];
  }) => getUrl("/admin/suspicious-txs/export-csv", query);
  /**
   * 取得可疑交易詳細資訊
   * @param id Suspicious Transactions ID
   */
  const getSuspiciousTransactionDetail = async (id: string) =>
    api.get(`/admin/suspicious-txs/${id}`, {});
  /**
   * 更新可疑交易詳細資訊
   * @param id Suspicious Transactions ID
   * @param body
   * @param body.comment 備註
   * @param body.dedicatedReviewResult 狀態 `0` 全部 | `1` 待審核 | `2` 通過 | `3` 駁回
   * @param body.reportMJIBAt 呈報調查局日期
   * @param body.riskReviewResult 風控審查結果 `1` 資訊審核 (所需欄位：備註) | `2` 風控審查結果 (所需欄位：風控審查結果) | `3` 專責審查 (所需欄位：專責審查結果、備註) | `4` 呈報調查局日期 (所需欄位：呈報調查局日期)
   */
  const updateSuspiciousTransactionDetail = async (
    id: string,
    body: {
      comment: string;
      dedicatedReviewResult: 0 | 1 | 2 | 3;
      reportMJIBAt: string;
      riskReviewResult: 1 | 2 | 3 | 4;
    }
  ) => api.patch(`/admin/suspicious-txs/${id}`, body);
  /**
   * 下載可疑交易相關檔案
   * @param data
   * @param data.id Suspicious Transactions ID
   * @param data.fileType 檔案類型
   * - `1` 資訊審核
   * - `2` 風控審查
   * @param data.filename 檔案名稱
   */
  const downloadSuspiciousTransactionFiles = async (data: {
    id: string;
    fileType: string;
    filename: string;
  }) =>
    api.get(
      `/admin/suspicious-txs/${data.id}/file`,
      {
        fileType: data.fileType,
        filename: data.filename,
      },
      { parseResultAsBuffer: true }
    );
  /**
   * 上傳可疑交易相關檔案
   * @param data FormData
   * @param data.id Suspicious Transactions ID
   * @param data.uploadType 檔案類型
   * - `1` 資訊審核
   * - `2` 風控審查
   * @param data.file 檔案
   */
  const uploadSuspiciousTransactionFiles = async (data: FormData) =>
    api.post(`/admin/suspicious-txs/${data.get("id")}/file`, data, {
      isFormData: true,
    });

  /**
   * 刪除可疑交易相關檔案
   * @param data
   * @param data.id Suspicious Transactions ID
   * @param data.fileType 檔案類型
   * - `1` 資訊審核
   * - `2` 風控審查
   * @param data.filename 檔案名稱
   */
  const deleteSuspiciousTransactionFiles = async (data: {
    id: string;
    fileType: string;
    filename: string;
  }) =>
    api.delete(
      `/admin/suspicious-txs/${data.id}/file?fileType=${
        data.fileType
      }&filename=${encodeURIComponent(data.filename)}`
    );

  return {
    login,
    forgotPassword,
    logout,
    changePassword,
    // admin api
    getAdmins,
    createAdmin,
    updateAdmin,
    deleteAdmin,
    // transactionpairs api
    getTransactionpairs,
    updateTransactionpairs,
    // banner api
    getBanners,
    createBanner,
    updateBanner,
    deleteBanner,
    updateBannerPriority,
    // mainnet api
    getMainnets,
    updateMainnets,
    // users api
    getUsers,
    createJuridicalUser,
    getExportUsersUrl,
    getUserInfo,
    getUserKycInfo,
    updateUserLevel,
    // user Statuslog
    getUserStatuslog,
    updateUserStatus,
    getExportUserStatuslogUrl,
    // user withdrawal whitelist
    getWithdrawalWhitelist,
    getExportWithdrawalWhitelistUrl,
    deleteWithdrawalWhitelist,
    // user loginlog
    getUserLoginLog,
    getExportUserLoginLogUrl,
    // user invite
    getUserInviteRecords,
    getUserInviteInfo,
    getUserInviteRewards,
    getExportUserInviteRewardsUrl,
    // check auth
    checkAuth,
    // KYC
    getKycList,
    getExportKycListUrl,
    updateUserKycComplianceReview,
    updateUserKycFinalReview,
    updateUserKycIdvAuditStatus,
    updateUserKycKryptoReview,
    resendUserKycKryptoReview,
    updateUserKycKryptoTaskID,
    getUserKycReviewslogs,
    getExportUserKycReviewslogsUrl,
    // KYC Name check
    updateKycNameCheck,
    downloadKycNameCheck,
    // KYC Risks
    getKycRisks,
    createKycRisk,
    updateKycRisk,
    deleteKycRisk,
    // Annual KYC Reviews
    getAnnualKYCReviews,
    getExportAnnualKYCReviewsUrl,
    // User KYC Risks
    getUserKycRisks,
    updateUserKycRisks,
    // User Bank
    getUserBank,
    updateUserBank,
    getUserBankLogs,
    getExportUserBankLogsUrl,
    // spots
    getSpots,
    getExportSpotsUrl,
    getExportSpotsAegisUrl,
    importSpotsAegis,
    // transactions
    getTransactions,
    getExportTransactionsUrl,
    // Receipts API
    getReceipts,
    getReceipt,
    getExportReceiptsUrl,
    issueReceipts,
    // User options
    getUserOptions,
    // 可疑交易管理
    getSuspiciousTransactions,
    getExportSuspiciousTransactionsUrl,
    getSuspiciousTransactionDetail,
    updateSuspiciousTransactionDetail,
    // 可疑交易檔案管理
    downloadSuspiciousTransactionFiles,
    uploadSuspiciousTransactionFiles,
    deleteSuspiciousTransactionFiles,
    // utils
    api,
  };
}
