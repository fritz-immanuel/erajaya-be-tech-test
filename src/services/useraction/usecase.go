package useraction

import (
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
	"github.com/fritz-immanuel/erajaya-be-tech-test/models"
	"github.com/gin-gonic/gin"
)

// Usecase is the contract between Repository and usecase
type Usecase interface {
	FindAll(*gin.Context, models.FindAllActionHistory) ([]*models.UserAction, *types.Error)
	CreateManual(*gin.Context, models.UserAction) *types.Error
}
