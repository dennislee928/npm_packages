package suspicioustransactionspkg

import (
	"bityacht-exchange-api-server/internal/database/sql/usersspottransfers"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"errors"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

func Detect(action Action, spotTransfer *usersspottransfers.Model) ([]Result, *errpkg.Error) {
	var output []Result

	switch action {
	case ActionDepositCryptocurrency:
		if spotTransfer == nil {
			return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("spot transfer is nil")}
		}

		if results, err := detectMultipleWithdrawalAndSameAmount(action, spotTransfer.UsersID, spotTransfer.CurrenciesSymbol, spotTransfer.Amount); err != nil {
			return nil, err
		} else if len(results) > 0 {
			output = append(output, results...)
		}
	case ActionWithdrawCryptocurrency:
		if spotTransfer == nil {
			return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("spot transfer is nil")}
		}

		if results, err := detectMultipleWithdrawalAndSameAmount(action, spotTransfer.UsersID, spotTransfer.CurrenciesSymbol, spotTransfer.Amount); err != nil {
			return nil, err
		} else if len(results) > 0 {
			output = append(output, results...)
		}

		if results, err := detectMultipleSameWithdrawalAddress(spotTransfer.ToAddress); err != nil {
			return nil, err
		} else if len(results) > 0 {
			output = append(output, results...)
		}
	default:
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeNotImplement}
	}

	return output, nil
}

func detectMultipleWithdrawalAndSameAmount(action Action, usersID int64, currency string, amount decimal.Decimal) ([]Result, *errpkg.Error) {
	spotTransfers, err := queryTransferByUser(usersID, withActions(usersspottransfers.ActionWithdraw), withDuration(time.Hour))
	if err != nil {
		return nil, err
	}

	var fiatTransfers []any
	// TODO Get Fiat Transfers

	var output []Result
	if len(spotTransfers)+len(fiatTransfers) >= 3 {
		output = append(output, Result{
			Type: TypeMultipleWithdrawal,
			Information: Information{
				SpotTransfers: spotTransfers,
				// TODO FiatTransfers: fiatTransfers,
			},
		})
	}

	sameAmountResult := Result{Type: TypeMultipleSameAmount}
	if action == ActionWithdrawCryptocurrency {
		for _, transfer := range spotTransfers {
			if transfer.CurrenciesSymbol == currency && transfer.Amount.Equal(amount) {
				sameAmountResult.Information.SpotTransfers = append(sameAmountResult.Information.SpotTransfers, transfer)
			}
		}
	}
	// TODO else if action == ActionWithdrawFiat {
	// TODO }

	if len(sameAmountResult.Information.SpotTransfers) >= 3 { // TODO || len(sameAmountResult.Information.FiatTransfers) >= 3
		output = append(output, sameAmountResult)
	}

	return output, nil
}

func detectMultipleSameWithdrawalAddress(address string) ([]Result, *errpkg.Error) {
	spotTransfers, err := queryWithdrawTransferByToAddress(address, withDuration(24*time.Hour))
	if err != nil {
		return nil, err
	}

	accountMap := make(map[int64]struct{})
	for _, transfer := range spotTransfers {
		accountMap[transfer.UsersID] = struct{}{}
	}

	if len(accountMap) >= 2 {
		return []Result{{
			Type: TypeMultipleSameWithdrawalAddress,
			Information: Information{
				SpotTransfers: spotTransfers,
			},
		}}, nil
	}

	return nil, nil
}
