package databases

import (
	"database/sql"
	"log"

	rice "github.com/GeertJohan/go.rice"
	"github.com/fritz-immanuel/erajaya-be-tech-test/configs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
)

// MigrateUp migrates the database up
func MigrateUp() {
	// Setup the database
	//
	cfg, err := configs.GetConfiguration()
	if err != nil {
		log.Fatal("error when getting configuration: ", err)
	}

	db, err := sql.Open("mysql", cfg.DBConnectionString)
	if err != nil {
		log.Fatal("error when open postgres connection: ", err)
	}

	// Setup the source driver
	//
	sourceDriver := &RiceBoxSource{}
	sourceDriver.PopulateMigrations(rice.MustFindBox("./migrations"))
	if err != nil {
		log.Fatal("error when creating source driver: ", err)
	}

	// Setup the database driver
	//
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal("error when creating postgres instance: ", err)
	}

	m, err := migrate.NewWithInstance(
		"go.rice", sourceDriver,
		"mysql", driver)

	if err != nil {
		log.Fatal("error when creating database instance: ", err)
	}

	if err := m.Up(); err != nil {
		if err.Error() != "no change" {
			log.Fatal("error when migrate up: ", err)
		}
	}

	defer m.Close()
}
