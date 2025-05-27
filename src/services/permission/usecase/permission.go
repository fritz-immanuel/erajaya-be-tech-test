package usecase

import (
	"time"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"

	"github.com/fritz-immanuel/erajaya-be-tech-test/models"
	"github.com/fritz-immanuel/erajaya-be-tech-test/src/services/permission"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
)

type PermissionUsecase struct {
	permissionRepo permission.Repository
	contextTimeout time.Duration
	db             *sqlx.DB
}

func NewPermissionUsecase(db *sqlx.DB, permissionRepo permission.Repository) permission.Usecase {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	return &PermissionUsecase{
		permissionRepo: permissionRepo,
		contextTimeout: timeoutContext,
		db:             db,
	}
}

func (u *PermissionUsecase) FindAll(ctx *gin.Context, filterFindAllParams models.FindAllPermissionParams) ([]*models.Permission, *types.Error) {
	result, err := u.permissionRepo.FindAll(ctx, filterFindAllParams)
	if err != nil {
		err.Path = ".PermissionService->FindAll()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *PermissionUsecase) Find(ctx *gin.Context, id int) (*models.Permission, *types.Error) {
	result, err := u.permissionRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".PermissionService->Find()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *PermissionUsecase) Count(ctx *gin.Context, filterFindAllParams models.FindAllPermissionParams) (int, *types.Error) {
	result, err := u.permissionRepo.FindAll(ctx, filterFindAllParams)
	if err != nil {
		err.Path = ".PermissionService->Count()" + err.Path
		return 0, err
	}

	return len(result), nil
}
