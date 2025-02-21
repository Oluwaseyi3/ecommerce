package models

import (
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	UserID  uint `json:"user_id"`
	IsAdmin bool `json:"is_admin"`
	jwt.RegisteredClaims
}

type User struct {
	ID       uint   `gorm:"primaryKey`
	Email    string `gorm:"unique"`
	Password string
	Name     string
	Address  string
	IsAdmin  bool `gorm::"default:false"`
}

func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	u.Password = string(bytes)
	return err
}

func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
