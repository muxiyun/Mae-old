package model

import (
	"github.com/jinzhu/gorm"
)

type Service struct {
	gorm.Model
	AppID   uint   `json:"app_id" gorm:"column:app_id;not null"`
	SvcName string `json:"svc_name" gorm:"column:svc_name;not null;unique;type:varchar(50)"`
	SvcDesc string `json:"svc_desc" gorm:"column:svc_desc;type:varchar(512)"`

	Versions []Version `gorm:"foreignkey:ServiceID"` //service表不会多任何字段，Version表多一个ServiceID
}

func (c *Service) TableName() string {
	return "services"
}

// Create creates a new Service.
func (svc *Service) Create() error {
	return DB.RWdb.Create(&svc).Error
}

// DeleteService deletes the service by the service id.
func DeleteService(id uint) error {
	svc := Service{}
	svc.ID = id
	return DB.RWdb.Delete(&svc).Error
}

// Update updates a Service information.
func (svc *Service) Update() error {
	return DB.RWdb.Save(svc).Error
}

func GetServiceByName(svc_name string) (*Service, error) {
	svc := &Service{}
	d := DB.RWdb.Where("svc_name = ?", svc_name).First(&svc)
	return svc, d.Error
}

func GetServiceByID(id int64) (*Service, error) {
	svc := &Service{}
	d := DB.RWdb.Where("id = ?", id).First(&svc)
	return svc, d.Error
}

// ListService List all services
func ListService(offset, limit int) ([]*Service, uint64, error) {

	svcs := make([]*Service, 0)
	var count uint64
	if err := DB.RWdb.Model(&Service{}).Count(&count).Error; err != nil {
		return svcs, count, err
	}

	if err := DB.RWdb.Offset(offset).Limit(limit).Order("id desc").Find(&svcs).Error; err != nil {
		return svcs, count, err
	}

	return svcs, count, nil
}

// List　all service belongs to an app
func ListServiceByAppID(offset,limit int,app_id uint)([]*Service, uint64, error){
	svcs := make([]*Service, 0)
	var count uint64
	if err := DB.RWdb.Model(&Service{}).Where("app_id = ?",app_id).Count(&count).Error; err != nil {
		return svcs, count, err
	}

	if err := DB.RWdb.Where("app_id = ?",app_id).Offset(offset).Limit(limit).Order("id desc").Find(&svcs).Error; err != nil {
		return svcs, count, err
	}

	return svcs, count, nil
}