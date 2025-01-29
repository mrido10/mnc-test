package sqlRepository

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type Clause struct {
	Limit  int
	Offset int
	Order  string
}

func SqlConnection(addr string, maxIdle, maxOpenConn int, maxLifeTimeConn int64) *gorm.DB {
	db, err := gorm.Open(postgres.Open(addr), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(maxLifeTimeConn) * time.Hour)

	return db
}

type SqlDB struct {
	DB *gorm.DB
}

func NewSql(db *gorm.DB) SqlDB {
	return SqlDB{DB: db}
}

func (s SqlDB) DBBeginTX() *gorm.DB {
	return s.DB.Begin()
}
