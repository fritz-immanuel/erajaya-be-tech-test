package repository

import (
	"fmt"
	"net/http"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library/appcontext"
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/data"
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
	"github.com/fritz-immanuel/erajaya-be-tech-test/models"
	"github.com/gin-gonic/gin"
)

type UserActionRepository struct {
	repository data.GenericStorage
}

func NewUserActionRepository(repository data.GenericStorage) UserActionRepository {
	return UserActionRepository{repository: repository}
}

func (s UserActionRepository) FindAll(ctx *gin.Context, params models.FindAllActionHistory) ([]*models.UserAction, *types.Error) {
	table_name := fmt.Sprintf("%s_status", params.TableName)
	data := []*models.UserAction{}
	bulks := []*models.UserActionBulk{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.DataFinder)
	}

	if params.FindAllParams.StatusID != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.StatusID)
	}

	if params.UserID != "" {
		where = fmt.Sprintf("%s AND user_actions.user_id = :user_id", where)
	}

	if params.RefID != "" {
		where = fmt.Sprintf("%s AND user_actions.ref_id = :ref_id", where)
	}

	if params.TableName != "" {
		where = fmt.Sprintf("%s AND user_actions.table_name = '%s'", where, params.TableName)
	}

	if params.FindAllParams.SortBy != "" {
		where = fmt.Sprintf("%s ORDER BY %s", where, params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where = fmt.Sprintf(`%s LIMIT :limit OFFSET :offset`, where)
	}

	query := ""

	if params.UsingStatusTable == 1 {
		query = fmt.Sprintf(`
    SELECT
      user_actions.*,
      IF(user_actions.action != 'Update Status', '-', status.name) status_name
    FROM user_actions
    LEFT JOIN status ON user_actions.action_value = status.id
    WHERE %s
    `, where)
	} else {
		query = fmt.Sprintf(`
    SELECT
      user_actions.*,
      IF(user_actions.action != 'Update Status', '-', %s.name) status_name
    FROM user_actions
    LEFT JOIN %s ON user_actions.action_value = %s.id
    WHERE %s
    `, table_name, table_name, table_name, where)
	}

	// fmt.Println(query)

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"limit":     params.FindAllParams.Size,
		"offset":    ((params.FindAllParams.Page - 1) * params.FindAllParams.Size),
		"status_id": params.FindAllParams.StatusID,
		"user_id":   params.UserID,
		"ref_id":    params.RefID,
	})

	if err != nil {
		return nil, &types.Error{
			Path:       ".UserActionStorage->FindAll()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	for _, v := range bulks {
		obj := &models.UserAction{
			ID:          v.ID,
			UserID:      v.UserID,
			TableName:   v.TableName,
			Action:      v.Action,
			ActionValue: v.ActionValue,
			CreatedAt:   v.CreatedAt,
			StatusName:  v.StatusName,
			RefID:       v.RefID,
		}

		data = append(data, obj)
	}

	return data, nil
}

