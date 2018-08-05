package model

import (
	"github.com/jinzhu/gorm"
)

type Service struct {
	gorm.Model
	AppId uint
	token string
	SvcName string
	Port uint
	Versions []Version
}
