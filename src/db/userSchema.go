package db

type User struct {
	Id           uint64 `gorm:"primaryKey;autoIncrement"`
	Username     string `gorm:"size:256"`
	Nickname     string `gorm:"size:256"`
	PasswordHash string
	PasswordSalt string
	Created      int64  `gorm:"autoCreateTime"`
	Hubs         []*Hub `gorm:"many2many:user_hubs"`
	Messages     []Message
}
