package sqlRepository

import (
	"gorm.io/gorm"
	"test2/model"
	"test2/model/entity"
)

type Transactions interface {
	Insert(tx *gorm.DB, tr *entity.Transaction) *model.Error
}
type Transaction struct {
	SqlDB
}

func NewTransaction(db *gorm.DB) Transactions {
	return Transaction{
		SqlDB: SqlDB{DB: db},
	}
}

func (t Transaction) Insert(tx *gorm.DB, tr *entity.Transaction) *model.Error {
	if err := tx.Create(tr).Error; err != nil {
		return model.NewError(500, "internal server error", err)
	}
	return nil
}
