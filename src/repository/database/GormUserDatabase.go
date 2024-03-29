package database

import (
	"fmt"
	"github.com/strikersk/user-auth/config"
	"github.com/strikersk/user-auth/constants"
	"github.com/strikersk/user-auth/src/domain"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(configuration config.DatabaseConfiguration) *GormUserRepository {
	dialector := resolveDatabase(configuration)

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatal("Database initialization: ", err.Error())
	}

	if err = db.AutoMigrate(&domain.UserDTO{}); err != nil {
		log.Fatal("error migrating struct: ", err)
	}

	return &GormUserRepository{
		db: db,
	}
}

func (r *GormUserRepository) CreateEntry(user *domain.UserDTO) (err error) {
	if err = r.db.Create(user).Error; err != nil {
		return err
	}
	return
}

func (r *GormUserRepository) ReadEntry(user *domain.UserDTO) (bool, error) {
	err := r.db.Where("username = ?", user.Username).First(user).Error

	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

func resolveDatabase(configuration config.DatabaseConfiguration) gorm.Dialector {
	switch configuration.Type {
	case constants.SQLite:
		return createSQLiteDialector(configuration)
	case constants.Postgres:
		return createPostgresDialector(configuration)
	default:
		log.Fatalf("No persistent database instance selected, available options are: %s, %s", constants.SQLite, constants.Postgres)
		return nil
	}
}

func createSQLiteDialector(configuration config.DatabaseConfiguration) gorm.Dialector {
	if configuration.URL == "" {
		log.Fatal("Database address not selected")
	}
	return sqlite.Open(configuration.URL)
}

func createPostgresDialector(configuration config.DatabaseConfiguration) gorm.Dialector {
	host := configuration.Host
	port := configuration.Port
	dbName := configuration.Name
	username := configuration.Username
	password := configuration.Password
	sslMode := "disable"

	args := fmt.Sprintf("host=%s port=%s dbname=%s user='%s' password=%s sslmode=%s", host, port, dbName, username, password, sslMode)
	return postgres.Open(args)
}
