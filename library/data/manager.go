package data

import (
	"fmt"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Manager represents the manager to manage the data consistency
type Manager struct {
	db *sqlx.DB
}

// RunInTransaction runs the f with the transaction queryable inside the context
func (m *Manager) RunInTransaction(ctx *gin.Context, f func(tctx *gin.Context) *types.Error) *types.Error {
	tx, err := m.db.Beginx()
	if err != nil {
		tx.Rollback()

		err := &types.Error{
			Path:    ".DealingHandler->Create()",
			Message: fmt.Sprintf("error when creating transction: %v", err),
			Error:   fmt.Errorf("error when creating transction: %v", err),
			Type:    "golang-error",
		}
		return err
	}

	ctx = NewContext(ctx, tx)
	if err != nil {
		fmt.Printf("\n[RunInTransaction - Prepare] Error: %v\n", err)
	}
	errTransaction := f(ctx)
	if errTransaction != nil {
		tx.Rollback()
		return errTransaction
	}

	err = tx.Commit()
	if err != nil {
		err := &types.Error{
			Path:    ".DealingHandler->Create()",
			Message: fmt.Sprintf("error when committing transaction: %v", err),
			Error:   fmt.Errorf("error when committing transaction: %v", err),
			Type:    "golang-error",
		}
		return err
	}

	return nil
}

// NewManager creates a new manager
func NewManager(
	db *sqlx.DB,
) *Manager {
	return &Manager{
		db: db,
	}
}
