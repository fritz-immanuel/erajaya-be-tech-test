package permission

import (
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
	"github.com/fritz-immanuel/erajaya-be-tech-test/models"
	"github.com/gin-gonic/gin"
)

// Repository is the contract between Repository and usecase
type Repository interface {
	FindAll(*gin.Context, models.FindAllPermissionParams) ([]*models.Permission, *types.Error)
	Find(*gin.Context, int) (*models.Permission, *types.Error)
}
