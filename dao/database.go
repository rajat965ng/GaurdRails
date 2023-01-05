package dao

import (
	"findings/model"
	"fmt"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	err  error
	lock = sync.Mutex{}
)

type IDatabase interface {
	GetConnection() *gorm.DB
}

type database struct {
}

func NewDatabaseInstance() *database {
	return &database{}
}

func (dbConfig *database) GetConnection() *gorm.DB {
	if db == nil {
		lock.Lock()
		defer lock.Unlock()
		if db == nil {
			fmt.Println("Creating db connection !!")
			db, err = gorm.Open(sqlite.Open("./resources/dev/test.db"), &gorm.Config{})
			if err != nil {
				panic("failed to connect database")
			}
			db.AutoMigrate(&model.Repository{}, &model.ScanDetail{}, &model.Finding{})
		}
	} else {
		fmt.Println("Using already created db connection !!")
	}
	return db
}
