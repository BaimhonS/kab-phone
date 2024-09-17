package scripts

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BaimhonS/kab-phone/configs"
	"github.com/BaimhonS/kab-phone/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"gorm.io/gorm"
)

func MigrateUp() {
	db := configs.ConnectDB()

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("error get database connection : %v", err)
	}

	if err := db.AutoMigrate(
		models.User{},
		models.Phone{},
		models.Cart{},
		models.Item{},
		models.Order{},
	); err != nil {
		log.Fatalf("error migrating database : %v", err)
	}

	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		log.Fatalf("error get mysql driver : %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatalf("error get migrate instance : %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("error migrate up database : %v", err)
	}

	MigrateImage(db)

	log.Println("migrate up success")
}

func MigrateDown(arg string) {
	arg = strings.Replace(arg, "migrate-down", "", 1)

	targetVersion, err := strconv.Atoi(arg)
	if err != nil {
		log.Fatalf("error convert target version : %v", err)
	}

	db := configs.ConnectDB()

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("error getting database connection: %v", err)
	}

	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	if err != nil {
		log.Fatalf("error getting mysql driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatalf("error getting migrate instance: %v", err)
	}

	currentVersion, _, err := m.Version()
	if err != nil {
		log.Fatalf("error getting current migration version: %v", err)
	}

	if targetVersion >= int(currentVersion) {
		log.Printf("target version %d is greater than or equal to current version %d, nothing to do", targetVersion, currentVersion)
		return
	}

	// Calculate the number of steps to roll back
	steps := int(currentVersion) - int(targetVersion)
	if err := m.Steps(-steps); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("error rolling back database to version %d: %v", targetVersion, err)
	}

	log.Printf("rollback database to version %d success", targetVersion)

}

func MigrateImage(db *gorm.DB) {
	filepath.Walk("./data/phones", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("error walking path : %v", err)
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				log.Fatalf("error opening file : %v", err)
			}

			imageData, err := io.ReadAll(file)
			if err != nil {
				log.Fatalf("error reading file : %v", err)
			}

			fileName := strings.ReplaceAll(info.Name(), ".jpg", "")

			if err := db.Model(&models.Phone{}).Where("model_name = ?", fileName).Updates(models.Phone{
				Image: imageData,
			}).Error; err != nil {
				log.Fatalf("error updating image : %v", err)
			}
		}
		return nil
	})
}
