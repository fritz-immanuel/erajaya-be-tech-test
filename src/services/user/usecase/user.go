package usecase

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library"
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/data"
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/helpers"
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"
	"github.com/fritz-immanuel/erajaya-be-tech-test/src/services/user"
	"github.com/google/uuid"

	"github.com/fritz-immanuel/erajaya-be-tech-test/models"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/jmoiron/sqlx"
)

type UserUsecase struct {
	userRepo           user.Repository
	userpermissionRepo user.PermissionRepository
	contextTimeout     time.Duration
	db                 *sqlx.DB
}

func NewUserUsecase(db *sqlx.DB, userRepo user.Repository, userpermissionRepo user.PermissionRepository) user.Usecase {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second

	return &UserUsecase{
		userRepo:           userRepo,
		userpermissionRepo: userpermissionRepo,
		contextTimeout:     timeoutContext,
		db:                 db,
	}
}

func (u *UserUsecase) FindAll(ctx *gin.Context, params models.FindAllUserParams) ([]*models.User, *types.Error) {
	result, err := u.userRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".UserUsecase->FindAll()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *UserUsecase) Find(ctx *gin.Context, id string) (*models.User, *types.Error) {
	result, err := u.userRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".UserUsecase->Find()" + err.Path
		return nil, err
	}

	var permissionParams models.FindAllUserPermissionParams
	permissionParams.UserID = id
	permissionParams.FindAllParams.SortBy = "permissions.module_name, permissions.sequence_number_detail ASC"
	result.Permissions, err = u.userpermissionRepo.FindAll(ctx, permissionParams)
	if err != nil {
		err.Path = ".UserUsecase->Find()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *UserUsecase) Count(ctx *gin.Context, params models.FindAllUserParams) (int, *types.Error) {
	result, err := u.userRepo.FindAll(ctx, params)
	if err != nil {
		err.Path = ".UserUsecase->Count()" + err.Path
		return 0, err
	}

	return len(result), nil
}

func (u *UserUsecase) Create(ctx *gin.Context, obj models.User) (*models.User, *types.Error) {
	err := helpers.ValidateStruct(obj)
	if err != nil {
		err.Path = ".UserUsecase->Create()" + err.Path
		return nil, err
	}

	// check for duplicate username
	users, err := u.userRepo.FindAll(ctx, models.FindAllUserParams{Username: obj.Username})
	if err != nil {
		err.Path = ".UserUsecase->Create()" + err.Path
		return nil, err
	}

	if len(users) > 0 {
		return nil, &types.Error{
			Path:       ".UserUsecase->Create()",
			Message:    "Username already exists",
			Error:      data.ErrNotFound,
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "mysql-error",
		}
	}

	data := models.User{}
	data.ID = uuid.New().String()
	data.Name = obj.Name
	data.Email = obj.Email
	data.Username = obj.Username
	data.Password = obj.Password
	data.StatusID = models.DEFAULT_STATUS_CODE

	result, err := u.userRepo.Create(ctx, &data)
	if err != nil {
		err.Path = ".UserUsecase->Create()" + err.Path
		return nil, err
	}

	// create permission
	var permssions []string
	for _, v := range obj.Permissions {
		permssions = append(permssions, fmt.Sprintf(`%d`, v.PermissionID))
	}

	var permissionParams models.FindAllUserPermissionParams
	permissionParams.PermissionIDString = strings.Join(permssions, ",")
	permissionParams.Not = 1
	err = u.userpermissionRepo.CreateBunch(ctx, data.ID, permissionParams)
	if err != nil {
		err.Path = ".UserUsecase->Create()" + err.Path
		return nil, err
	}

	return result, nil
}

