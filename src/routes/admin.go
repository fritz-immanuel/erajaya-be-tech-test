package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/fritz-immanuel/erajaya-be-tech-test/library/data"
	"github.com/fritz-immanuel/erajaya-be-tech-test/src/app/admin"
	"github.com/jmoiron/sqlx"
)

// RegisterWebRoutes  is a function to register all WEB Routes in the projectbase
func RegisterAdminRoutes(db *sqlx.DB, dataManager *data.Manager, router *gin.Engine) {
	v1 := router.Group("/admin/v1")
	{
		admin.RegisterRoutes(db, dataManager, router, v1)
	}
}
