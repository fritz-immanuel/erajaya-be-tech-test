package repository

import (
	"fmt"
	"net/http"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library/data"
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
	"github.com/fritz-immanuel/erajaya-be-tech-test/models"
	"github.com/gin-gonic/gin"
)

type UserPermissionRepository struct {
	repository data.GenericStorage
}

func NewUserPermissionRepository(repository data.GenericStorage) UserPermissionRepository {
	return UserPermissionRepository{repository: repository}
}

func (s UserPermissionRepository) FindAll(ctx *gin.Context, params models.FindAllUserPermissionParams) ([]*models.UserPermission, *types.Error) {
	result := []*models.UserPermission{}
	bulks := []*models.UserPermissionBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.DataFinder)
	}

	if params.UserID != "" {
		where += ` AND user_permissions.user_id = :user_id`
	}

	if params.FindAllParams.SortBy != "" {
		where = fmt.Sprintf("%s ORDER BY %s", where, params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where = fmt.Sprintf(`%s LIMIT :limit OFFSET :offset`, where)
	}

	query := fmt.Sprintf(`
  SELECT
    user_permissions.user_id,
    user_permissions.permission_id,
    permissions.package AS permission_package,
    permissions.module_name AS permission_module_name,
    permissions.action_name AS permission_action_name,
    permissions.http_method AS permission_http_method,
    permissions.route AS permission_route
  FROM user_permissions
  JOIN permissions on permissions.id = user_permissions.permission_id
  WHERE %s
  `, where)

	// fmt.Println(query)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"user_id": params.UserID,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".UserPermissionStorage->FindAll()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	for _, v := range bulks {
		obj := &models.UserPermission{
			UserID:       v.UserID,
			PermissionID: v.PermissionID,
			Permission: models.Permission{
				ID:         v.PermissionID,
				Package:    v.PermissionPackage,
				ModuleName: v.PermissionModuleName,
				ActionName: v.PermissionActionName,
				HTTPMethod: v.PermissionHTTPMethod,
				Route:      v.PermissionRoute,
			},
		}

		result = append(result, obj)
	}

	return result, nil
}

func (s UserPermissionRepository) Find(ctx *gin.Context, id string) (*models.UserPermission, *types.Error) {
	result := &models.UserPermission{}
	bulks := []*models.UserPermissionBulk{}

	var err error

	query := `
  SELECT
    user_permissions.user_id,
    user_permissions.permission_id,
    permissions.package AS permission_package,
    permissions.module_name AS permission_module_name,
    permissions.action_name AS permission_action_name,
    permissions.display_module_name AS permission_display_module_name,
    permissions.display_action_name AS permission_display_action_name,
    permissions.http_method AS permission_http_method,
    permissions.route AS permission_route
  FROM user_permissions
  JOIN permissions on permissions.id = user_permissions.permission_id
  WHERE user_permissions.id = :id`

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return nil, &types.Error{
			Path:       ".UserPermissionStorage->Find()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		v := bulks[0]
		result = &models.UserPermission{
			UserID:       v.UserID,
			PermissionID: v.PermissionID,
			Permission: models.Permission{
				ID:                v.PermissionID,
				Package:           v.PermissionPackage,
				ModuleName:        v.PermissionModuleName,
				ActionName:        v.PermissionActionName,
				DisplayModuleName: v.PermissionDisplayModuleName,
				DisplayActionName: v.PermissionDisplayActionName,
				HTTPMethod:        v.PermissionHTTPMethod,
				Route:             v.PermissionRoute,
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

	return result, nil
}

func (s UserPermissionRepository) Create(ctx *gin.Context, obj *models.CreateUpdateUserPermission) (*models.UserPermission, *types.Error) {
	data := models.UserPermission{}
	_, err := s.repository.Insert(ctx, obj)
	if err != nil {
		return nil, &types.Error{
			Path:       ".UserPermissionStorage->Create()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	data.UserID = obj.UserID
	data.PermissionID = obj.PermissionID
	return &data, nil
}

func (s UserPermissionRepository) DeleteByUserID(ctx *gin.Context, id string) *types.Error {
	args := make(map[string]interface{})
	err := s.repository.ExecQuery(ctx, fmt.Sprintf("DELETE FROM user_permissions WHERE user_id = '%s'", id), args)
	if err != nil {
		return &types.Error{
			Path:       ".UserPermissionStorage->DeleteByUserID()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return nil
}

func (s UserPermissionRepository) CreateBunch(ctx *gin.Context, userID string, params models.FindAllUserPermissionParams) *types.Error {
	args := make(map[string]interface{})

	where := "TRUE"
	not := ""

	if params.Package != "" {
		where = fmt.Sprintf("%s AND package = '%s'", where, params.Package)
	}

	if params.PermissionIDString != "" {
		where = fmt.Sprintf("%s AND id IN (%s)", where, params.PermissionIDString)
	}

	if params.Not != 0 {
		not = "NOT"
	}

	query := fmt.Sprintf(`
  INSERT INTO user_permissions (id, user_id, permission_id, created_at, updated_at)
  SELECT UUID(), "%s", id, UTC_TIMESTAMP + INTERVAL 7 hour, UTC_TIMESTAMP + INTERVAL 7 HOUR
  FROM (
    SELECT id FROM permissions
    WHERE %s AND id %s IN (
      SELECT permission_id FROM user_permissions
      WHERE user_id = "%s"
    )
  ) permission`, userID, where, not, userID)

	// fmt.Println(query)

	err := s.repository.ExecQuery(ctx, query, args)
	if err != nil {
		return &types.Error{
			Path:       ".UserPermissionStorage->CreateBunch()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return nil
}
