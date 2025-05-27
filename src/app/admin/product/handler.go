package product

import (
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library/helpers"
	"github.com/fritz-immanuel/erajaya-be-tech-test/middleware"
	"github.com/fritz-immanuel/erajaya-be-tech-test/models"
	"github.com/fritz-immanuel/erajaya-be-tech-test/src/services/product"
	"github.com/gin-gonic/gin"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library/data"
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/http/response"
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/types"

	productRepository "github.com/fritz-immanuel/erajaya-be-tech-test/src/services/product/repository"
	productUsecase "github.com/fritz-immanuel/erajaya-be-tech-test/src/services/product/usecase"
)

type ProductHandler struct {
	ProductUsecase product.Usecase
	dataManager    *data.Manager
	Result         gin.H
	Status         int
}

func (h ProductHandler) RegisterAPI(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine, v *gin.RouterGroup) {
	productRepo := productRepository.NewProductRepository(
		data.NewMySQLStorage(db, "products", models.Product{}, data.MysqlConfig{}),
		data.NewMySQLStorage(db, "status", models.Status{}, data.MysqlConfig{}),
	)

	uProduct := productUsecase.NewProductUsecase(db, productRepo)

	base := &ProductHandler{ProductUsecase: uProduct, dataManager: dataManager}

	rs := v.Group("/products")
	{
		rs.GET("", middleware.Auth, base.FindAll)
		rs.GET("/:id", middleware.Auth, base.Find)
		rs.POST("", middleware.Auth, base.Create)
		rs.PUT("/:id", middleware.Auth, base.Update)

		rs.PUT("/:id/status", middleware.Auth, base.UpdateStatus)
	}

	rss := v.Group("/statuses")
	{
		rss.GET("/products", base.FindStatus)
	}
}

func (h *ProductHandler) FindAll(c *gin.Context) {
	var params models.FindAllProductParams
	page, size := helpers.FilterFindAll(c)
	filterFindAllParams := helpers.FilterFindAllParam(c)
	params.FindAllParams = filterFindAllParams
	if c.Query("SortName") == "" || c.Query("SortBy") == "" {
		params.FindAllParams.SortBy = "products.name ASC"
	}
	datas, err := h.ProductUsecase.FindAll(c, params)
	if err != nil {
		if err.Error != data.ErrNotFound {
			response.Error(c, err.Message, err.StatusCode, *err)
			return
		}
	}

	params.FindAllParams.Page = -1
	params.FindAllParams.Size = -1
	length, err := h.ProductUsecase.Count(c, params)
	if err != nil {
		err.Path = ".ProductHandler->FindAll()" + err.Path
		if err.Error != data.ErrNotFound {
			response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
			return
		}
	}

	dataresponse := types.ResultAll{Status: "Sukses", StatusCode: http.StatusOK, Message: "Product Data fetched!", TotalData: length, Page: page, Size: size, Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(h.Status, h.Result)
}

func (h *ProductHandler) Find(c *gin.Context) {
	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".ProductHandler->Find()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	result, err := h.ProductUsecase.Find(c, id)
	if err != nil {
		err.Path = ".ProductHandler->Find()" + err.Path
		if err.Error == data.ErrNotFound {
			response.Error(c, "Product not found", http.StatusUnprocessableEntity, *err)
			return
		}
		response.Error(c, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Product Data fetched!", Data: result}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var err *types.Error
	var product models.Product
	var dataProduct *models.Product

	product.Name = c.PostForm("Name")
	product.Price, _ = strconv.ParseFloat(c.PostForm("Price"), 64)
	product.Description = c.PostForm("Description")
	product.Quantity, _ = strconv.Atoi(c.PostForm("Quantity"))

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		dataProduct, err = h.ProductUsecase.Create(c, product)
		if err != nil {
			return err
		}

		return nil
	})
	if errTransaction != nil {
		errTransaction.Path = ".ProductHandler->Create()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Product Data created!", Data: dataProduct}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *ProductHandler) Update(c *gin.Context) {
	var err *types.Error
	var product models.Product
	var data *models.Product

	id, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".ProductHandler->Update()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	product.Name = c.PostForm("Name")
	product.Price, _ = strconv.ParseFloat(c.PostForm("Price"), 64)
	product.Description = c.PostForm("Description")
	product.Quantity, _ = strconv.Atoi(c.PostForm("Quantity"))

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.ProductUsecase.Update(c, id, product)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".ProductHandler->Update()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Product Data updated!", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}

func (h *ProductHandler) FindStatus(c *gin.Context) {
	var datas []*models.Status
	datas = append(datas, &models.Status{ID: models.STATUS_INACTIVE, Name: "Inactive"})
	datas = append(datas, &models.Status{ID: models.STATUS_ACTIVE, Name: "Active"})

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Product Status Data fetched!", Data: datas}
	h.Result = gin.H{
		"result": dataresponse,
	}
	c.JSON(http.StatusOK, h.Result)
}

func (h *ProductHandler) UpdateStatus(c *gin.Context) {
	var err *types.Error
	var data *models.Product

	productID, err := helpers.ValidateUUID(c.Param("id"))
	if err != nil {
		err.Path = ".ProductHandler->UpdateStatus()" + err.Path
		response.Error(c, err.Message, err.StatusCode, *err)
		return
	}

	newStatusID := c.PostForm("StatusID")

	errTransaction := h.dataManager.RunInTransaction(c, func(tctx *gin.Context) *types.Error {
		data, err = h.ProductUsecase.UpdateStatus(c, productID, newStatusID)
		if err != nil {
			return err
		}

		return nil
	})

	if errTransaction != nil {
		errTransaction.Path = ".ProductHandler->UpdateStatus()" + errTransaction.Path
		response.Error(c, errTransaction.Message, errTransaction.StatusCode, *errTransaction)
		return
	}

	dataresponse := types.Result{Status: "Sukses", StatusCode: http.StatusOK, Message: "Product Status has been updated!", Data: data}
	h.Result = gin.H{
		"result": dataresponse,
	}

	c.JSON(http.StatusOK, h.Result)
}
