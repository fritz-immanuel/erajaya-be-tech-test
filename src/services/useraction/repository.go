package useraction

import (
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
	"github.com/fritz-immanuel/erajaya-be-tech-test/models"
	"github.com/gin-gonic/gin"
)

// Repository is the contract between Repository and usecase
type Repository interface {
	FindAll(*gin.Context, models.FindAllActionHistory) ([]*models.UserAction, *types.Error)
	Find(*gin.Context, int) (*models.UserAction, *types.Error)
	FindPermission(*gin.Context, string, string) (*models.Permission, *types.Error)
	CreateManual(*gin.Context, *models.UserAction) *types.Error
	Update(*gin.Context, *models.UserAction) (*models.UserAction, *types.Error)
}
