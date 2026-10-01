package handlers

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"candyzflorist-go/config"
	"candyzflorist-go/helpers"
	"candyzflorist-go/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const MaxImageSizeBytes = 2 * 1024 * 1024 // 2MB

// Helper to check valid image extensions
func isValidImageExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp"
}

// Helper to format full image URL
func formatImageURL(c *gin.Context, filePath string) string {
	if strings.HasPrefix(filePath, "http://") || strings.HasPrefix(filePath, "https://") {
		return filePath
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	host := c.Request.Host
	return fmt.Sprintf("%s://%s/%s", scheme, host, strings.TrimPrefix(filePath, "./"))
}

// Save uploaded file into uploads/<category-slug>/<uuid><ext>
func saveUploadedImage(file *multipart.FileHeader, categorySlug string, uploadDir string) (string, error) {
	if file.Size > MaxImageSizeBytes {
		return "", fmt.Errorf("ukuran file %s melebihi batas 2MB", file.Filename)
	}

	if !isValidImageExtension(file.Filename) {
		return "", fmt.Errorf("format file %s harus berupa JPG, PNG, atau WEBP", file.Filename)
	}

	if categorySlug == "" {
		categorySlug = "general"
	}

	targetDir := filepath.Join(uploadDir, categorySlug)
	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("gagal membuat direktori upload: %v", err)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	savePath := filepath.Join(targetDir, newFileName)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := dst.ReadFrom(src); err != nil {
		return "", err
	}

	// Normalisasi path separator untuk web URL
	webPath := filepath.ToSlash(savePath)
	return webPath, nil
}

func GetAllProduk(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 {
		limit = 10
	} else if limit > 100 {
		limit = 100
	}

	search := strings.TrimSpace(c.Query("search"))
	kategoriIDStr := strings.TrimSpace(c.Query("kategori_id"))

	query := config.DB.Model(&models.Produk{})

	if search != "" {
		query = query.Where("nama_produk LIKE ? OR deskripsi LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if kategoriIDStr != "" && kategoriIDStr != "ALL" {
		if kID, err := strconv.ParseUint(kategoriIDStr, 10, 32); err == nil {
			query = query.Joins("JOIN produk_kategoris ON produk_kategoris.produk_id = produks.id").
				Where("produk_kategoris.kategori_id = ?", kID)
		}
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
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&products).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data produk")
		return
	}

	// Transform data to ensure full URLs
	var items []gin.H
	for _, p := range products {
		var photoUrls []string
		for _, f := range p.FotoProduk {
			photoUrls = append(photoUrls, formatImageURL(c, f.FilePath))
		}

		items = append(items, gin.H{
			"id":          p.ID,
			"nama_produk": p.NamaProduk,
			"slug":        p.Slug,
			"harga":       p.Harga,
			"deskripsi":   p.Deskripsi,
			"kategori":    p.Kategori,
			"foto_produk": photoUrls,
			"created_at":  p.CreatedAt,
			"updated_at":  p.UpdatedAt,
		})
	}

	helpers.SuccessResponse(c, http.StatusOK, "Daftar produk katalog", gin.H{
		"items": items,
		"pagination": gin.H{
			"current_page": page,
			"per_page":     limit,
			"total_items":  totalItems,
			"total_pages":  totalPages,
		},
	})
}

func GetProdukByID(c *gin.Context) {
	id := c.Param("id")
	var product models.Produk

	if err := config.DB.Preload("Kategori").Preload("FotoProduk").First(&product, id).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusNotFound, "Produk tidak ditemukan")
		return
	}

	var photoUrls []string
	for _, f := range product.FotoProduk {
		photoUrls = append(photoUrls, formatImageURL(c, f.FilePath))
	}

	helpers.SuccessResponse(c, http.StatusOK, "Detail produk", gin.H{
		"id":          product.ID,
		"nama_produk": product.NamaProduk,
		"slug":        product.Slug,
		"harga":       product.Harga,
		"deskripsi":   product.Deskripsi,
		"kategori":    product.Kategori,
		"foto_produk": photoUrls,
		"created_at":  product.CreatedAt,
		"updated_at":  product.UpdatedAt,
	})
}

