package suspicioustxs

import (
	"bityacht-exchange-api-server/internal/database/sql/suspicioustransactions"
	"mime/multipart"
)

type GetDetailRequest struct {
	ID int64 `uri:"ID" binding:"required,gt=0"`
}

type UploadFileRequest struct {
	ID int64 `uri:"ID" binding:"required,gt=0" swaggerignore:"true"`

	// 上傳類型
	// * 1: 資訊審核
	// * 2: 風控審查
	UploadType suspicioustransactions.UpdateFilesType `form:"uploadType" binding:"required,gte=1,lte=2"`

	File *multipart.FileHeader `form:"file" binding:"required" swaggerignore:"true"`
}

type DownloadFileRequest struct {
	ID int64 `uri:"ID" binding:"required,gt=0" swaggerignore:"true"`

	// 檔案類型
	// * 1: 資訊審核
	// * 2: 風控審查
	FileType suspicioustransactions.UpdateFilesType `form:"fileType" binding:"required,gte=1,lte=2"`

	// 檔案名稱
	Filename string `form:"filename" binding:"required"`
}

type DeleteFileRequest DownloadFileRequest
