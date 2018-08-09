package model

import (
	"github.com/jinzhu/gorm"
)

type Version struct {
	gorm.Model
	ServiceID   uint   `json:"svc_id" gorm:"column:svc_id;not null"`
	VersionName string `json:"version_name" gorm:"column:version_name;type:varchar(50)"`
	VersionDesc string `json:"version_desc" gorm:"column:version_desc;type:varchar(512)"`
	VersionConf string `json:"version_conf" gorm:"column:version_conf;type:varchar(4096)"`
}

//VersionConf字段包含了配置结构序列化之后的内容

func (c *Version) TableName() string {
	return "versions"
}
