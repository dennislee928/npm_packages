package usersmodifylogs

type Type int32

const (
	TypeStatusLog Type = iota + 1
	TypeReviewLog
	TypeBankAccountLog
)

const badStatusChinese = "未知錯誤"
