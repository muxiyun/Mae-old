package model

import (
	"github.com/jinzhu/gorm"
)

type Version struct {
	gorm.Model
	ServiceID   uint   `json:"svc_id" gorm:"column:svc_id;not null"`
	VersionName string `json:"version_name" gorm:"column:version_name;type:varchar(50)"`
	VersionDesc string `json:"version_desc" gorm:"column:version_desc;type:varchar(512)"`
}

func (c *Version) TableName() string {
	return "versions"
}
