package product

import (
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
	"github.com/fritz-immanuel/erajaya-be-tech-test/models"
	"github.com/gin-gonic/gin"
)

// Repository is the contract between Repository and usecase
type Repository interface {
	FindAll(*gin.Context, models.FindAllProductParams) ([]*models.Product, *types.Error)
	Find(*gin.Context, string) (*models.Product, *types.Error)
	Create(*gin.Context, *models.Product) (*models.Product, *types.Error)
	Update(*gin.Context, *models.Product) (*models.Product, *types.Error)
	UpdateStatus(*gin.Context, string, string) (*models.Product, *types.Error)
}