func (s UserActionRepository) FindPermission(ctx *gin.Context, ModuleName string, PackageName string) (*models.Permission, *types.Error) {
	data := []*models.Permission{}
	var err error

	where := fmt.Sprintf(`name = '%s'`, ModuleName)

	if PackageName != "" {
		where = fmt.Sprintf(`%s AND package = '%s'`, where, PackageName)
	}

	query := fmt.Sprintf(`
  SELECT permissions.*
  FROM permissions
  WHERE %s
  `, where)

	// fmt.Println(query)

	err = s.repository.SelectWithQuery(ctx, &data, query, map[string]interface{}{})
	if err != nil {
		return nil, &types.Error{
			Path:       ".UserActionStorage->FindPermission()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(data) == 0 {
		return nil, &types.Error{
			Path:       ".UserActionStorage->FindPermission()",
			Message:    "Data Not Found",
			Error:      fmt.Errorf("Data Not Found"),
			StatusCode: http.StatusNotFound,
			Type:       "mysql-error",
		}
	}

	obj := &models.Permission{}
	for _, v := range data {
		obj = &models.Permission{
			ID:         v.ID,
			ModuleName: v.ModuleName,
			TableName:  v.TableName,
		}
		data = append(data, obj)
	}

	return obj, nil
}

func (s UserActionRepository) CreateManual(ctx *gin.Context, obj *models.UserAction) *types.Error {
	currentUserID := appcontext.UserID(ctx)
	if currentUserID == nil {
		currentUserID = appcontext.EmployeeID(ctx)
	}

	insertQuery := fmt.Sprintf("INSERT INTO user_actions (id, user_id, table_name, action, ref_id, created_at) VALUES (UUID(), '%s', '%s', '%s', '%s', UTC_TIMESTAMP + INTERVAL 7 HOUR)", *currentUserID, obj.TableName, obj.Action, obj.RefID)

	err := s.repository.ExecQuery(ctx, insertQuery, nil)
	if err != nil {
		return &types.Error{
			Path:       ".UserActionStorage->CreateManual()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return nil

}

func (s UserActionRepository) Find(ctx *gin.Context, id int) (*models.UserAction, *types.Error) {
	result := models.UserAction{}
	bulks := []*models.UserAction{}

	var err error

	query := `
  SELECT
    user_actions.id, user_actions.user_id, user_actions.table_name,
    user_actions.action, user_actions.action_value, user_actions.ref_id
  FROM user_actions
  WHERE user_actions.id = :id FOR UPDATE`

	err = s.repository.SelectWithQuery(ctx, &bulks, query, map[string]interface{}{
		"id": id,
	})
	if err != nil {

		return nil, &types.Error{
			Path:       ".UserActionStorage->Find()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	if len(bulks) > 0 {
		result = *bulks[0]
	}

	return &result, nil
}

func (s UserActionRepository) Update(ctx *gin.Context, obj *models.UserAction) (*models.UserAction, *types.Error) {
	data := models.UserAction{}
	err := s.repository.Update(ctx, obj)
	if err != nil {

		return nil, &types.Error{
			Path:       ".UserActionStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	err = s.repository.FindByID(ctx, &data, obj.ID)
	if err != nil {
		return nil, &types.Error{
			Path:       ".UserActionStorage->Update()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}
	return &data, nil
}

func (s UserActionRepository) FindAllQueueMaster(ctx *gin.Context, params models.FindAllActionHistory) ([]*models.UserAction, *types.Error) {
	data := []*models.UserAction{}

	var err error

	where := `TRUE`

	if params.FindAllParams.DataFinder != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.DataFinder)
	}

	if params.FindAllParams.StatusID != "" {
		where = fmt.Sprintf("%s AND %s", where, params.FindAllParams.StatusID)
	}

	if params.UserID != "" {
		where = fmt.Sprintf("%s AND user_actions.user_id = :user_id", where)
	}

	if params.RefID != "" {
		where = fmt.Sprintf("%s AND user_actions.ref_id = :ref_id", where)
	}

	if params.TableName != "" {
		where = fmt.Sprintf("%s AND user_actions.table_name = '%s'", where, params.TableName)
	}

	if params.GroupBy != "" {
		where = fmt.Sprintf("GROUP BY %s", params.GroupBy)
	}

	if params.FindAllParams.SortBy != "" {
		where = fmt.Sprintf("%s ORDER BY %s", where, params.FindAllParams.SortBy)
	}

	if params.FindAllParams.Page > 0 && params.FindAllParams.Size > 0 {
		where = fmt.Sprintf(`%s LIMIT :limit OFFSET :offset`, where)
	}

	query := ""

	query = fmt.Sprintf(`
  SELECT
    user_actions.id, user_actions.user_id, user_actions.table_name
    user_actions.action, user_actions.action_value, user_actions.ref_id
  FROM user_actions
  WHERE user_actions.ref_id != 0 AND %s
  `, where)

	err = s.repository.SelectWithQuery(ctx, &data, query, map[string]interface{}{
		"limit":     params.FindAllParams.Size,
		"offset":    ((params.FindAllParams.Page - 1) * params.FindAllParams.Size),
		"status_id": params.FindAllParams.StatusID,
		"ref_id":    params.RefID,
		"user_id":   params.UserID,
	})

	if err != nil {
		return nil, &types.Error{
			Path:       ".UserActionStorage->FindAllQueueMaster()",
			Message:    err.Error(),
			Error:      err,
			StatusCode: http.StatusInternalServerError,
			Type:       "mysql-error",
		}
	}

	return data, nil
}
