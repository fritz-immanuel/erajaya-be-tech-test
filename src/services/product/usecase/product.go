package usecase

import (
	"time"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
	"github.com/fritz-immanuel/erajaya-be-tech-test/src/services/product"
	"github.com/google/uuid"

	"github.com/fritz-immanuel/erajaya-be-tech-test/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
)

type ProductUsecase struct {
	productRepo    product.Repository
	contextTimeout time.Duration
	db             *sqlx.DB
}

func NewProductUsecase(db *sqlx.DB, productRepo product.Repository) product.Usecase {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	return &ProductUsecase{
		productRepo:    productRepo,
		contextTimeout: timeoutContext,
		db:             db,
	}
}

func (u *ProductUsecase) FindAll(ctx *gin.Context, params models.FindAllProductParams) ([]*models.Product, *types.Error) {
	result, err := u.productRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".ProductUsecase->FindAll()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *ProductUsecase) Find(ctx *gin.Context, id string) (*models.Product, *types.Error) {
	result, err := u.productRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".ProductUsecase->Find()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *ProductUsecase) Count(ctx *gin.Context, params models.FindAllProductParams) (int, *types.Error) {
	result, err := u.productRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".ProductUsecase->Count()" + err.Path
		return 0, err
	}

	return len(result), nil
}

func (u *ProductUsecase) Create(ctx *gin.Context, obj models.Product) (*models.Product, *types.Error) {
	data := models.Product{
		ID:          uuid.New().String(),
		Name:        obj.Name,
		Price:       obj.Price,
		Description: obj.Description,
		Quantity:    obj.Quantity,
		StatusID:    models.DEFAULT_STATUS_CODE,
	}

	result, err := u.productRepo.Create(ctx, &data)
	if err != nil {
		err.Path = ".ProductUsecase->Create()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *ProductUsecase) Update(ctx *gin.Context, id string, obj models.Product) (*models.Product, *types.Error) {
	data, err := u.productRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".ProductUsecase->Update()" + err.Path
		return nil, err
	}

	data.Name = obj.Name
	data.Price = obj.Price
	data.Description = obj.Description
	data.Quantity = obj.Quantity

	result, err := u.productRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".ProductUsecase->Update()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *ProductUsecase) UpdateStatus(ctx *gin.Context, id string, newStatusID string) (*models.Product, *types.Error) {
	result, err := u.productRepo.UpdateStatus(ctx, id, newStatusID)
	if err != nil {
		err.Path = ".ProductUsecase->UpdateStatus()" + err.Path
		return nil, err
	}

	return result, nil
}
