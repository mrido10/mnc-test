package sqlRepository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"test2/model"
	"test2/model/entity"
)

type User struct {
	SqlDB
}

type Users interface {
	Insert(tx *gorm.DB, user *entity.User) *model.Error
	GetByPhoneNumber(phoneNumber string) (entity.User, *model.Error)
	GetByID(id string) (entity.User, *model.Error)
	UpdateBalance(tx *gorm.DB, id string, balance float64) *model.Error
}

func NewUser(db *gorm.DB) Users {
	return &User{SqlDB{DB: db}}
}

func (u User) Insert(tx *gorm.DB, user *entity.User) *model.Error {
	if err := tx.Create(user).Error; err != nil {
		if strings.Contains(err.Error(), "SQLSTATE 23505") {
			return model.NewError(400, "phone number has been used", nil)
		}
		return model.NewError(500, "internal server error", err)
	}
	return nil
}

func (u User) GetByPhoneNumber(phoneNumber string) (entity.User, *model.Error) {
	var usr entity.User
	if err := u.DB.Where("phone_number = ?", phoneNumber).Find(&usr).Error; err != nil {
		return usr, model.NewError(500, "internal server error", err)
	}
	return usr, nil
}

func (u User) GetByID(id string) (entity.User, *model.Error) {
	var usr entity.User
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return usr, model.NewError(500, "internal server error", err)
	}
	if err := u.DB.Where("id = ?", idUUID).Find(&usr).Error; err != nil {
		return usr, model.NewError(500, "internal server error", err)
	}
	return usr, nil
}

func (u User) UpdateBalance(tx *gorm.DB, id string, balance float64) *model.Error {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return model.NewError(500, "internal server error", err)
	}
	if err := tx.Model(&entity.User{}).
		Where("id = ?", idUUID).
		Updates(map[string]interface{}{
			"balance": balance,
		}).Error; err != nil {
		return model.NewError(500, "internal server error", err)
	}
	return nil
}
