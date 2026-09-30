package handlers

import (
	"net/http"

	"candyzflorist-go/config"
	"candyzflorist-go/helpers"
	"candyzflorist-go/models"
	"github.com/gin-gonic/gin"
)

type PengaturanInput struct {
	Whatsapp   string `json:"whatsapp"`
	Instagram  string `json:"instagram"`
	Tiktok     string `json:"tiktok"`
	Email      string `json:"email"`
	Alamat     string `json:"alamat"`
	TemplateWa string `json:"template_wa"`
}

func GetPengaturan(c *gin.Context) {
	var pengaturan models.Pengaturan
	if err := config.DB.First(&pengaturan).Error; err != nil {
		// Buat row default jika belum ada
		pengaturan = models.Pengaturan{
			Whatsapp:   "089688035866",
			Instagram:  "https://instagram.com/crandyzflorist",
			Tiktok:     "https://tiktok.com/@crandyzflorist",
			Email:      "crandyzflorist@gmail.com",
			Alamat:     "Blok F No. 528, Perumahan Bumi Telukjambe, Kec. Telukjambe Timur, Karawang, Jawa Barat 41361",
			TemplateWa: "Halo CandyzFlorist, saya tertarik untuk memesan produk *{nama_produk}* dengan harga *Rp {harga}*.",
		}
		config.DB.Create(&pengaturan)
	}

	helpers.SuccessResponse(c, http.StatusOK, "Data pengaturan sistem", pengaturan)
}

func UpdatePengaturan(c *gin.Context) {
	var input PengaturanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Format data tidak valid")
		return
	}

	var pengaturan models.Pengaturan
	if err := config.DB.First(&pengaturan).Error; err != nil {
		pengaturan = models.Pengaturan{}
		config.DB.Create(&pengaturan)
	}

	pengaturan.Whatsapp = input.Whatsapp
	pengaturan.Instagram = input.Instagram
	pengaturan.Tiktok = input.Tiktok
	pengaturan.Email = input.Email
	pengaturan.Alamat = input.Alamat
	pengaturan.TemplateWa = input.TemplateWa

	if err := config.DB.Save(&pengaturan).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal menyimpan pengaturan")
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "Pengaturan berhasil diperbarui", pengaturan)
}