func CreateProduk(c *gin.Context) {
	namaProduk := strings.TrimSpace(c.PostForm("nama_produk"))
	if namaProduk == "" {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Nama produk wajib diisi")
		return
	}

	hargaStr := c.PostForm("harga")
	harga, err := strconv.ParseFloat(hargaStr, 64)
	if err != nil || harga <= 0 {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Harga produk harus berupa angka valid")
		return
	}

	deskripsi := c.PostForm("deskripsi")
	slugInput := strings.TrimSpace(c.PostForm("slug"))
	if slugInput == "" {
		slugInput = helpers.GenerateSlug(namaProduk)
	} else {
		slugInput = helpers.GenerateSlug(slugInput)
	}

	// Parse Kategori IDs (comma-separated or multiple form values)
	kategoriIDsStr := c.PostForm("kategori_ids")
	var kategoriIDs []uint
	if kategoriIDsStr != "" {
		parts := strings.Split(kategoriIDsStr, ",")
		for _, part := range parts {
			if idNum, err := strconv.ParseUint(strings.TrimSpace(part), 10, 32); err == nil {
				kategoriIDs = append(kategoriIDs, uint(idNum))
			}
		}
	}

	var categories []models.Kategori
	categorySlug := "general"
	if len(kategoriIDs) > 0 {
		config.DB.Where("id IN ?", kategoriIDs).Find(&categories)
		if len(categories) > 0 {
			categorySlug = categories[0].Slug
		}
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	// Handle Image Files
	form, _ := c.MultipartForm()
	var files []*multipart.FileHeader
	if form != nil {
		files = form.File["images"]
		if len(files) == 0 {
			files = form.File["foto_produk"]
		}
	}

	// Create Product inside Transaction
	tx := config.DB.Begin()

	product := models.Produk{
		NamaProduk: namaProduk,
		Slug:       slugInput,
		Harga:      harga,
		Deskripsi:  deskripsi,
		Kategori:   categories,
	}

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat produk katalog: "+err.Error())
		return
	}

	// Upload & Save Images
	var savedImages []models.ProdukImage
	for _, file := range files {
		savedPath, err := saveUploadedImage(file, categorySlug, uploadDir)
		if err != nil {
			tx.Rollback()
			helpers.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		img := models.ProdukImage{
			ProdukID: product.ID,
			FilePath: savedPath,
		}
		if err := tx.Create(&img).Error; err != nil {
			tx.Rollback()
			helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal menyimpan data gambar")
			return
		}
		savedImages = append(savedImages, img)
	}

	tx.Commit()

	// Preload full response
	config.DB.Preload("Kategori").Preload("FotoProduk").First(&product, product.ID)

	var photoUrls []string
	for _, f := range product.FotoProduk {
		photoUrls = append(photoUrls, formatImageURL(c, f.FilePath))
	}

	helpers.SuccessResponse(c, http.StatusCreated, "Produk berhasil ditambahkan", gin.H{
		"id":          product.ID,
		"nama_produk": product.NamaProduk,
		"slug":        product.Slug,
		"harga":       product.Harga,
		"deskripsi":   product.Deskripsi,
		"kategori":    product.Kategori,
		"foto_produk": photoUrls,
	})
}

func UpdateProduk(c *gin.Context) {
	id := c.Param("id")
	var product models.Produk

	if err := config.DB.Preload("Kategori").Preload("FotoProduk").First(&product, id).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusNotFound, "Produk tidak ditemukan")
		return
	}

	namaProduk := strings.TrimSpace(c.PostForm("nama_produk"))
	if namaProduk != "" {
		product.NamaProduk = namaProduk
	}

	if hargaStr := c.PostForm("harga"); hargaStr != "" {
		if harga, err := strconv.ParseFloat(hargaStr, 64); err == nil && harga > 0 {
			product.Harga = harga
		}
	}

	if deskripsi := c.PostForm("deskripsi"); deskripsi != "" {
		product.Deskripsi = deskripsi
	}

	slugInput := strings.TrimSpace(c.PostForm("slug"))
	if slugInput != "" {
		product.Slug = helpers.GenerateSlug(slugInput)
	} else if namaProduk != "" {
		product.Slug = helpers.GenerateSlug(namaProduk)
	}

	// Update Categories
	kategoriIDsStr := c.PostForm("kategori_ids")
	categorySlug := "general"
	if kategoriIDsStr != "" {
		var kategoriIDs []uint
		parts := strings.Split(kategoriIDsStr, ",")
		for _, part := range parts {
			if idNum, err := strconv.ParseUint(strings.TrimSpace(part), 10, 32); err == nil {
				kategoriIDs = append(kategoriIDs, uint(idNum))
			}
		}

		var categories []models.Kategori
		config.DB.Where("id IN ?", kategoriIDs).Find(&categories)
		if len(categories) > 0 {
			categorySlug = categories[0].Slug
		}
		// Replace associations
		config.DB.Model(&product).Association("Kategori").Replace(categories)
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	// Synchronize Existing Photos (Delete removed photos)
	form, _ := c.MultipartForm()
	if form != nil {
		if rawVals, exists := form.Value["existing_photos"]; exists {
			var keptPhotos []string
			for _, v := range rawVals {
				v = strings.TrimSpace(v)
				if v == "" {
					continue
				}
				// If JSON array string: ["url1", "url2"]
				if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") {
					var list []string
					if err := json.Unmarshal([]byte(v), &list); err == nil {
						keptPhotos = append(keptPhotos, list...)
						continue
					}
				}
				// Split by comma if any
				for _, part := range strings.Split(v, ",") {
					part = strings.TrimSpace(part)
					if part != "" {
						keptPhotos = append(keptPhotos, part)
					}
				}
			}

			// Delete photos not present in keptPhotos
			for _, img := range product.FotoProduk {
				isKept := false
				imgBase := filepath.Base(img.FilePath)
				for _, kept := range keptPhotos {
					if strings.Contains(kept, img.FilePath) || (imgBase != "" && strings.Contains(kept, imgBase)) {
						isKept = true
						break
					}
				}

				if !isKept {
					if !strings.HasPrefix(img.FilePath, "http") {
						_ = os.Remove(img.FilePath)
					}
					config.DB.Delete(&img)
				}
			}
		}

		// Handle New Image Uploads if any
		files := form.File["images"]
		if len(files) == 0 {
			files = form.File["foto_produk"]
		}

		for _, file := range files {
			savedPath, err := saveUploadedImage(file, categorySlug, uploadDir)
			if err == nil {
				img := models.ProdukImage{
					ProdukID: product.ID,
					FilePath: savedPath,
				}
				config.DB.Create(&img)
			}
		}
	}

	config.DB.Save(&product)

	// Fetch updated product
	config.DB.Preload("Kategori").Preload("FotoProduk").First(&product, product.ID)

	var photoUrls []string
	for _, f := range product.FotoProduk {
		photoUrls = append(photoUrls, formatImageURL(c, f.FilePath))
	}

	helpers.SuccessResponse(c, http.StatusOK, "Produk berhasil diperbarui", gin.H{
		"id":          product.ID,
		"nama_produk": product.NamaProduk,
		"slug":        product.Slug,
		"harga":       product.Harga,
		"deskripsi":   product.Deskripsi,
		"kategori":    product.Kategori,
		"foto_produk": photoUrls,
	})
}