func (u *UserUsecase) Update(ctx *gin.Context, id string, obj models.User) (*models.User, *types.Error) {
	err := helpers.ValidateStruct(obj)
	if err != nil {
		err.Path = ".UserUsecase->Update()" + err.Path
		return nil, err
	}

	// check for duplicate username
	var dupeParams models.FindAllUserParams
	dupeParams.Username = obj.Username
	dupeParams.FindAllParams.DataFinder = fmt.Sprintf(`users.id != '%s'`, id)
	users, err := u.userRepo.FindAll(ctx, dupeParams)
	if err != nil {
		err.Path = ".UserUsecase->Update()" + err.Path
		return nil, err
	}

	if len(users) > 0 {
		return nil, &types.Error{
			Path:       ".UserUsecase->Update()",
			Message:    "Username already exists",
			Error:      fmt.Errorf("Username already exists"),
			StatusCode: http.StatusUnprocessableEntity,
			Type:       "mysql-error",
		}
	}

	data, err := u.userRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".UserUsecase->Update()" + err.Path
		return nil, err
	}

	data.Name = obj.Name
	data.Email = obj.Email
	data.Username = obj.Username

	result, err := u.userRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".UserUsecase->Update()" + err.Path
		return nil, err
	}

	// update permission
	err = u.userpermissionRepo.DeleteByUserID(ctx, id)
	if err != nil {
		err.Path = ".UserUsecase->Update()" + err.Path
		return nil, err
	}

	var permssions []string
	for _, v := range obj.Permissions {
		permssions = append(permssions, fmt.Sprintf(`%d`, v.PermissionID))
	}

	var permissionParams models.FindAllUserPermissionParams
	permissionParams.PermissionIDString = strings.Join(permssions, ",")
	permissionParams.Not = 1
	err = u.userpermissionRepo.CreateBunch(ctx, data.ID, permissionParams)
	if err != nil {
		err.Path = ".UserUsecase->Update()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *UserUsecase) UpdatePassword(ctx *gin.Context, id string, newPassword string) (*models.User, *types.Error) {
	data, err := u.userRepo.Find(ctx, id)
	if err != nil {
		err.Path = ".UserUsecase->UpdatePassword()" + err.Path
		return nil, err
	}

	hash := md5.New()
	io.WriteString(hash, newPassword)
	data.Password = fmt.Sprintf("%x", hash.Sum(nil))

	result, err := u.userRepo.Update(ctx, data)
	if err != nil {
		err.Path = ".UserUsecase->UpdatePassword()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *UserUsecase) UpdateStatus(ctx *gin.Context, id string, newStatusID string) (*models.User, *types.Error) {
	if newStatusID != models.STATUS_ACTIVE && newStatusID != models.STATUS_INACTIVE {
		return nil, &types.Error{
			Path:       ".UserUsecase->UpdateStatus()",
			Message:    "StatusID is not valid",
			Error:      fmt.Errorf("StatusID is not valid"),
			StatusCode: http.StatusBadRequest,
		}
	}

	result, err := u.userRepo.UpdateStatus(ctx, id, newStatusID)
	if err != nil {
		err.Path = ".UserUsecase->UpdateStatus()" + err.Path
		return nil, err
	}

	return result, err
}

func (u *UserUsecase) Login(ctx *gin.Context, creds models.UserLogin) (*models.UserLogin, *types.Error) {
	err := helpers.ValidateStruct(creds)
	if err != nil {
		err.Path = ".UserUsecase->Login()" + err.Path
		return nil, err
	}

	var userParams models.FindAllUserParams
	userParams.Username = creds.Username
	userParams.Password = creds.Password
	userParams.FindAllParams.StatusID = `status_id = 1`
	users, err := u.FindAll(ctx, userParams)
	if err != nil {
		err.Path = ".UserUsecase->Login()" + err.Path
		return nil, err
	}

	if len(users) == 0 {
		return nil, &types.Error{
			StatusCode: http.StatusUnauthorized,
			Message:    "Username / Password is incorrect",
			Error:      data.ErrNotFound,
			Path:       ".UserUsecase->Login()",
		}
	}

	user := users[0]

	credentials := library.Credential{ID: user.ID, Username: user.Username, Name: user.Name, Type: "WebAdmin"}

	token, errorJwtSign := library.JwtSignString(credentials)
	if errorJwtSign != nil {
		return nil, &types.Error{
			Error:      errorJwtSign,
			Message:    "Error JWT Sign String",
			Path:       ".UserUsecase->Login()",
			StatusCode: http.StatusInternalServerError,
		}
	}

	var permissionParams models.FindAllUserPermissionParams
	permissionParams.UserID = user.ID
	permissionParams.FindAllParams.SortBy = "permissions.module_name, permissions.sequence_number_detail ASC"
	creds.Permissions, err = u.userpermissionRepo.FindAll(ctx, permissionParams)
	if err != nil {
		err.Path = ".UserUsecase->Find()" + err.Path
		return nil, err
	}

	creds.ID = user.ID
	creds.Name = user.Name
	creds.Token = token
	creds.Password = ""

	return &creds, nil
}
