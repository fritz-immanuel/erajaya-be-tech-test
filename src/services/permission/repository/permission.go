package repository

import (
	"fmt"
	"net/http"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library/data"
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
	"github.com/fritz-immanuel/erajaya-be-tech-test/models"
	"github.com/gin-gonic/gin"
)

// OutletRepository initialize object from model Outlet, to be used in database operation
type PermissionRepository struct {
	repository data.GenericStorage
}

// NewOutletRepository initialize service that provide connection to Database
func NewPermissionRepository(repository data.GenericStorage) PermissionRepository {
	//db := &models.DB{DB: configs.ActiveDB}
	return PermissionRepository{repository: repository}
}

// FindAll is a function to get all Data
func (s PermissionRepository) FindAll(ctx *gin.Context, params models.FindAllPermissionParams) ([]*models.Permission, *types.Error) {
	data := []*models.Permission{}

	var err error

	where := `TRUE`

	if params.Package != "" {
		where = fmt.Sprintf("%s AND package = '%s'", where, params.Package)
	}

	if params.Name != "" {
		where = fmt.Sprintf("%s AND name = '%s'", where, params.Name)
	}

	if params.IsHidden != 0 {
		where = fmt.Sprintf("%s AND is_hidden = %d", where, params.IsHidden)
	}

	if params.FindAllParams.DataFinder != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.DataFinder)
	}

	if params.FindAllParams.SortBy != "" {
		where = fmt.Sprintf("%s ORDER BY %s", where, params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where = fmt.Sprintf(`%s LIMIT :limit OFFSET :offset`, where)
	}

	query := fmt.Sprintf(`
  SELECT *
  FROM permissions
  WHERE %s
  `, where)

	err = s.repository.SelectWithQuery(ctx, &data, query, map[string]interface{}{
		"limit":  params.FindAllParams.Size,
		"offset": ((params.FindAllParams.Page - 1) * params.FindAllParams.Size),
	})

	if err != nil {
		return nil, &types.Error{
			Path:       ".PermissionStorage->FindAll()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return data, nil
}

// Find is a function to get by ID
func (s PermissionRepository) Find(ctx *gin.Context, id int) (*models.Permission, *types.Error) {
	var err error

	result := []*models.Permission{}

	query := `
  SELECT * FROM permissions
  WHERE users.id = :id`

	err = s.repository.SelectWithQuery(ctx, &result, query, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".PermissionStorage->Find()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(result) == 0 {
		return nil, &types.Error{
			Path:       ".PermissionStorage->Find()",
			Message:    "Data Not Found",
			Error:      data.ErrNotFound,
			StatusCode: http.StatusNotFound,
			Type:       "mysql-error",
		}
	}

	return result[0], nil
}
