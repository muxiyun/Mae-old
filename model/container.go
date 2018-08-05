package model

import (
	"github.com/jinzhu/gorm"
)

type Container struct{
	gorm.Model
	ContainerName string
	ImageURL string
	ImageTag string
	RunCmd string
	Port uint
	Envs []Env
	Volumes []Volume
}
