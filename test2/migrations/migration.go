package migrations

import (
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type migrate struct {
	*gorm.DB
}

func NewMigrate(db *gorm.DB) migrate {
	return migrate{DB: db.Begin()}
}

func (m migrate) Migrate() error {
	var err error
	defer func() {
		if err != nil {
			m.Rollback()
			return
		}
		m.Commit()
		return
	}()
	err = m.migratev1_0_0()
	return err
}

func (m migrate) createTableIfNotExistAutoMigrate(table ...interface{}) error {
	for _, tbl := range table {
		if !m.Migrator().HasTable(&tbl) {
			if err := m.AutoMigrate(&tbl); err != nil {
				return err
			}
			log.Debug("Create table", tbl)
		}
	}
	return nil
}
