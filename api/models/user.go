package models

type User struct {
	Id       int64  `gorm:"primarykey" json:"id"`
	Name     string `gorm:"not null" json:"name"`
	Username string `gorm:"not null;unique" json:"username"`
	Email    string `gorm:"not null;unique" json:"email"`
	Password string `gorm:"not null" json:"password"`
	IsAdmin  bool
}
