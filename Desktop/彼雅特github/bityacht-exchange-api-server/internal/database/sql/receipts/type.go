package receipts

type Status int32

const (
	StatusUnknown Status = iota
	StatusPending
	StatusIssuing
	StatusIssued
	StatusFailed
)

func (s Status) Int32() int32 {
	return int32(s)
}

func (s Status) Chinese() string {
	switch s {
	case StatusUnknown:
		return "未知"
	case StatusPending:
		return "未開立"
	case StatusIssuing:
		return "開立中"
	case StatusIssued:
		return "已開立"
	case StatusFailed:
		return "已失敗"
	}

	return "未知錯誤"
}

func canIssueStatus(status Status) bool {
	return status == StatusIssuing || status == StatusFailed
}
