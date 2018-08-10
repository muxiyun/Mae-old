package model

import (
	"github.com/jinzhu/gorm"
)

type Version struct {
	gorm.Model
	ServiceID   uint   `json:"svc_id" gorm:"column:svc_id;not null"`
	VersionName string `json:"version_name" gorm:"column:version_name;type:varchar(50)"`
	VersionDesc string `json:"version_desc" gorm:"column:version_desc;type:varchar(512)"`
	VersionConfig string `json:"version_conf" gorm:"column:version_conf;type:varchar(4096)"`
}

//VersionConf字段包含了配置结构序列化之后的内容

func (c *Version) TableName() string {
	return "versions"
}


// Create creates a new Version.
func (v *Version) Create() error {
	return DB.RWdb.Create(&v).Error
}

// DeleteVersion deletes the version by the version id.
func DeleteVersion(id uint) error {
	v := Version{}
	v.ID = id
	return DB.RWdb.Delete(&v).Error
}

// Update updates a Version information.
func (v *Version) Update() error {
	return DB.RWdb.Save(v).Error
}

func GetVersionByName(version_name string) (*Version, error) {
	v := &Version{}
	d := DB.RWdb.Where("version_name = ?", version_name).First(&v)
	return v, d.Error
}

func GetVersionByID(id int64) (*Version, error) {
	v := &Version{}
	d := DB.RWdb.Where("id = ?", id).First(&v)
	return v, d.Error
}

// ListVersion List all versions
func ListVersion(offset, limit int) ([]*Version, uint64, error) {

	vs := make([]*Version, 0)
	var count uint64
	if err := DB.RWdb.Model(&Version{}).Count(&count).Error; err != nil {
		return vs, count, err
	}

	if err := DB.RWdb.Offset(offset).Limit(limit).Order("id desc").Find(&vs).Error; err != nil {
		return vs, count, err
	}

	return vs, count, nil
}

// List　all versions belongs to a service
func ListVersionByServiceID(offset,limit int,service_id uint)([]*Version, uint64, error){
	vs := make([]*Version, 0)
	var count uint64
	if err := DB.RWdb.Model(&Version{}).Where("svc_id = ?",service_id).Count(&count).Error; err != nil {
		return vs, count, err
	}

	if err := DB.RWdb.Where("svc_id = ?",service_id).Offset(offset).Limit(limit).Order("id desc").Find(&vs).Error; err != nil {
		return vs, count, err
	}

	return vs, count, nil
}