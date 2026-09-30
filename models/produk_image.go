package models

import (
	"time"
)

type ProdukImage struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProdukID  uint      `gorm:"index;not null" json:"produk_id"`
	FilePath  string    `gorm:"type:varchar(500);not null" json:"file_path"`
	CreatedAt time.Time `json:"created_at"`
}
