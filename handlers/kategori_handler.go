package handlers

import (
	"net/http"
	"os"
	"path/filepath"

	"candyzflorist-go/config"
	"candyzflorist-go/helpers"
	"candyzflorist-go/models"
	"github.com/gin-gonic/gin"
)

type KategoriInput struct {
	NamaKategori string `json:"nama_kategori" binding:"required"`
	Slug         string `json:"slug"`
}

func GetAllKategori(c *gin.Context) {
	var categories []models.Kategori
	if err := config.DB.Order("created_at desc").Find(&categories).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data kategori")
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "Daftar kategori", categories)
}

func GetKategoriByID(c *gin.Context) {
	id := c.Param("id")
	var kategori models.Kategori

	if err := config.DB.First(&kategori, id).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusNotFound, "Kategori tidak ditemukan")
		return
	}

	helpers.SuccessResponse(c, http.StatusOK, "Detail kategori", kategori)
}

func CreateKategori(c *gin.Context) {
	var input KategoriInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Nama kategori wajib diisi")
		return
	}

	slug := input.Slug
	if slug == "" {
		slug = helpers.GenerateSlug(input.NamaKategori)
	} else {
		slug = helpers.GenerateSlug(slug)
	}

	// Cek apakah slug sudah ada
	var existing models.Kategori
	if err := config.DB.Where("slug = ?", slug).First(&existing).Error; err == nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Kategori dengan nama atau slug serupa sudah ada")
		return
	}

	kategori := models.Kategori{
		NamaKategori: input.NamaKategori,
		Slug:         slug,
	}

	if err := config.DB.Create(&kategori).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal menyimpan kategori")
		return
	}

	// Buat folder uploads/<slug>
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	categoryFolder := filepath.Join(uploadDir, slug)
	_ = os.MkdirAll(categoryFolder, os.ModePerm)

	helpers.SuccessResponse(c, http.StatusCreated, "Kategori berhasil ditambahkan", kategori)
}

func UpdateKategori(c *gin.Context) {
	id := c.Param("id")
	var kategori models.Kategori

	if err := config.DB.First(&kategori, id).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusNotFound, "Kategori tidak ditemukan")
		return
	}

	var input KategoriInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Nama kategori wajib diisi")
		return
	}

	oldSlug := kategori.Slug
	newSlug := input.Slug
	if newSlug == "" {
		newSlug = helpers.GenerateSlug(input.NamaKategori)
	} else {
		newSlug = helpers.GenerateSlug(newSlug)
	}

	// Cek duplikasi slug jika berubah
	if newSlug != oldSlug {
		var existing models.Kategori
		if err := config.DB.Where("slug = ? AND id != ?", newSlug, id).First(&existing).Error; err == nil {
			helpers.ErrorResponse(c, http.StatusBadRequest, "Kategori dengan nama atau slug serupa sudah ada")
			return
		}
	}

	kategori.NamaKategori = input.NamaKategori
	kategori.Slug = newSlug

	if err := config.DB.Save(&kategori).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui kategori")
		return
	}

	// Rename folder jika slug berubah
	if oldSlug != newSlug {
		uploadDir := os.Getenv("UPLOAD_DIR")
		if uploadDir == "" {
			uploadDir = "./uploads"
		}
		oldFolder := filepath.Join(uploadDir, oldSlug)
		newFolder := filepath.Join(uploadDir, newSlug)
		if _, err := os.Stat(oldFolder); err == nil {
			_ = os.Rename(oldFolder, newFolder)
		} else {
			_ = os.MkdirAll(newFolder, os.ModePerm)
		}
	}

	helpers.SuccessResponse(c, http.StatusOK, "Kategori berhasil diperbarui", kategori)
}

func DeleteKategori(c *gin.Context) {
	id := c.Param("id")
	var kategori models.Kategori

	if err := config.DB.First(&kategori, id).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusNotFound, "Kategori tidak ditemukan")
		return
	}

	// Cek apakah kategori terpakai di produk
	var count int64
	config.DB.Table("produk_kategoris").Where("kategori_id = ?", id).Count(&count)
	if count > 0 {
		helpers.ErrorResponse(c, http.StatusBadRequest, "Kategori tidak dapat dihapus karena masih digunakan oleh produk katalog")
		return
	}

	if err := config.DB.Delete(&kategori).Error; err != nil {
		helpers.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus kategori")
		return
	}

	// Hapus folder upload jika kosong
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	categoryFolder := filepath.Join(uploadDir, kategori.Slug)
	_ = os.Remove(categoryFolder) // hanya terhapus jika kosong

	helpers.SuccessResponse(c, http.StatusOK, "Kategori berhasil dihapus", nil)
}
