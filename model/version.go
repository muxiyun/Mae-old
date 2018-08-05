package model

import (
	"github.com/jinzhu/gorm"
)

type Version struct {
	gorm.Model
	SvcId uint
	VersionName string
	VersionConf VersionConfig
}
