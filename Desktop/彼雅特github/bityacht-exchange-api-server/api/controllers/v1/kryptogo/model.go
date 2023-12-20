package kryptogo

import (
	"bityacht-exchange-api-server/internal/database/sql/idverifications"
	"bityacht-exchange-api-server/internal/pkg/kyc"
	"strconv"
)

func idvCallbackReqToModel(req kyc.IDVCallbackRequest) idverifications.Model {
	record := idverifications.Model{
		TaskID: strconv.FormatInt(req.IDVTaskID, 10),
	}

	record.State.ParseFromKryptoGO(req.State)
	record.AuditStatus.ParseFromKryptoGO(req.AuditStatus)

	if req.AuditTimestamp == 0 {
		return record
	}

	if record.AuditTime.Time = kyc.ParseKryptoGOTimestamp(req.AuditTimestamp); !record.AuditTime.Time.IsZero() {
		record.AuditTime.Valid = true
	}

	return record
}

//! Deprecated (Message from TG 2023/11/28)
// func idvCallbackReqToUpdateArgs(req kyc.IDVCallbackRequest) (users.UpdateFromIDVCallbackArgs, *errpkg.Error) {
// 	args := users.UpdateFromIDVCallbackArgs{
// 		FirstName:  req.FirstName,
// 		LastName:   req.LastName,
// 		NationalID: req.IDNumber,
// 	}

//! Deprecated (Message from TG 2023/11/28)
// 	args.Gender.ParseFromKryptoGO(req.Gender)

// 	if req.Birthday != "" {
// 		if err := args.BirthDate.Parse(time.DateOnly, req.Birthday); err != nil {
// 			return args, err
// 		}
// 	}

// 	return args, nil
// }
