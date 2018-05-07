package models

import (
	"crypto/md5"
	"io"
	"os"
	"time"
	"strings"
	"log"
	"strconv"
	"errors"
	"encoding/base64"

	"github.com/muxiyun/MAE/tools/aestool"
)

type User struct {
	Username string
	Uid int
	Role string
	PasswordHash string
	Token string
}

func (u User) GeneratePasswordHash(password string) {
	key := []byte(os.Getenv("SECRETKEY"))
	passrwordhash := md5.New()
	io.WriteString(passrwordhash, password)
	u.PasswordHash = string(passrwordhash.Sum(key))
	return
}

func (u User) Login(password string) (token string, b bool){
	key := []byte(os.Getenv("SECRETKEY"))
	passwordhash := md5.New()
	io.WriteString(passwordhash, password)
	hash := string(passwordhash.Sum(key))
	if u.PasswordHash == hash {
		u.GenerateToken()
		return u.Token, true
	} else {
		return "", false
	}
}


/*base64(aes(uid)) + "." + base64(aes(role)) + "." + base64(aes(generated unix timestamp))*/
func (u User) GenerateToken() {
	key := []byte(os.Getenv("SECRETKEY"))
	goaes := aestool.NewGoAES(key)

	userinfo := []byte(string(u.Uid))
	encrypt1, err := goaes.Encrypt(userinfo)
	if err != nil {
		log.Fatalln(u.Uid, "Generating tokens", err)
	}

	role := []byte(u.Role)
	encrypt2, err := goaes.Encrypt([]byte(role))
	if err != nil {
		log.Fatalln(u.Uid, "Generating tokens", err)
	}

	timestamp :=  []byte(string(time.Now().Unix()))
	encrypt3, err:= goaes.Encrypt(timestamp)
	if err != nil {
		log.Fatalln(u.Uid, "Generating tokens", err)
	}

	token := base64.StdEncoding.EncodeToString(encrypt1) + "." +
		     base64.StdEncoding.EncodeToString(encrypt2) + "." +
		     base64.StdEncoding.EncodeToString(encrypt3)
	u.Token = token
}

func (u User) CheckToken(token string, role string) (err error) {
	// 过期时间 3h 秒
	expire := int64(3 * 60 * 60)

	tokens := strings.Split(token, ".")
	if len(tokens) != 3 {
		err = errors.New("Invalid Token")
		return err
	}

	// check uid
	uidstr, err := base64.StdEncoding.DecodeString(tokens[0])
	if err != nil {
		return err
	}
	uid, err := strconv.Atoi(string(uidstr))
	if err != nil {
		return err
	}
	if uid != u.Uid {
		err = errors.New("Uid Not Match")
		return err
	}

	// check role
	tokenrole, err := base64.StdEncoding.DecodeString(tokens[1])
	if err != nil {
		return err
	}
	if string(tokenrole) != u.Role {
		err = errors.New("User Role Not Match")
		return err
	}

	// check time
	tokentime, err := base64.StdEncoding.DecodeString(tokens[2])
	if err != nil {
		return err
	}

	timestamp, err := strconv.Atoi(string(tokentime))
	if err != nil {
		return err
	}
	result := time.Now().Unix() - int64(timestamp)
	if result > expire || result < 0 {
		err = errors.New("Time Expired")
		return err
	}

	// Token有效， 更新Token 时间
	u.GenerateToken()
	// 错误返回空
	return nil
}