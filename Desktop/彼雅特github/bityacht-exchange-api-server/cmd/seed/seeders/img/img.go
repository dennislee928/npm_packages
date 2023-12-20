package img

import (
	"bityacht-exchange-api-server/internal/pkg/datauri"
	"io"
	"os"
	"path"
)

var imgMap = make(map[string][]byte)

func Identification() []byte          { return getImg("identification.jpeg") }
func IdentificationBack() []byte      { return getImg("identification_back.jpeg") }
func Passport() []byte                { return getImg("passport.jpeg") }
func ResidentCertificate() []byte     { return getImg("resident_certificate.png") }
func ResidentCertificateBack() []byte { return getImg("resident_certificate_back.png") }
func ResultPass() []byte              { return getImg("result_pass.jpg") }
func ResultReject() []byte            { return getImg("result_reject.jpg") }

func getImg(filename string) []byte {
	img, ok := imgMap[filename]
	if ok {
		return img
	}
	defer func() {
		imgMap[filename] = img
	}()

	imgFile, err := os.Open(path.Clean("cmd/seed/seeders/img/" + filename))
	if err != nil {
		return nil
	}

	imgBytes, err := io.ReadAll(imgFile)
	if err != nil {
		return nil
	}

	img, _ = datauri.BytesToDataURI(imgBytes)
	return img
}
