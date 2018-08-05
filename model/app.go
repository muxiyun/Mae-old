package model

import (
	"github.com/jinzhu/gorm"
)
type App struct {
	gorm.Model
	AppName string
	SvcList []Service
}
