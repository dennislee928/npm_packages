package wallets

import (
	"bityacht-exchange-api-server/internal/database/sql/walletsaddresses"
	"bityacht-exchange-api-server/internal/pkg/jwt"
	"bityacht-exchange-api-server/internal/pkg/wallet"
	"errors"
	"net/http"

	errpkg "bityacht-exchange-api-server/internal/pkg/err"

	qrcode "github.com/skip2/go-qrcode"

	"github.com/gin-gonic/gin"
)

// @Summary 	取得存款地址
// @Description 取得存款地址
// @Tags 		User-Wallets
// @Security	BearerAuth
// @Produce		json
// @Param 		query query GetDepositAddressRequest true "query"
// @Success 	200 {object} GetDepositAddressResponse
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/wallets/deposit/address [get]
func GetDepositAddressHandler(ctx *gin.Context) {
	var req GetDepositAddressRequest

	claims, wrapErr := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, wrapErr) {
		return
	}

	if err := ctx.ShouldBindQuery(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadQuery, err) {
		return
	}

	userID := claims.UserPayload.ID
	addressRecord, wrapErr := walletsaddresses.GetByUserMainnet(userID, req.Mainnet.BinanceNetwork())
	if errpkg.Handler(ctx, wrapErr) {
		return
	}

	if !addressRecord.Address.Valid || addressRecord.Address.String == "" {
		if addressRecord.TxID != "" {
			address, wrapErr := wallet.Cybavo.GetContractDepositAddress(ctx.Request.Context(), req.CurrencyType, req.Mainnet, addressRecord.TxID)
			if errpkg.Handler(ctx, wrapErr) {
				return
			}

			if wrapErr := walletsaddresses.UpdateAddress(&addressRecord, address); errpkg.Handler(ctx, wrapErr) {
				return
			}
		} else {
			errpkg.HandlerWithCode(ctx, http.StatusNotFound, errpkg.CodeWalletAddressNotGen, errors.New("wallet address not gen"))
			return
		}
	}

	bs, err := addressToQRCodeData(addressRecord.Address.String)
	if err != nil {
		errpkg.Handler(ctx, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeGenQRCode, Err: err})
		return
	}

	ctx.JSON(http.StatusOK, GetDepositAddressResponse{
		Address:    addressRecord.Address.String,
		QRCodeData: bs,
	})
}

// @Summary 	產生存款地址
// @Description 產生存款地址
// @Tags 		User-Wallets
// @Security	BearerAuth
// @Produce		json
// @Param 		body body GetDepositAddressRequest true "Request Body"
// @Success 	200 {object} GetDepositAddressResponse
// @Failure 	400 {object} errpkg.JsonError
// @Failure 	500 {object} errpkg.JsonError
// @Router 		/user/wallets/deposit/address [post]
func GenDepositAddressHandler(ctx *gin.Context) {
	var req GetDepositAddressRequest

	claims, wrapErr := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
	if errpkg.Handler(ctx, wrapErr) {
		return
	}

	if err := ctx.ShouldBindJSON(&req); errpkg.HandlerWithCode(ctx, http.StatusBadRequest, errpkg.CodeBadQuery, err) {
		return
	}

	_, wrapErr = walletsaddresses.GetByUserMainnet(claims.ID(), req.Mainnet.BinanceNetwork())
	if wrapErr == nil {
		errpkg.HandlerWithCode(ctx, http.StatusConflict, errpkg.CodeWalletAddressAlreadySet, errors.New("wallet address already set"))
		return
	}

	if wrapErr != nil &&
		wrapErr.Code != errpkg.CodeWalletAddressNotGen &&
		errpkg.Handler(ctx, wrapErr) {
		return
	}

	address, txID, wrapErr := wallet.Cybavo.CreateDepositAddress(ctx.Request.Context(), req.Mainnet)
	if errpkg.Handler(ctx, wrapErr) {
		return
	}

	if wrapErr := walletsaddresses.CreateAddress(claims.ID(), req.Mainnet.BinanceNetwork(), address, txID); errpkg.Handler(ctx, wrapErr) {
		return
	}

	if address == "" {
		ctx.Status(http.StatusCreated)
		return
	}

	bs, err := addressToQRCodeData(address)
	if err != nil {
		errpkg.Handler(ctx, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeGenQRCode, Err: err})
		return
	}

	ctx.JSON(http.StatusOK, GetDepositAddressResponse{
		Address:    address,
		QRCodeData: bs,
	})
}

func addressToQRCodeData(address string) ([]byte, error) {
	return qrcode.Encode(address, qrcode.Medium, 256)
}
