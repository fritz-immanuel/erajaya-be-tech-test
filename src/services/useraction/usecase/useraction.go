package usecase

import (
	"time"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
	"github.com/fritz-immanuel/erajaya-be-tech-test/src/services/useraction"

	"github.com/fritz-immanuel/erajaya-be-tech-test/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
)

type UserAction struct {
	repo           useraction.Repository
	contextTimeout time.Duration
	db             *sqlx.DB
}

func NewUserAction(db *sqlx.DB, repo useraction.Repository) useraction.Usecase {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	return &UserAction{
		repo:           repo,
		contextTimeout: timeoutContext,
		db:             db,
	}
}

func (u *UserAction) FindAll(ctx *gin.Context, filterFindAllParams models.FindAllActionHistory) ([]*models.UserAction, *types.Error) {
	resultPermission, err := u.repo.FindPermission(ctx, filterFindAllParams.ModuleName, filterFindAllParams.PackageName)
	if err != nil {
		err.Path = ".UserActionService->FindAll()" + err.Path
		return nil, err
	}

	filterFindAllParams.TableName = resultPermission.TableName
	result, err := u.repo.FindAll(ctx, filterFindAllParams)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *UserAction) CreateManual(ctx *gin.Context, obj models.UserAction) *types.Error {
	err := u.repo.CreateManual(ctx, &obj)
	if err != nil {
		err.Path = ".UserActionService->CreateManual()" + err.Path
		return err
	}

	return nil
}
