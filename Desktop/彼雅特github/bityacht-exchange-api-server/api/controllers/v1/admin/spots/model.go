package spots

import (
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"bityacht-exchange-api-server/internal/pkg/wallet"
	"mime/multipart"
)

type AegisExportRequest struct {
	// 主網
	// * 1: BTC
	// * 2: ETH
	// * 3: ERC20
	// * 4: TRC20
	Mainnet wallet.Mainnet `form:"mainnet" binding:"required,gte=1,lte=4"`

	// 交易時間（開始）
	StartAt modelpkg.Date `form:"startAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`

	// 交易時間（結束）
	EndAt modelpkg.Date `form:"endAt" swaggertype:"string" format:"date(YYYY/MM/DD)"`
}

type AegisImportRequest struct {
	CSV *multipart.FileHeader `form:"csv" swaggerignore:"true" binding:"required"`
}
