package model

import (
	"github.com/jinzhu/gorm"
)

type VersionConfig struct {
	gorm.Model
	Replicas uint
	Containers []Container
}
