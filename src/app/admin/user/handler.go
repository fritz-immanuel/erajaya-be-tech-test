package user

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library"
	"github.com/fritz-immanuel/erajaya-be-tech-test/middleware"
	"github.com/fritz-immanuel/erajaya-be-tech-test/models"
	"github.com/gin-gonic/gin"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library/data"
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/helpers"
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/http/response"
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"

	"github.com/fritz-immanuel/erajaya-be-tech-test/src/services/user"

	userRepository "github.com/fritz-immanuel/erajaya-be-tech-test/src/services/user/repository"
	userUsecase "github.com/fritz-immanuel/erajaya-be-tech-test/src/services/user/usecase"
)

var ()

type UserHandler struct {
	UserUsecase user.Usecase
	dataManager *data.Manager
	Result      gin.H
	Status      int
}

func (h UserHandler) RegisterAPI(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine, v *gin.RouterGroup) {
	userRepo := userRepository.NewUserRepository(
		data.NewMySQLStorage(db, "users", models.User{}, data.MysqlConfig{}),
		data.NewMySQLStorage(db, "status", models.Status{}, data.MysqlConfig{}),
	)

	userpermissionRepo := userRepository.NewUserPermissionRepository(
		data.NewMySQLStorage(db, "user_permissions", models.UserPermission{}, data.MysqlConfig{}),
	)

	uUser := userUsecase.NewUserUsecase(db, userRepo, userpermissionRepo)

	base := &UserHandler{UserUsecase: uUser, dataManager: dataManager}

	rs := v.Group("/users")
	{
		rs.GET("", middleware.Auth, base.FindAll)
		rs.GET("/:id", middleware.Auth, base.Find)
		rs.POST("", middleware.Auth, base.Create)
		rs.PUT("/:id", middleware.Auth, base.Update)

		rs.PUT("/:id/password", middleware.Auth, base.UpdatePassword)
		rs.PUT("/:id/reset-password", middleware.Auth, base.ResetPassword)

		rs.PUT("/:id/status", middleware.Auth, base.UpdateStatus)
	}

	rsa := v.Group("/users/auth")
	{
		rsa.POST("/login", base.Login)
	}

	rss := v.Group("/statuses")
	{
		rss.GET("/users", base.FindStatus)
	}
}

func (h *UserHandler) FindAll(c *gin.Context) {
	var params models.FindAllUserParams
	page, size := helpers.FilterFindAll(c)
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams
	datas, err := h.UserUsecase.FindAll(c, params)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, err.Message, http.StatusInternalServerError, *err)
			return
		}
	}

	for _, u := range datas {
		u.Password = ""
	}

	params.FindAllParams.Page = -1
	params.FindAllParams.Size = -1
	length, err := h.UserUsecase.Count(c, params)
	if err != nil {
		err.Path = ".UserHandler->FindAll()" + err.Path
		if err.Error != data.ErrNotFound {
			response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
			return
		}
	}

	dataresponse := types.ResultAll{Status: "Sukses", StatusCode: http.StatusOK, Message: "User Data fetched!", TotalData: length, Page: page, Size: size, Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(h.Status, h.Result)
}

