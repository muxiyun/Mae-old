package model

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"
	"golang.org/x/crypto/bcrypt"
)



type User struct {
	gorm.Model
	UserName string `json:"username" gorm:"column:username;not null;unique"`
	Email string `json:"email" gorm:"column:email;not null;unique"`
	PasswordHash string `json:"password" gorm:"column:passwordhash;not null"`
	Role string `json:"role" gorm:"column:role;default:'user'"`
}


func (c *User) TableName() string {
	return "users"
}

// Create creates a new user account.
func (u *User) Create() error {
	return DB.RWdb.Create(&u).Error
}

// DeleteUser deletes the user by the user identifier.
func DeleteUser(id uint) error {
	user := User{}
	user.ID = id
	return DB.RWdb.Delete(&user).Error
}

// Update updates an user account information.
func (u *User) Update() error {
	return DB.RWdb.Save(u).Error
}

// GetUser gets an user by the user name.
func GetUserByName(username string) (*User, error) {
	u := &User{}
	d := DB.RWdb.Where("username = ?", username).First(&u)
	return u, d.Error
}

// GetUser gets an user by the user's email.
func GetUserByEmail(email string)(*User,error) {
	u := &User{}
	d := DB.RWdb.Where("email = ?", email).First(&u)
	return u, d.Error
}

//GetUser gets an user by the user id
func GetUserByID(id uint)(*User,error){
	u:=&User{}
	d:=DB.RWdb.Where("id = ?",id).First(&u)
	return u,d.Error
}

// ListUser List all users
func ListUser(offset, limit int) ([]*User, uint64, error) {

	users := make([]*User, 0)
	var count uint64

	if err := DB.RWdb.Model(&User{}).Count(&count).Error; err != nil {
		return users, count, err
	}

	if err := DB.RWdb.Offset(offset).Limit(limit).Order("id desc").Find(&users).Error; err != nil {
		return users, count, err
	}

	return users, count, nil
}

// Compare with the plain text password. Returns true if it's the same as the encrypted one (in the `User` struct).
func (u *User) Compare(pwd string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(pwd))
}

// Encrypt the user password.
func (u *User) Encrypt() (err error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(u.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash=string(hashedBytes)
	return nil
}

// Validate the fields.
func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}