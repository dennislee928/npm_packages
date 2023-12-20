export default function useReviewStatus() {
  const idvAuditStatus = ["未知", "未驗證", "接受", "拒絕"];
  const idvAuditStatusBadge = ["in-progress", "warn", "completed", "cancelled"];
  const idvState = [
    "未知",
    "通過",
    "需人工驗證",
    "不通過",
    "驗證中",
    "驗證初始化",
  ];
  const idvStateBadge = [
    "in-progress",
    "completed",
    "cancelled",
    "cancelled",
    "warn",
    "checking",
  ];
  return {
    idvAuditStatus,
    idvAuditStatusBadge,
    idvState,
    idvStateBadge,
  };
}
