package user

import (
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
	"github.com/fritz-immanuel/erajaya-be-tech-test/models"
	"github.com/gin-gonic/gin"
)

// Usecase is the contract between Repository and usecase
type Usecase interface {
	FindAll(context *gin.Context, params models.FindAllUserParams) ([]*models.User, *types.Error)
	Find(context *gin.Context, userID string) (*models.User, *types.Error)
	Count(context *gin.Context, params models.FindAllUserParams) (int, *types.Error)
	Create(context *gin.Context, data models.User) (*models.User, *types.Error)
	Update(context *gin.Context, userID string, data models.User) (*models.User, *types.Error)
	UpdateStatus(context *gin.Context, userID string, newStatusID string) (*models.User, *types.Error)

	// PASSWORD
	UpdatePassword(*gin.Context, string, string) (*models.User, *types.Error)

	Login(context *gin.Context, creds models.UserLogin) (*models.UserLogin, *types.Error)
}
