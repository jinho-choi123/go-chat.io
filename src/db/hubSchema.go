package db

type Hub struct {
	Id       uint64 `gorm:"primaryKey;autoIncrement"`
	Users    []User `gorm:"many2many:user_hubs"`
	Messages []Message
	Created  int64 `gorm:"autoCreateTime"`
}
