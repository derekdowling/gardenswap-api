package API

import (
	"path"
	"runtime"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

// GetDB loads the db from file
func (a *API) GetDB() (*xorm.Engine, error) {
	_, filename, _, _ := runtime.Caller(0)
	dbPath := path.Join(path.Dir(filename), "./tmp/gorm.db")

	engine, err := xorm.NewEngine("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	engine.Logger = a.Logger
	return engine, nil
}
