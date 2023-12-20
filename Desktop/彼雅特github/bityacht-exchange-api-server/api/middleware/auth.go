package middleware

import (
	"bityacht-exchange-api-server/internal/cache/redis/managers"
	"bityacht-exchange-api-server/internal/cache/redis/users"
	sqlusers "bityacht-exchange-api-server/internal/database/sql/users"
	"bityacht-exchange-api-server/internal/database/sql/usersmodifylogs"
	errpkg "bityacht-exchange-api-server/internal/pkg/err"
	"bityacht-exchange-api-server/internal/pkg/jwt"
	"bityacht-exchange-api-server/internal/pkg/logger"
	"bityacht-exchange-api-server/internal/pkg/rbac"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const authHeaderKey = "Authorization"

// Authorization will check the Authorization header and Set JWT Claims to gin.Context.
func Authorization(jwtType jwt.Type) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get(authHeaderKey)
		if authHeader == "" {
			// TODO: Maybe better implement for export Authorization
			if ctx.Request.Method != http.MethodGet || !strings.Contains(ctx.Request.RequestURI, "export") {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errpkg.Error{Code: errpkg.CodeBadAuthorizationToken})
				return
			}

			authHeader = "Bearer " + ctx.Query("accessToken")
		}

		splitedAuthHeader := strings.Split(authHeader, " ")
		if len(splitedAuthHeader) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errpkg.Error{Code: errpkg.CodeBadAuthorizationToken})
			return
		}

		if splitedAuthHeader[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errpkg.Error{Code: errpkg.CodeBadAuthorizationToken})
			return
		}

		switch jwtType {
		case jwt.TypeManager:
			if claims, err := jwt.ValidateManager(splitedAuthHeader[1]); errpkg.Handler(ctx, err) {
				ctx.Abort()
			} else if err := managers.Validate(ctx, claims.ManagerPayload.ID, claims.RegisteredClaims.ID); errpkg.Handler(ctx, err) {
				ctx.Abort()
			} else if claims.Type != jwtType {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errpkg.Error{Code: errpkg.CodeBadAuthorizationToken})
			} else {
				ctx.Set(jwt.ClaimsKeyInGin, claims)
			}
		case jwt.TypeUser:
			if claims, err := jwt.ValidateUser(splitedAuthHeader[1]); errpkg.Handler(ctx, err) {
				ctx.Abort()
			} else if err := users.Validate(ctx, claims.UserPayload.ID, claims.RegisteredClaims.ID); errpkg.Handler(ctx, err) {
				ctx.Abort()
			} else if claims.Type != jwtType {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errpkg.Error{Code: errpkg.CodeBadAuthorizationToken})
			} else {
				ctx.Set(jwt.ClaimsKeyInGin, claims)
			}
		default:
			logger.Logger.Err(errors.New("bad jwt type")).Int("jwt type", int(jwtType)).Msg("bad usage of middleware.Authorization")
			ctx.Abort()
		}
	}
}

func ManagerRoleCheck(object string, action string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, err := jwt.GetClaimsFromGin[jwt.ManagerClaims](ctx)
		if errpkg.Handler(ctx, err) {
			return
		}

		// 檢查用戶是否擁有訪問該路由的權限
		if ok, err := rbac.Enforce(claims.ManagersRolesID, object, action); err != nil || !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errpkg.Error{Code: errpkg.CodePermissionDenied, Err: err})
		}
	}
}

func UserPermissionCheck(minNatureLevel int32, minJuridicalLevel int32, forbiddenForzen bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims, err := jwt.GetClaimsFromGin[jwt.UserClaims](ctx)
		if errpkg.Handler(ctx, err) {
			return
		}

		minLevel := minNatureLevel
		if claims.UserPayload.Type == sqlusers.TypeJuridicalPerson {
			minLevel = minJuridicalLevel
		}
		if minLevel < 0 {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errpkg.Error{Code: errpkg.CodePermissionDenied, Err: errors.New("bad path")})
			return
		} else if minLevel > 0 && claims.Level < minLevel {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errpkg.Error{Code: errpkg.CodePermissionDenied, Err: errors.New("bad level")})
			return
		}

		if forbiddenForzen && claims.Status == usersmodifylogs.SLStatusForzen {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errpkg.Error{Code: errpkg.CodePermissionDenied, Err: errors.New("bad status")})
			return
		}
	}
}
