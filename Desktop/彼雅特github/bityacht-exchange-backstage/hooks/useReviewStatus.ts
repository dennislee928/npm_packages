export default function useReviewStatus() {
  const status = [
    "全部",
    "待審核",
    "已通過",
    "已拒絕",
    "身分驗證審核未通過",
    "待複核",
    "待審核",
    "通過",
    "未通過",
  ];
  const badgeStatus = [
    "in-progress",
    "in-progress",
    "completed",
    "cancelled",
    "cancelled",
    "warn",
    "in-progress",
    "completed",
    "cancelled",
  ];
  return {
    status,
    badgeStatus,
  };
}
