package data

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type key int

const (
	txKey key = 0
)

// Queryer represents the database commands interface
type Queryer interface {
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	Rebind(query string) string
	MustExec(query string, args ...interface{}) sql.Result
	Select(dest interface{}, query string, args ...interface{}) error
	Get(dest interface{}, query string, args ...interface{}) error
}

// NewContext creates a new data context
func NewContext(ctx *gin.Context, q Queryer) *gin.Context {
	ctx.Set("transaction", q)
	return ctx
}

// TxFromContext returns the trasanction object from the context
func TxFromContext(ctx *gin.Context) (Queryer, bool) {
	q := (*ctx).Value("transaction")
	if q == nil {
		return nil, false
	}
	return q.(Queryer), true
}
