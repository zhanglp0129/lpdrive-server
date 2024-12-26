package model

type Model struct {
	ID        int64 `gorm:"primaryKey"`
	CreatedAt int64 `gorm:"autoUpdateTime:milli"`
	UpdatedAt int64 `gorm:"autoCreateTime:milli"`
}
