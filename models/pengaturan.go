package models

import (
	"time"
)

type Pengaturan struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Whatsapp   string    `gorm:"type:varchar(30)" json:"whatsapp"`
	Instagram  string    `gorm:"type:varchar(255)" json:"instagram"`
	Tiktok     string    `gorm:"type:varchar(255)" json:"tiktok"`
	Email      string    `gorm:"type:varchar(150)" json:"email"`
	Alamat     string    `gorm:"type:text" json:"alamat"`
	TemplateWa string    `gorm:"type:text" json:"template_wa"`
	UpdatedAt  time.Time `json:"updated_at"`
}
