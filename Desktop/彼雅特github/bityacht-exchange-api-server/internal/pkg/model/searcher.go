package modelpkg

import (
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetSearcherFromQuery(ctx *gin.Context) Searcher {
	return Searcher{Value: ctx.Query("search")}
}

type Searcher struct {
	Value string `form:"search"` // Default: ''
}

// AddToQuery will check search value is empty or not, if not -> add Where condition to query. Don't add ` to fields's prefix and suffix, this function will do this for you.
func (s *Searcher) AddToQuery(db *gorm.DB, fields []string) *gorm.DB {
	if s.Value == "" || len(fields) == 0 {
		return db
	}

	wheres := make([]string, len(fields))
	args := make([]interface{}, len(fields))
	for i := range args {
		wheres[i] = addBackQuote(fields[i]) + " LIKE ?"
		args[i] = "%" + s.Value + "%"
	}

	return db.Where(strings.Join(wheres, " OR "), args...)
}
