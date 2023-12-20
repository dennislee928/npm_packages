export default function useSuspiciousReviewResult() {
  const status = ["全部", "待審核", "通過", "駁回"];
  const badgeStatus = ["in-progress", "in-progress", "completed", "cancelled"];
  return {
    status,
    badgeStatus,
  };
}
