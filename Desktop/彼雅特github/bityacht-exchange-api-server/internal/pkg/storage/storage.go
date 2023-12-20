package storage

import (
	"bityacht-exchange-api-server/configs"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/logger"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func GetBannerPath(filename string) string {
	return path.Join(configs.Config.Server.StoragePath, "banners", filename)
}

func GetNameCheckPdfPath(filename string) string {
	return path.Join(configs.Config.Server.StoragePath, "nameCheckPdfs", filename)
}

func GetSuspiciousTxPath(folder string, filename string) string {
	return path.Join(configs.Config.Server.StoragePath, "suspiciousTx", folder, filename)
}

func CheckAndSaveUploadFile(allowedExtensions map[string]struct{}, fileHeader *multipart.FileHeader, filePath string, openFileParams int) *errpkg.Error {
	if fileHeader == nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("file header cannot be nil")}
	}
	src, err := fileHeader.Open()
	if err != nil {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadParam, Err: err}
	}
	defer func() {
		if err := src.Close(); err != nil {
			logger.Logger.Err(err).Str("filename", fileHeader.Filename).Msg("CheckAndSaveUploadFile Close src error")
		}
	}()

	if len(allowedExtensions) > 0 {
		// we only check the file extension.
		contentType := mime.TypeByExtension(path.Ext(fileHeader.Filename))
		if _, ok := allowedExtensions[contentType]; !ok {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadParam, Err: errors.New("file type not allowed")}
		}

		// Check Content Type from file content
		buffer := make([]byte, 512)
		if _, err = src.Read(buffer); err != nil {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadParam, Err: err}
		}
		bufContentType := http.DetectContentType(buffer)
		if contentType != bufContentType {
			return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadParam, Err: errors.New("bad content type of uploaded file")}
		}
	}

	// Save File
	if err := os.MkdirAll(filepath.Dir(filePath), 0750); err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSaveUploadedFile, Err: err}
	}
	var out *os.File
	switch openFileParams {
	case os.O_EXCL:
		out, err = os.OpenFile(path.Clean(filePath), os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	case os.O_TRUNC:
		out, err = os.OpenFile(path.Clean(filePath), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	default:
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSaveUploadedFile, Err: err}
	}
	if err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSaveUploadedFile, Err: err}
	}
	defer func() {
		if err := out.Close(); err != nil {
			logger.Logger.Err(err).Str("filePath", filePath).Msg("CheckAndSaveUploadFile Close out error")
		}
	}()
	if _, err = src.Seek(0, io.SeekStart); err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSaveUploadedFile, Err: err}
	}
	if _, err = io.Copy(out, src); err != nil {
		return &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSaveUploadedFile, Err: err}
	}
	return nil
}
