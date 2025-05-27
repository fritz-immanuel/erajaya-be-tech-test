package repository

import (
	"fmt"
	"net/http"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library/data"
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
	"github.com/fritz-immanuel/erajaya-be-tech-test/models"
	"github.com/gin-gonic/gin"
)

type ProductRepository struct {
	repository       data.GenericStorage
	statusRepository data.GenericStorage
}

func NewProductRepository(repository data.GenericStorage, statusRepository data.GenericStorage) ProductRepository {
	return ProductRepository{repository: repository, statusRepository: statusRepository}
}

// A function to get all Data that matches the filter provided
func (s ProductRepository) FindAll(ctx *gin.Context, params models.FindAllProductParams) ([]*models.Product, *types.Error) {
	result := []*models.Product{}
	bulks := []*models.ProductBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.DataFinder)
	}

	if params.FindAllParams.BusinessID != "" {
		where += fmt.Sprintf(` AND products.%s`, params.FindAllParams.BusinessID)
	}

	if params.FindAllParams.StatusID != "" {
		where += fmt.Sprintf(` AND products.%s`, params.FindAllParams.StatusID)
	}

	if params.FindAllParams.SortBy != "" {
		where = fmt.Sprintf("%s ORDER BY %s", where, params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where = fmt.Sprintf(`%s LIMIT :limit OFFSET :offset`, where)
	}

	query := fmt.Sprintf(`
  SELECT
    products.id, products.name, products.price, products.description, products.quantity,
    products.status_id, products.created_at,
    status.name AS status_name
  FROM products
  JOIN status ON products.status_id = status.id
  WHERE %s
  `, where)

	// fmt.Println(query)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"limit":  params.FindAllParams.Size,
		"offset": ((params.FindAllParams.Page - 1) * params.FindAllParams.Size),
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".ProductStorage->FindAll()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		for _, v := range bulks {
			result = append(result, &models.Product{
				ID:          v.ID,
				Name:        v.Name,
				Price:       v.Price,
				Description: v.Description,
				Quantity:    v.Quantity,
				StatusID:    v.StatusID,
				Status: models.Status{
					ID:   v.StatusID,
					Name: v.StatusName,
				},
				CreatedAt: v.CreatedAt,
			})
		}
	}

	return result, nil
}

// A function to get a row of data specified by the given ID
func (s ProductRepository) Find(ctx *gin.Context, id string) (*models.Product, *types.Error) {
	result := models.Product{}
	bulks := []*models.ProductBulk{}
	var err error

	query := `
  SELECT
    products.id, products.name, products.price, products.description, products.quantity,
    products.status_id, products.created_at,
    status.name AS status_name
  FROM products
  JOIN status ON products.status_id = status.id
  WHERE products.id = :id`

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".ProductStorage->Find()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		v := bulks[0]
		result = models.Product{
			ID:          v.ID,
			Name:        v.Name,
			Price:       v.Price,
			Description: v.Description,
			Quantity:    v.Quantity,
			StatusID:    v.StatusID,
			Status: models.Status{
				ID:   v.StatusID,
				Name: v.StatusName,
			},
			CreatedAt: v.CreatedAt,
		}
	} else {
		return nil, &types.Error{
			Path:       ".ProductStorage->Find()",
			Message:    "Data Not Found",
			Error:      data.ErrNotFound,
			StatusCode: http.StatusNotFound,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

// Inserts a new row of data
func (s ProductRepository) Create(ctx *gin.Context, obj *models.Product) (*models.Product, *types.Error) {
	data := models.Product{}
	_, err := s.repository.Insert(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ProductStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ProductStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

// Updates a row of data specified by the given ID inside the obj struct
func (s ProductRepository) Update(ctx *gin.Context, obj *models.Product) (*models.Product, *types.Error) {
	data := models.Product{}
	err := s.repository.Update(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ProductStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ProductStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s ProductRepository) UpdateStatus(ctx *gin.Context, id string, statusID string) (*models.Product, *types.Error) {
	data := models.Product{}
	err := s.repository.UpdateStatus(ctx, id, statusID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ProductStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, id)
	if err != nil {
		return nil, &types.Error{
			Path:       ".ProductStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &data, nil
}