func (h *UserHandler) Find(c *gin.Context) {
	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".EmployeeRoleHandler->Find()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	result, err := h.UserUsecase.Find(c, id)
	if err != nil {
		err.Path = ".UserHandler->Find()" + err.Path
		if err.Error == data.ErrNotFound {
			response.Error(c, "User not found", http.StatusNotFound, *err)
			return
		}
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	result.Password = ""

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "User Data fetched!", Data: result}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *UserHandler) Create(c *gin.Context) {
	var err *types.Error
	var user models.User
	var dataUser *models.User

	hash := md5.New()
	io.WriteString(hash, c.PostForm("Password"))

	user.Name = c.PostForm("Name")
	user.Email, err = library.IsEmailValid(c.PostForm("Email"))
	if err != nil {
		err.Path = ".UserHandler->Create()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}
	user.Username = c.PostForm("Username")
	user.Password = fmt.Sprintf("%x", hash.Sum(nil))

	errJson := json.Unmarshal([]byte(c.PostForm("Permissions")), &user.Permissions)
	if errJson != nil {
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, types.Error{
			Path:  ".UserHandler->Create()",
			Error: errJson,
			Type:  "convert-error",
		})
		return
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		dataUser, err = h.UserUsecase.Create(c, user)
		if err != nil {
			return err
		}

		dataUser.Password = ""

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".UserHandler->Create()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "User Data created!", Data: dataUser}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *UserHandler) Update(c *gin.Context) {
	var err *types.Error
	var user models.User
	var data *models.User

	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".EmployeeRoleHandler->Update()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	user.Name = c.PostForm("Name")
	user.Email, err = library.IsEmailValid(c.PostForm("Email"))
	if err != nil {
		err.Path = ".UserHandler->Update()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}
	user.Username = c.PostForm("Username")

	errJson := json.Unmarshal([]byte(c.PostForm("Permissions")), &user.Permissions)
	if errJson != nil {
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, types.Error{
			Path:  ".UserHandler->Update()",
			Error: errJson,
			Type:  "convert-error",
		})
		return
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.UserUsecase.Update(c, id, user)
		if err != nil {
			return err
		}

		data.Password = ""

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".UserHandler->Update()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "User Data updated!", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	var err *types.Error
	var dataUser *models.User

	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".EmployeeRoleHandler->UpdatePassword()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	var oldPassword = c.PostForm("OldPassword")
	var newPassword = c.PostForm("NewPassword")

	if newPassword == "" {
		err = &types.Error{
			Path:    ".UserHandler->UpdatePassword()",
			Message: "Password baru tidak boleh kosong",
			Type:    "mysql-error",
		}
		response.Error(c, err.Message, http.StatusUnprocessableEntity, *err)
		return
	}

	if newPassword == oldPassword {
		err = &types.Error{
			Path:    ".UserHandler->UpdatePassword()",
			Message: "Password baru tidak boleh sama dengan password lama",
			Type:    "mysql-error",
		}
		response.Error(c, err.Message, http.StatusUnprocessableEntity, *err)
		return
	}

	if newPassword != c.PostForm("ConfirmNewPassword") {
		err = &types.Error{
			Path:    ".UserHandler->UpdatePassword()",
			Message: "Gagal mengkonfirmasi password baru",
			Type:    "mysql-error",
		}
		response.Error(c, err.Message, http.StatusUnprocessableEntity, *err)
		return
	}

	modelUser, err := h.UserUsecase.Find(c, id)
	if err != nil {
		err.Path = ".UserHandler->UpdatePassword()" + err.Path
		if err.Error == data.ErrNotFound {
			response.Error(c, "Data not found", http.StatusUnprocessableEntity, *err)
			return
		}
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	var currentPassword = modelUser.Password
	hash := md5.New()
	io.WriteString(hash, oldPassword)
	hashedOldPassword := fmt.Sprintf("%x", hash.Sum(nil))

	if currentPassword != hashedOldPassword {
		err = &types.Error{
			Path:    ".UserHandler->UpdatePassword()",
			Message: "Incorrect Previous Password",
			Type:    "mysql-error",
		}
		response.Error(c, err.Message, http.StatusUnprocessableEntity, *err)
		return
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		dataUser, err = h.UserUsecase.UpdatePassword(c, id, newPassword)
		if err != nil {
			return err
		}

		dataUser.Password = ""

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".UserHandler->UpdatePassword()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "User Password has been updated!", Data: dataUser}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	var err *types.Error
	var dataUser *models.User

	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".EmployeeRoleHandler->ResetPassword()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		dataUser, err = h.UserUsecase.UpdatePassword(c, id, "123456")
		if err != nil {
			return err
		}

		dataUser.Password = ""

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".UserHandler->ResetPassword()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "User Password reset successful", Data: dataUser}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *UserHandler) FindStatus(c *gin.Context) {
	var datas []*models.Status
	datas = append(datas, &models.Status{ID: models.STATUS_INACTIVE, Name: "Inactive"})
	datas = append(datas, &models.Status{ID: models.STATUS_ACTIVE, Name: "Active"})

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "User Status Data fetched!", Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(http.StatusOK, h.Result)
}

func (h *UserHandler) UpdateStatus(c *gin.Context) {
	var err *types.Error
	var data *models.User

	userID, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".UserRoleHandler->UpdateStatus()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	newStatusID := c.PostForm("StatusID")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.UserUsecase.UpdateStatus(c, userID, newStatusID)
		if err != nil {
			return err
		}

		data.Password = ""

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".UserHandler->UpdateStatus()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "User Status has been updated!", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *UserHandler) Login(c *gin.Context) {
	var err *types.Error
	var obj models.UserLogin
	var data *models.UserLogin

	hash := md5.New()
	io.WriteString(hash, c.PostForm("Password"))

	username := c.PostForm("Username")
	password := fmt.Sprintf("%x", hash.Sum(nil))

	obj.Username = username
	obj.Password = password

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.UserUsecase.Login(c, obj)
		if err != nil {
			return err
		}

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".UserHandler->Login()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Login successful", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}
