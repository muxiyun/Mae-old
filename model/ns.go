package model

import "github.com/jinzhu/gorm"

type NameSpace struct{
	gorm.Model
	NSName string

}