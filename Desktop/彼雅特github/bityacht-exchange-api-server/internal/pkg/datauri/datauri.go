package datauri

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/logger"
	"encoding/base64"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"image/jpeg"
	"image/png"

	"github.com/rs/zerolog"
)

const (
	// Ref: https://en.wikipedia.org/wiki/Data_URI_scheme
	imageDataPrefix       = "data:image/"
	imageDataPrefixLength = len(imageDataPrefix)

	dataEncoding = "base64"

	// Ref: https://mariadb.com/kb/en/mediumblob/
	mediumBlobMaxSize = 16777215
)

// only Accept => data:image/$FormatName;base64,$Data
func ValidateImage(dataURI string) *errpkg.Error {
	if imageSize := len(dataURI); imageSize == 0 {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadBody}
	} else if imageSize > mediumBlobMaxSize {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeImageOverSize}
	} else if !strings.HasPrefix(dataURI, imageDataPrefix) {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadDataURIImageFormat}
	}

	semicolonIndex := strings.IndexByte(dataURI, ';')
	commaIndex := strings.IndexByte(dataURI, ',')
	if semicolonIndex == -1 || commaIndex == -1 {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadDataURIImageFormat}
	}

	formatName := dataURI[imageDataPrefixLength:semicolonIndex]
	encoding := dataURI[semicolonIndex+1 : commaIndex]
	if encoding != dataEncoding {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadDataURIImageFormat, Err: errors.New("bad encoding")}
	}

	var err error
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(dataURI[commaIndex+1:]))

	switch formatName {
	case "png":
		_, err = png.DecodeConfig(reader)
	case "jpeg":
		_, err = jpeg.DecodeConfig(reader)
	default:
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadUploadedFileType}
	}

	if err != nil {
		return &errpkg.Error{HttpStatus: http.StatusBadRequest, Code: errpkg.CodeBadImageData, Err: err}
	}

	return nil
}

func BytesToDataURI(data []byte) ([]byte, *errpkg.Error) {
	base64Encoding := "data:image/"
	mimeType := http.DetectContentType(data)
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "jpeg;base64,"
	case "image/png":
		base64Encoding += "png;base64,"
	default:
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadImageData, Err: errors.New("bad content type: " + mimeType)}
	}

	output := make([]byte, len(base64Encoding)+base64.StdEncoding.EncodedLen(len(data)))
	copy(output, []byte(base64Encoding))
	base64.StdEncoding.Encode(output[len(base64Encoding):], data)

	return output, nil
}

func DownloadImage(errLogger zerolog.Logger, rawURL string) []byte {
	errLogger = errLogger.With().Str("url", rawURL).Logger()

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		errLogger.Err(err).Msg("Parse URL error")
		return nil
	}

	resp, err := http.Get(parsedURL.String())
	if err != nil {
		errLogger.Err(err).Msg("Get Image error")
		return nil
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			errLogger.Err(err).Msg("Close Body error")
		}
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		errLogger.Err(err).Msg("ReadAll error")
		return nil
	} else if output, err := BytesToDataURI(respBody); err != nil {
		errLogger.Err(err.Err).Send()
		return nil
	} else {
		return output
	}
}

func GetDataURIFromFileHeader(file *multipart.FileHeader) ([]byte, *errpkg.Error) {
	if file == nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadCoding, Err: errors.New("file is nil")}
	}

	imageFile, err := file.Open()
	if err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSaveUploadedFile, Err: err}
	}
	defer func() {
		if err := imageFile.Close(); err != nil {
			logger.Logger.Err(err).Str("service", "get data uri from file header").Msg("close file error")
		}
	}()

	if bytesImage, err := io.ReadAll(imageFile); err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeSaveUploadedFile, Err: err}
	} else if dataURIResultImage, err := BytesToDataURI(bytesImage); err != nil {
		return nil, err
	} else {
		return dataURIResultImage, nil
	}
}
