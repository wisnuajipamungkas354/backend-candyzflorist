package models

import (
	"time"
)

type Produk struct {
	ID          uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	NamaProduk  string        `gorm:"type:varchar(200);not null" json:"nama_produk"`
	Slug        string        `gorm:"type:varchar(250);uniqueIndex;not null" json:"slug"`
	Harga       float64       `gorm:"type:decimal(15,2);not null;default:0" json:"harga"`
	Deskripsi   string        `gorm:"type:text" json:"deskripsi"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`

	// Relasi
	Kategori    []Kategori    `gorm:"many2many:produk_kategoris;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"kategori"`
	FotoProduk  []ProdukImage `gorm:"foreignKey:ProdukID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"foto_produk"`
}
