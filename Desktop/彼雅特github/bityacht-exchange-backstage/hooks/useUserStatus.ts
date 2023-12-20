export default function useUserStatus() {
  const status = ["未啟用", "已啟用", "已停用", "凍結中"];
  const badgeStatus = ["in-progress", "completed", "cancelled", "frozen"];
  return {
    status,
    badgeStatus,
  };
}
