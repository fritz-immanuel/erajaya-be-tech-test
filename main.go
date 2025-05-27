package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fritz-immanuel/erajaya-be-tech-test/configs"
	"github.com/fritz-immanuel/erajaya-be-tech-test/databases"
	"github.com/fritz-immanuel/erajaya-be-tech-test/library/data"
	"github.com/fritz-immanuel/erajaya-be-tech-test/src/routes"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	os.Setenv("TZ", "Asia/Jakarta")

	config, err := configs.GetConfiguration()
	if err != nil {
		log.Fatalln("failed to get configuration: ", err)
	}

	configs.AppConfig = config

	db, err := sqlx.Open("mysql", config.DBConnectionString)
	if err != nil {
		log.Fatalln("failed to open database x: ", err)
	}
	defer db.Close()

	dataManager := data.NewManager(
		db,
	)

	databases.MigrateUp()

	fmt.Println("Server is running...")
	routes.RegisterRoutes(db, config, dataManager)
}
