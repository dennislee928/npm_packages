package modelpkg

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPaginatorFromQuery(ctx *gin.Context) Paginator {
	paginator := Paginator{PageSize: 25, Page: 1}

	if val := ctx.Query("pageSize"); val != "" {
		if iVal, err := strconv.Atoi(val); err == nil {
			switch iVal {
			case 25, 50, 100:
				paginator.PageSize = iVal
			}
		}
	}
	if val := ctx.Query("page"); val != "" {
		if iVal, err := strconv.Atoi(val); err == nil {
			paginator.Page = iVal
		}
	}

	return paginator
}

type Paginator struct {
	PageSize    int   `json:"pageSize,omitempty" form:"pageSize"` // Default: 25
	Page        int   `json:"page,omitempty" form:"page"`         // Default: 1
	TotalRecord int64 `json:"totalRecord,omitempty" swaggerignore:"true"`
}

func (p *Paginator) Offset() int {
	return (p.Page - 1) * p.PageSize
}
