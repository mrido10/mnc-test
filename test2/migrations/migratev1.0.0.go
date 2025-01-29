package migrations

import "test2/model/entity"

func (m migrate) migratev1_0_0() error {
	if err := m.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		return err
	}
	return m.createTableIfNotExistAutoMigrate(
		&entity.User{},
		&entity.Transaction{},
	)
}
