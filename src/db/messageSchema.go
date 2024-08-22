package db

type Message struct {
	Id      uint64 `gorm:"primaryKey;autoIncrement"`
	Content string
	Hash    string
	Created int64 `gorm:"autoCreateTime"`
	UserID  uint64
	HubID   uint64
}
