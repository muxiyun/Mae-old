

package model

import (
	"github.com/jinzhu/gorm"
)

type Service struct {
	gorm.Model
	AppID uint `json:"app_id" gorm:"column:app_id;not null"`
	SvcName string `json:"svc_name" gorm:"column:svc_name;not null;type:varchar(50)"`
	SvcDesc string `json:"svc_desc" gorm:"column:svc_desc;type:varchar(512)"`

	Versions []Version `gorm:"foreignkey:ServiceID"` //service表不会多任何字段，Version表多一个ServiceID
}


func (c *Service) TableName() string {
	return "services"
}