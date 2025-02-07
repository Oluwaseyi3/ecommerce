package models

type User struct {
	ID       uint   `gorm:"primaryKey`
	Email    string `gorm:"string"`
	Password string
	Name     string
	Address  string
	IsAdmin  bool `gorm::"default:false"`
}

// func (u )()  {

// }
