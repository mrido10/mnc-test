package sqlRepository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"test2/model"
	"test2/model/entity"
	"test2/util"
)

type Transactions interface {
	Insert(tx *gorm.DB, tr *entity.Transaction) *model.Error
	GetList(userID string, clause Clause) ([]entity.Transaction, *model.Error)
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

func (t Transaction) GetList(userID string, cl Clause) ([]entity.Transaction, *model.Error) {
	t.DB = t.DB.Limit(cl.Limit).Offset(cl.Offset)
	if util.IsEmptyStringWithTrimSpace(&cl.Order) {
		t.DB = t.DB.Order(cl.Order)
	}

	idUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, model.NewError(500, "internal server error", err)
	}
	var list []entity.Transaction
	if err := t.DB.Model(&entity.Transaction{}).
		Where("user_id = ?", idUUID).
		Find(&list).Error; err != nil {
		return list, model.NewError(500, "internal server error", err)
	}
	return list, nil
}
