package data

import (
	"os"
	"path/filepath"

	"github.com/reliefeffortslk/protest-tracker-api/pkg/configs"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite" // Sqlite driver based on GGO

	// "github.com/glebarez/sqlite" // Pure go SQLite driver, checkout https://github.com/glebarez/sqlite for details
	"github.com/xo/dburl"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var _db *gorm.DB

// var log = logrus.New()

func getDialector() gorm.Dialector {
	dbEnv := os.Getenv("DATABASE_DSN")
	if dbEnv == "" {
		panic("DATABASE_DSN was not set")
	}
	u, err := dburl.Parse(dbEnv)
	if err != nil {
		panic(err)
	}

	switch u.Scheme {
	case "sqlite":
		path := u.DSN
		if !filepath.IsAbs(path) {
			path = filepath.Join(configs.ProjectRoot(), path)
		}
		return sqlite.Open(path)
	case "postgres":
		return postgres.Open(u.DSN)
	default:
		panic("unsupported database scheme")
	}

}

func GetDb() *gorm.DB {
	if _db != nil {
		return _db
	}

	var err error
	dialector := getDialector()
	_db, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	_db.AutoMigrate(&Protest{}, &ProtestNotification{}, &GeoSubscription{}, &DataSource{})
	return _db
}
