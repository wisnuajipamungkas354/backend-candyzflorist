package handlers

import (
	"net/http"

	"candyzflorist-go/config"
	"candyzflorist-go/helpers"
	"candyzflorist-go/models"
	"github.com/gin-gonic/gin"
)

func GetDashboardStats(c *gin.Context) {
	var totalProduk int64
	var totalKategori int64

	config.DB.Model(&models.Produk{}).Count(&totalProduk)
	config.DB.Model(&models.Kategori{}).Count(&totalKategori)

	var pengaturan models.Pengaturan
	nomorWa := "089688035866"
	if err := config.DB.First(&pengaturan).Error; err == nil && pengaturan.Whatsapp != "" {
		nomorWa = pengaturan.Whatsapp
	}

	var recentProducts []models.Produk
	config.DB.Preload("FotoProduk").Preload("Kategori").Order("created_at desc").Limit(4).Find(&recentProducts)

	var recentList []gin.H
	for _, p := range recentProducts {
		var photoUrl string
		if len(p.FotoProduk) > 0 {
			photoUrl = formatImageURL(c, p.FotoProduk[0].FilePath)
		}
		recentList = append(recentList, gin.H{
			"id":          p.ID,
			"nama_produk": p.NamaProduk,
			"harga":       p.Harga,
			"foto":        photoUrl,
			"created_at":  p.CreatedAt,
		})
	}

	helpers.SuccessResponse(c, http.StatusOK, "Statistik Dashboard", gin.H{
		"total_produk":   totalProduk,
		"total_kategori": totalKategori,
		"nomor_wa":       nomorWa,
		"produk_terbaru": recentList,
	})
}
