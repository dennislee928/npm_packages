package userscommissions

type Action int32

const (
	ActionDeposit Action = iota + 1
	ActionWithdraw
)

func (a Action) Chinese() string {
	switch a {
	case ActionDeposit:
		return "返佣"
	case ActionWithdraw:
		return "提領"
	}

	return "未知錯誤"
}
