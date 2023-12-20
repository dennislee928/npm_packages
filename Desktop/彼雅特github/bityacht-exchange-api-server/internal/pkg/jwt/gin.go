package jwt

import (
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetClaimsFromGin[T ManagerClaims | UserClaims](ctx *gin.Context) (*T, *errpkg.Error) {
	if claims, ok := ctx.Get(ClaimsKeyInGin); !ok {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeGetClaims}
	} else if claims, ok := claims.(T); !ok {
		return nil, &errpkg.Error{HttpStatus: http.StatusInternalServerError, Code: errpkg.CodeBadClaimsType}
	} else {
		return &claims, nil
	}
}
