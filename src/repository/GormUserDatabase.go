package userRepository

import (
	"github.com/strikersk/user-auth/src/domain"
	"gorm.io/gorm"
	"log"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) GormUserRepository {
	if err := db.AutoMigrate(&domain.UserDTO{}); err != nil {
		log.Fatal("error migrating struct: ", err)
	}

	return GormUserRepository{
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
