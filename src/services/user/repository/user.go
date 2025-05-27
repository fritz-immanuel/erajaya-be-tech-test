package repository

import (
	"fmt"
	"net/http"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library/data"
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
	"github.com/fritz-immanuel/erajaya-be-tech-test/models"
	"github.com/gin-gonic/gin"
)

type UserRepository struct {
	repository       data.GenericStorage
	statusRepository data.GenericStorage
}

func NewUserRepository(repository data.GenericStorage, statusRepository data.GenericStorage) UserRepository {
	return UserRepository{repository: repository, statusRepository: statusRepository}
}

func (s UserRepository) FindAll(ctx *gin.Context, params models.FindAllUserParams) ([]*models.User, *types.Error) {
	result := []*models.User{}
	bulks := []*models.UserBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.DataFinder)
	}

	if params.FindAllParams.StatusID != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.StatusID)
	}

	if params.Email != "" {
		where += ` AND users.email = :email`
	}

	if params.Username != "" {
		where += ` AND users.username = :username`
	}

	if params.Password != "" {
		where += ` AND users.password = :password`
	}

	if params.FindAllParams.SortBy != "" {
		where = fmt.Sprintf("%s ORDER BY %s", where, params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where = fmt.Sprintf(`%s LIMIT :limit OFFSET :offset`, where)
	}

	query := fmt.Sprintf(`
  SELECT
    users.id, users.name, users.email, users.password, users.username,
    users.status_id, status.name AS status_name
  FROM users
  JOIN status ON status.id = users.status_id
  WHERE %s
  `, where)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"limit":     params.FindAllParams.Size,
		"offset":    ((params.FindAllParams.Page - 1) * params.FindAllParams.Size),
		"status_id": params.FindAllParams.StatusID,
		"email":     params.Email,
		"username":  params.Username,
		"password":  params.Password,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".UserStorage->FindAll()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	for _, v := range bulks {
		obj := &models.User{
			ID:       v.ID,
			Name:     v.Name,
			Email:    v.Email,
			Username: v.Username,
			Password: v.Password,
			StatusID: v.StatusID,
			Status: models.Status{
				ID:   v.StatusID,
				Name: v.StatusName,
			},
		}

		result = append(result, obj)
	}

	return result, nil
}

func (s UserRepository) Find(ctx *gin.Context, id string) (*models.User, *types.Error) {
	var err error

	result := models.User{}
	bulks := []*models.UserBulk{}

	query := `
  SELECT
    users.id, users.name, users.email, users.password, users.username,
    users.status_id, status.name AS status_name
  FROM users
  JOIN status on status.id = users.status_id
  WHERE users.id = :id`

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".UserStorage->Find()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		v := bulks[0]
		result = models.User{
			ID:       v.ID,
			Name:     v.Name,
			Email:    v.Email,
			Username: v.Username,
			Password: v.Password,
			StatusID: v.StatusID,
			Status: models.Status{
				ID:   v.StatusID,
				Name: v.StatusName,
			},
		}
	} else {
		return nil, &types.Error{
			Path:       ".UserStorage->Find()",
			Message:    "Data Not Found",
			Error:      data.ErrNotFound,
			StatusCode: http.StatusNotFound,
			Type:       "mysql-error",
		}
	}

	return &result, nil
}

func (s UserRepository) Create(ctx *gin.Context, obj *models.User) (*models.User, *types.Error) {
	data := models.User{}
	_, err := s.repository.Insert(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".UserStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".UserStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s UserRepository) Update(ctx *gin.Context, obj *models.User) (*models.User, *types.Error) {
	data := models.User{}
	err := s.repository.Update(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".UserStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".UserStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s UserRepository) UpdateStatus(ctx *gin.Context, id string, statusID string) (*models.User, *types.Error) {
	data := models.User{}
	err := s.repository.UpdateStatus(ctx, id, statusID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".UserStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, id)
	if err != nil {
		return nil, &types.Error{
			Path:       ".UserStorage->UpdateStatus()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return &data, nil
}
