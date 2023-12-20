export default function useOrderStatus() {
  const status = ["全部", "處理中", "已完成", "已失敗", "已取消", "審核中"];
  const badgeStatus = [
    "in-progress",
    "in-progress",
    "completed",
    "cancelled",
    "frozen",
    "checking",
  ];
  return {
    status,
    badgeStatus,
  };
}
