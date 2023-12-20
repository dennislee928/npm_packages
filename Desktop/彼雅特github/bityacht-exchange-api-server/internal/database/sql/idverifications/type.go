package idverifications

import (
	"bityacht-exchange-api-server/internal/database/sql/usersmodifylogs"
	"bityacht-exchange-api-server/internal/pkg/kyc"
)

type Type int32

const (
	TypeManual Type = iota + 1
	TypeKryptoGO
)

type State int32

const (
	StateUnknown State = iota
	StateAccept
	StateReview
	StateReject
	StatePending
	StateInitial
)

func (s *State) ParseFromKryptoGO(state kyc.IDVState) {
	switch state {
	case kyc.IDVStateAccept:
		*s = StateAccept
	case kyc.IDVStateReview:
		*s = StateReview
	case kyc.IDVStateReject:
		*s = StateReject
	case kyc.IDVStatePending:
		*s = StatePending
	case kyc.IDVStateInitial:
		*s = StateInitial
	default:
		*s = StateUnknown
	}
}

type AuditStatus int32

const (
	AuditStatusUnknown AuditStatus = iota
	AuditStatusPending
	AuditStatusAccepted
	AuditStatusRejected
)

func (as *AuditStatus) ParseFromKryptoGO(auditStatus kyc.AuditStatus) {
	switch auditStatus {
	case kyc.AuditStatusPending:
		*as = AuditStatusPending
	case kyc.AuditStatusAccepted:
		*as = AuditStatusAccepted
	case kyc.AuditStatusRejected:
		*as = AuditStatusRejected
	default:
		*as = AuditStatusUnknown
	}
}

func (as AuditStatus) ToRLStatus() usersmodifylogs.RLStatus {
	switch as {
	case AuditStatusPending:
		return usersmodifylogs.RLIDVStatusPending
	case AuditStatusAccepted:
		return usersmodifylogs.RLIDVStatusAccepted
	case AuditStatusRejected:
		return usersmodifylogs.RLIDVStatusRejected
	default:
		return usersmodifylogs.RLStatusUnknown
	}
}

type UpdateImagesByURLRequest struct {
	IDImage        string
	IDBackImage    string
	PassportImage  string
	FaceImage      string
	IDAndFaceImage string
}
