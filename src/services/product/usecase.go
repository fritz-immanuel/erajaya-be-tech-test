package product

import (
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
	"github.com/fritz-immanuel/erajaya-be-tech-test/models"
	"github.com/gin-gonic/gin"
)

// Usecase is the contract between Repository and usecase
type Usecase interface {
	FindAll(context *gin.Context, params models.FindAllProductParams) ([]*models.Product, *types.Error)
	Find(context *gin.Context, id string) (*models.Product, *types.Error)
	Count(context *gin.Context, params models.FindAllProductParams) (int, *types.Error)
	Create(context *gin.Context, newData models.Product) (*models.Product, *types.Error)
	Update(context *gin.Context, id string, updatedData models.Product) (*models.Product, *types.Error)

	UpdateStatus(*gin.Context, string, string) (*models.Product, *types.Error)
}