func DeleteProduk(c *gin.Context) {
	id := c.Param("id")
	var product models.Produk

	if err := config.DB.Preload("FotoProduk").First(&product, id).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusNotFound, "Produk tidak ditemukan")
		return
	}

	// Hapus file foto fisik dari disk
	for _, img := range product.FotoProduk {
		if !strings.HasPrefix(img.FilePath, "http") {
			_ = os.Remove(img.FilePath)
		}
	}

	// Hapus record gambar dan relasi kategori
	config.DB.Where("produk_id = ?", product.ID).Delete(&models.ProdukImage{})
	config.DB.Model(&product).Association("Kategori").Clear()

	// Hapus record produk
	if err := config.DB.Delete(&product).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus produk")
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "Produk berhasil dihapus", nil)
}

func UploadImages(c *gin.Context) {
	categorySlug := c.DefaultPostForm("kategori_slug", "general")
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	form, err := c.MultipartForm()
	if err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Gagal memproses form upload")
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		files = form.File["foto_produk"]
	}

	var uploadedPaths []string
	for _, file := range files {
		savedPath, err := saveUploadedImage(file, categorySlug, uploadDir)
		if err != nil {
			helpers.ErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		uploadedPaths = append(uploadedPaths, formatImageURL(c, savedPath))
	}

	helpers.SuccessResponse(c, http.StatusOK, "Gambar berhasil diunggah", gin.H{
		"files": uploadedPaths,
	})
}
