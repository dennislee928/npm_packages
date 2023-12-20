package emailtemplates

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	modelpkg "bityacht-exchange-api-server/internal/pkg/model"
	"bytes"
	"embed"
	"html/template"
	"net/http"
	"time"
)

//go:embed */*.html.tmpl
var templatesFS embed.FS

var (
	templates *template.Template
)

func init() {
	templates = template.Must(template.ParseFS(templatesFS, "*/*.html.tmpl"))
}

func executeTemplate(name string, payload any) ([]byte, *errpkg.Error) {
	bytesBuffer := new(bytes.Buffer)
	data := map[string]any{
		"Payload": payload,
		"Logo":    template.URL("cid:logo.png"),
	}

	if err := templates.ExecuteTemplate(bytesBuffer, name, data); err != nil {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeExecuteTemplate, Err: err}
	}

	return bytesBuffer.Bytes(), nil
}

func FormatTime(t time.Time) string {
	return t.In(modelpkg.DefaultTimeLoc).Format("2006-01-02 15:04:05 (UTC -07:00)")
}
