package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"candyzflorist-go/config"
	"candyzflorist-go/helpers"
	"candyzflorist-go/models"
	"github.com/gin-gonic/gin"
)

// GetPublicKategori returns all categories for public storefront
func GetPublicKategori(c *gin.Context) {
	var categories []models.Kategori
	if err := config.DB.Order("id asc").Find(&categories).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data kategori")
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "Daftar kategori publik", categories)
}

// GetPublicKatalog returns all products with optional filter by category slug, search, and pagination
func GetPublicKatalog(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "12"))
	if limit < 1 {
		limit = 12
	} else if limit > 100 {
		limit = 100
	}

	categorySlug := strings.TrimSpace(c.Query("kategori"))
	search := strings.TrimSpace(c.Query("search"))

	query := config.DB.Model(&models.Produk{})

	if search != "" {
		query = query.Where("nama_produk LIKE ? OR deskripsi LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if categorySlug != "" && categorySlug != "all" {
		query = query.Joins("JOIN produk_kategoris ON produk_kategoris.produk_id = produks.id").
			Joins("JOIN kategoris ON kategoris.id = produk_kategoris.kategori_id").
			Where("kategoris.slug = ?", categorySlug)
	}

	var totalItems int64
	query.Count(&totalItems)

	totalPages := int((totalItems + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	offset := (page - 1) * limit

	var products []models.Produk
	if err := query.Preload("Kategori").Preload("FotoProduk").
		Order("produks.created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&products).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil katalog produk")
		return
	}

	var items []gin.H
	for _, p := range products {
		var photoUrls []string
		for _, f := range p.FotoProduk {
			photoUrls = append(photoUrls, formatImageURL(c, f.FilePath))
		}

		// Category slug list
		var categorySlugs []string
		var categoryLabels []string
		for _, cat := range p.Kategori {
			categorySlugs = append(categorySlugs, cat.Slug)
			categoryLabels = append(categoryLabels, cat.NamaKategori)
		}

		primaryCategory := "ready-stock"
		if len(categorySlugs) > 0 {
			primaryCategory = categorySlugs[0]
		}

		mainImage := ""
		if len(photoUrls) > 0 {
			mainImage = photoUrls[0]
		}

		items = append(items, gin.H{
			"id":               p.ID,
			"name":             p.NamaProduk,
			"nama_produk":      p.NamaProduk,
			"slug":             p.Slug,
			"price":            p.Harga,
			"harga":            p.Harga,
			"deskripsi":        p.Deskripsi,
			"category":         primaryCategory,
			"categories":       categorySlugs,
			"category_labels":  categoryLabels,
			"image":            mainImage,
			"foto_produk":      photoUrls,
		})
	}

	helpers.SuccessResponse(c, http.StatusOK, "Katalog produk publik", gin.H{
		"items": items,
		"pagination": gin.H{
			"current_page": page,
			"per_page":     limit,
			"total_items":  totalItems,
			"total_pages":  totalPages,
		},
	})
}

// GetPublicProdukDetail returns single product detail by slug or ID
func GetPublicProdukDetail(c *gin.Context) {
	param := c.Param("identifier")
	var product models.Produk

	query := config.DB.Preload("Kategori").Preload("FotoProduk")
	if err := query.Where("slug = ? OR id = ?", param, param).First(&product).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusNotFound, "Produk tidak ditemukan")
		return
	}

	var photoUrls []string
	for _, f := range product.FotoProduk {
		photoUrls = append(photoUrls, formatImageURL(c, f.FilePath))
	}

	mainImage := ""
	if len(photoUrls) > 0 {
		mainImage = photoUrls[0]
	}

	helpers.SuccessResponse(c, http.StatusOK, "Detail produk", gin.H{
		"id":          product.ID,
		"name":        product.NamaProduk,
		"nama_produk": product.NamaProduk,
		"slug":        product.Slug,
		"price":       product.Harga,
		"harga":       product.Harga,
		"deskripsi":   product.Deskripsi,
		"kategori":    product.Kategori,
		"image":       mainImage,
		"foto_produk": photoUrls,
	})
}

// GetPublicPengaturan returns public store information (WhatsApp, contacts, address)
func GetPublicPengaturan(c *gin.Context) {
	var pengaturan models.Pengaturan
	if err := config.DB.First(&pengaturan).Error; err != nil {
		pengaturan = models.Pengaturan{
			Whatsapp:   "089688035866",
			Instagram:  "https://instagram.com/crandyzflorist",
			Tiktok:     "https://tiktok.com/@crandyzflorist",
			Email:      "crandyzflorist@gmail.com",
			Alamat:     "Blok F No. 528, Perumahan Bumi Telukjambe, Kec. Telukjambe Timur, Karawang, Jawa Barat 41361",
			TemplateWa: "Halo CandyzFlorist, saya tertarik untuk memesan produk *{nama_produk}* dengan harga *Rp {harga}*.",
		}
	}

	helpers.SuccessResponse(c, http.StatusOK, "Informasi publik toko", pengaturan)
}
