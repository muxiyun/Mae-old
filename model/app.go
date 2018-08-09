package model

import (
	"github.com/jinzhu/gorm"
)

type App struct {
	gorm.Model
	AppName string `json:"app_name" gorm:"column:app_name;not null;unique;type:varchar(50)"`
	AppDesc string `json:"app_desc" gorm:"column:app_desc;type:varchar(512)"`

	Services []Service `gorm:"foreignkey:AppID"` //App表不会多任何字段，Service表会多app_id
}

func (a *App) TableName() string {
	return "apps"
}

// Create creates a new App.
func (app *App) Create() error {
	return DB.RWdb.Create(&app).Error
}

// DeleteApp deletes the app by the user identifier.
func DeleteApp(id uint) error {
	app:= App{}
	app.ID = id
	return DB.RWdb.Delete(&app).Error
}

// Update updates an App information.
func (app *App) Update() error {
	return DB.RWdb.Save(app).Error
}

func GetAppByName(appname string) (*App, error) {
	app := &App{}
	d := DB.RWdb.Where("app_name = ?", appname).First(&app)
	return app, d.Error
}

func GetAppByID(id int64) (*App, error) {
	app := &App{}
	d := DB.RWdb.Where("id = ?", id).First(&app)
	return app, d.Error
}

// ListApp List all apps
func ListApp(offset, limit int) ([]*App, uint64, error) {

	apps := make([]*App, 0)
	var count uint64
	if err := DB.RWdb.Model(&App{}).Count(&count).Error; err != nil {
		return apps, count, err
	}

	if err := DB.RWdb.Offset(offset).Limit(limit).Order("id desc").Find(&apps).Error; err != nil {
		return apps, count, err
	}

	return apps, count, nil
}
