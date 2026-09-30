package models

import (
	"time"
)

type Kategori struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	NamaKategori string    `gorm:"type:varchar(150);not null" json:"nama_kategori"`
	Slug         string    `gorm:"type:varchar(200);uniqueIndex;not null" json:"slug"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
