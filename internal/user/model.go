package user

type User struct {
	ID           int64  `gorm:"primaryKey"`
	Username     string `gorm:"unique;not null"`
	Email        string `gorm:"size:255;unique;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	Role         string `gorm:"column:user_role;size:20;not null"`
}
