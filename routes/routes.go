package routes

import (
	"candyzflorist-go/handlers"
	"candyzflorist-go/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "running",
			"message": "Backend API is fully operational!",
		})
	})

	api := r.Group("/api")

	// 1. Static Uploads File Serving
	r.Static("/uploads", "./uploads")
	api.Static("/uploads", "./uploads")

	// 2. Public Storefront Routes (No Auth Required)
	public := api.Group("/public")
	{
		public.GET("/kategori", handlers.GetPublicKategori)
		public.GET("/katalog", handlers.GetPublicKatalog)
		public.GET("/katalog/:identifier", handlers.GetPublicProdukDetail)
		public.GET("/pengaturan", handlers.GetPublicPengaturan)
	}

	// 3. Auth Routes (Public)
	auth := api.Group("/auth")
	{
		auth.POST("/login", handlers.Login)
	}

	// 4. Auth Routes (Protected)
	authProtected := api.Group("/auth")
	authProtected.Use(middleware.AuthMiddleware())
	{
		authProtected.POST("/logout", handlers.Logout)
		authProtected.GET("/me", handlers.Me)
		authProtected.PUT("/change-password", handlers.ChangePassword)
	}

	// 5. Kategori Routes (Protected Admin)
	kategori := api.Group("/kategori")
	kategori.Use(middleware.AuthMiddleware())
	{
		kategori.GET("", handlers.GetAllKategori)
		kategori.GET("/:id", handlers.GetKategoriByID)
		kategori.POST("", handlers.CreateKategori)
		kategori.PUT("/:id", handlers.UpdateKategori)
		kategori.DELETE("/:id", handlers.DeleteKategori)
	}

	// 6. Katalog / Produk Routes (Protected Admin)
	katalog := api.Group("/katalog")
	katalog.Use(middleware.AuthMiddleware())
	{
		katalog.GET("", handlers.GetAllProduk)
		katalog.GET("/:id", handlers.GetProdukByID)
		katalog.POST("", handlers.CreateProduk)
		katalog.PUT("/:id", handlers.UpdateProduk)
		katalog.DELETE("/:id", handlers.DeleteProduk)
		katalog.POST("/upload", handlers.UploadImages)
	}

	// 7. Pengaturan Sistem Routes (Protected Admin)
	pengaturan := api.Group("/pengaturan")
	pengaturan.Use(middleware.AuthMiddleware())
	{
		pengaturan.GET("", handlers.GetPengaturan)
		pengaturan.PUT("", handlers.UpdatePengaturan)
	}

	// 8. Dashboard Routes (Protected Admin)
	dashboard := api.Group("/dashboard")
	dashboard.Use(middleware.AuthMiddleware())
	{
		dashboard.GET("/stats", handlers.GetDashboardStats)
	}
}
