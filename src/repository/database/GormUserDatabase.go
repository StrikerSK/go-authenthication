package database

import (
	"github.com/strikersk/user-auth/src/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository() *GormUserRepository {
	dialector := sqlite.Open("./test.db")

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
	result := r.db.Where("username = ?", user.Username).First(&user)

	if result.Error == nil {
		return true, nil
	} else if result.Error.Error() == gorm.ErrRecordNotFound.Error() {
		return false, nil
	} else {
		return false, result.Error
	}
}
