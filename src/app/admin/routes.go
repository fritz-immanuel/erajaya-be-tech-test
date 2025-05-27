package admin

import (
	http_product "github.com/fritz-immanuel/erajaya-be-tech-test/src/app/admin/product"
	http_user "github.com/fritz-immanuel/erajaya-be-tech-test/src/app/admin/user"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library/data"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

var (
	productHandler http_product.ProductHandler
	userHandler    http_user.UserHandler
)

func RegisterRoutes(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine, v *gin.RouterGroup) {
	v1 := v.Group("")
	{
		productHandler.RegisterAPI(db, dataManager, router, v1)
		userHandler.RegisterAPI(db, dataManager, router, v1)
	}
}
