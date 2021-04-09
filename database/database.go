package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// This is list of external variables that are
// accessible from packages importing this one
var (
	DBConn *gorm.DB
)
