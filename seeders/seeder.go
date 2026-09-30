package seeders

import (
	"log"
	"os"
	"path/filepath"

	"candyzflorist-go/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ProductSeed struct {
	Name         string
	Slug         string
	CategorySlug string
	Price        float64
	ImagePath    string
	Description  string
}

func SeedDatabase(db *gorm.DB) {
	if db == nil {
		return
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	// 1. Seed Admin User
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		admin := models.User{
			Username: "admin",
			Email:    "admin@candyzflorist.com",
			Password: string(hashedPassword),
		}
		if err := db.Create(&admin).Error; err == nil {
			log.Println("🌱 Seeder: Berhasil membuat akun admin default (admin@candyzflorist.com / admin123)")
		}
	}

	// 2. 11 Kategori dari halaman katalog CandyzFlorist
	catalogCategories := []models.Kategori{
		{NamaKategori: "Bouquet Ready Stock", Slug: "ready-stock"},
		{NamaKategori: "Bouquet Wisuda", Slug: "wisuda"},
		{NamaKategori: "Bag Charm Bouquet", Slug: "bag-charm"},
		{NamaKategori: "Snack Bouquet", Slug: "snack"},
		{NamaKategori: "Mini Bouquet", Slug: "mini"},
		{NamaKategori: "Hand Bouquet Wedding", Slug: "wedding"},
		{NamaKategori: "Flower Box / Bloom Box", Slug: "flower-box"},
		{NamaKategori: "Vase Bouquet", Slug: "vase"},
		{NamaKategori: "Custom Bouquet", Slug: "custom"},
		{NamaKategori: "Karangan Bunga", Slug: "karangan"},
		{NamaKategori: "Money Bouquet", Slug: "money"},
	}

	categoryMap := make(map[string]models.Kategori)
	for _, cat := range catalogCategories {
		var existing models.Kategori
		if err := db.Where("slug = ?", cat.Slug).First(&existing).Error; err != nil {
			db.Create(&cat)
			categoryMap[cat.Slug] = cat
		} else {
			existing.NamaKategori = cat.NamaKategori
			db.Save(&existing)
			categoryMap[cat.Slug] = existing
		}

		folder := filepath.Join(uploadDir, cat.Slug)
		_ = os.MkdirAll(folder, os.ModePerm)
	}
	log.Println("🌱 Seeder: Berhasil sinkronisasi 11 Kategori Koleksi")

	// 3. Seed Seluruh Produk Katalog
	productSeeds := []ProductSeed{
		// Ready Stock
		{Name: "COLORFUL", Slug: "colorful", CategorySlug: "ready-stock", Price: 175000, ImagePath: "uploads/ready-stock/01-colorful.png", Description: "Buket bunga segar warna-warni ceria pilihan terbaik."},
		{Name: "TROPHY FLOWERS", Slug: "trophy-flowers-1", CategorySlug: "ready-stock", Price: 225000, ImagePath: "uploads/ready-stock/02-trophy-flowers.png", Description: "Rangkaian bunga bentuk piala elegan untuk penghargaan & selebrasi."},
		{Name: "CIEL", Slug: "ciel", CategorySlug: "ready-stock", Price: 100000, ImagePath: "uploads/ready-stock/03-ciel.png", Description: "Buket bunga nuansa biru langit yang menenangkan."},
		{Name: "MOCHI", Slug: "mochi", CategorySlug: "ready-stock", Price: 115000, ImagePath: "uploads/ready-stock/04-mochi.png", Description: "Rangkaian buket manis dan lembut seperti mochi."},
		{Name: "Mila", Slug: "mila", CategorySlug: "ready-stock", Price: 115000, ImagePath: "uploads/ready-stock/05-mila.png", Description: "Buket cantik bernuansa feminin dan anggun."},
		{Name: "THUMBELINA BLOOMBOX", Slug: "thumbelina-bloombox", CategorySlug: "ready-stock", Price: 265000, ImagePath: "uploads/ready-stock/06-thumbelina-bloombox.png", Description: "Bloombox eksklusif dengan paduan bunga premium."},
		{Name: "Baby Blue", Slug: "baby-blue", CategorySlug: "ready-stock", Price: 145000, ImagePath: "uploads/ready-stock/07-baby-blue.png", Description: "Buket bunga mawar biru pastel dengan wrapping senada."},
		{Name: "SOFT PEACH", Slug: "soft-peach", CategorySlug: "ready-stock", Price: 147000, ImagePath: "uploads/ready-stock/08-soft-peach.png", Description: "Bunga peach lembut cocok untuk hadiah ulang tahun dan anniversary."},
		{Name: "FAIRY", Slug: "fairy", CategorySlug: "ready-stock", Price: 135000, ImagePath: "uploads/ready-stock/09-fairy.png", Description: "Buket magis dengan sentuhan warna lilac dan pink pastel."},
		{Name: "HAPPY KILA", Slug: "happy-kila", CategorySlug: "ready-stock", Price: 235000, ImagePath: "uploads/ready-stock/10-happy-kila.png", Description: "Rangkaian bunga ceria menghadirkan kebahagiaan setiap momen."},
		{Name: "TROPHY FLOWERS PREMIUM", Slug: "trophy-flowers-premium", CategorySlug: "ready-stock", Price: 225000, ImagePath: "uploads/ready-stock/11-trophy-flowers.png", Description: "Trophy flowers edisi spesial dengan wrapping gold mewah."},
		{Name: "SMILE", Slug: "smile", CategorySlug: "ready-stock", Price: 137000, ImagePath: "uploads/ready-stock/12-smile.png", Description: "Buket bunga pembawa senyum kebahagiaan."},
		{Name: "FUSHIA", Slug: "fushia", CategorySlug: "ready-stock", Price: 100000, ImagePath: "uploads/ready-stock/13-fushia.png", Description: "Buket warna fushia menyala nan menawan."},
		{Name: "POLKA", Slug: "polka", CategorySlug: "ready-stock", Price: 85000, ImagePath: "uploads/ready-stock/14-polka.png", Description: "Buket mini polkadot lucu dan ekonomis."},
		{Name: "VINTAGE", Slug: "vintage", CategorySlug: "ready-stock", Price: 85000, ImagePath: "uploads/ready-stock/15-vintage.png", Description: "Buket nuansa rustic vintage elegan."},
		{Name: "NANA", Slug: "nana", CategorySlug: "ready-stock", Price: 97000, ImagePath: "uploads/ready-stock/16-nana.png", Description: "Buket bunga segar praktis dan manis."},
		{Name: "WHITE BLOOM", Slug: "white-bloom", CategorySlug: "ready-stock", Price: 150000, ImagePath: "uploads/ready-stock/17-white-bloom.png", Description: "Rangkaian bunga putih suci nan elegan."},
		{Name: "FEELING", Slug: "feeling", CategorySlug: "ready-stock", Price: 155000, ImagePath: "uploads/ready-stock/18-feeling.png", Description: "Buket bunga untuk mengungkapkan perasaan terdalam."},

		// Wisuda
		{Name: "BUKET NIMO", Slug: "buket-nimo", CategorySlug: "wisuda", Price: 150000, ImagePath: "uploads/wisuda/01-buket-nimo.png", Description: "Buket wisuda spesial boneka dan bunga artificial berkualitas."},
		{Name: "BUKET LOLY", Slug: "buket-loly", CategorySlug: "wisuda", Price: 135000, ImagePath: "uploads/wisuda/02-buket-loly.png", Description: "Buket wisuda ceria warna pastel."},
		{Name: "BUKET BUNNY", Slug: "buket-bunny", CategorySlug: "wisuda", Price: 150000, ImagePath: "uploads/wisuda/03-buket-bunny.png", Description: "Buket kelinci lucu dengan bunga segar/artificial."},
		{Name: "BUKET BELLE", Slug: "buket-belle", CategorySlug: "wisuda", Price: 150000, ImagePath: "uploads/wisuda/04-buket-belle.png", Description: "Buket anggun untuk momen kelulusan dan sidang skripsi."},
		{Name: "BUKET PEACHY", Slug: "buket-peachy", CategorySlug: "wisuda", Price: 150000, ImagePath: "uploads/wisuda/05-buket-peachy.png", Description: "Buket wisuda warna peach lembut."},

		// Bag Charm
		{Name: "BUKET POLKA BAG CHARM", Slug: "buket-polka-bag-charm", CategorySlug: "bag-charm", Price: 100000, ImagePath: "uploads/bag-charm/01-buket-polka.png", Description: "Buket gantungan tas mini aesthetic & trendi."},
		{Name: "BUKET CIEL BAG CHARM", Slug: "buket-ciel-bag-charm", CategorySlug: "bag-charm", Price: 100000, ImagePath: "uploads/bag-charm/02-buket-ciel.png", Description: "Bag charm buket nuansa biru lembut."},
		{Name: "BUKET MOCHI BAG CHARM", Slug: "buket-mochi-bag-charm", CategorySlug: "bag-charm", Price: 115000, ImagePath: "uploads/bag-charm/03-buket-mochi.png", Description: "Bag charm buket puffy gemas."},
		{Name: "BUKET MILA BAG CHARM", Slug: "buket-mila-bag-charm", CategorySlug: "bag-charm", Price: 100000, ImagePath: "uploads/bag-charm/04-buket-mila.png", Description: "Bag charm bunga mini feminin."},
		{Name: "BUKET LUMI BAG CHARM", Slug: "buket-lumi-bag-charm", CategorySlug: "bag-charm", Price: 115000, ImagePath: "uploads/bag-charm/05-buket-lumi.png", Description: "Bag charm buket aesthetic viral."},

		// Snack Bouquet
		{Name: "BUKET LILAC SNACK", Slug: "buket-lilac-snack", CategorySlug: "snack", Price: 150000, ImagePath: "uploads/snack/01-buket-lilac.png", Description: "Buket snack premium nuansa lilac cantik."},
		{Name: "BUKET BLUE SNACK", Slug: "buket-blue-snack", CategorySlug: "snack", Price: 75000, ImagePath: "uploads/snack/02-buket-blue.png", Description: "Buket snack serba biru hemat & lezat."},
		{Name: "BUKET MIX SNACK", Slug: "buket-mix-snack", CategorySlug: "snack", Price: 85000, ImagePath: "uploads/snack/03-buket-mix-snack.png", Description: "Buket kombinasi aneka camilan favorit."},
		{Name: "BUKET SISI SNACK", Slug: "buket-sisi-snack", CategorySlug: "snack", Price: 50000, ImagePath: "uploads/snack/04-buket-sisi.png", Description: "Buket snack mini terjangkau untuk kado sahabat."},
		{Name: "BUKET SNOW SNACK", Slug: "buket-snow-snack", CategorySlug: "snack", Price: 75000, ImagePath: "uploads/snack/05-buket-snow.png", Description: "Buket snack putih bersih elegan."},

		// Mini Bouquet
		{Name: "BUKET SOFT VINTAGE MINI", Slug: "buket-soft-vintage-mini", CategorySlug: "mini", Price: 85000, ImagePath: "uploads/mini/01-buket-soft-vintage.png", Description: "Mini buket bunga gaya vintage rustic."},
		{Name: "BUKET TULIP MINI", Slug: "buket-tulip-mini", CategorySlug: "mini", Price: 50000, ImagePath: "uploads/mini/02-buket-tulip.png", Description: "Bunga tulip artificial mini nan cantik."},
		{Name: "BUKET BUTTERFLY MINI", Slug: "buket-butterfly-mini", CategorySlug: "mini", Price: 50000, ImagePath: "uploads/mini/03-buket-butterfly.png", Description: "Buket mini dengan aksen kupu-kupu berkilau."},
		{Name: "BUKET HAIRCLIP MINI", Slug: "buket-hairclip-mini", CategorySlug: "mini", Price: 25000, ImagePath: "uploads/mini/04-buket-hairclip.png", Description: "Buket unik berisi jepitan rambut modis."},
		{Name: "BUKET ROSE SINGLE MINI", Slug: "buket-rose-single-mini", CategorySlug: "mini", Price: 20000, ImagePath: "uploads/mini/05-buket-rose-single.png", Description: "Single rose mini untuk hadiah simpel berkesan."},
		{Name: "BUKET TEDY BEAR MINI", Slug: "buket-tedy-bear-mini", CategorySlug: "mini", Price: 100000, ImagePath: "uploads/mini/06-buket-tedy-bear.png", Description: "Boneka beruang kecil dengan buket bunga mini."},
		{Name: "BUKET HOLO SINGLE MINI", Slug: "buket-holo-single-mini", CategorySlug: "mini", Price: 20000, ImagePath: "uploads/mini/07-buket-holo-single.png", Description: "Bunga single dengan wrapping hologram modern."},
		{Name: "BUKET MIX SNACK MINI", Slug: "buket-mix-snack-mini", CategorySlug: "mini", Price: 45000, ImagePath: "uploads/mini/08-buket-mix-snack.png", Description: "Mini snack buket pas di kantong."},
		{Name: "BUKET TULIP SINGLE MINI", Slug: "buket-tulip-single-mini", CategorySlug: "mini", Price: 20000, ImagePath: "uploads/mini/09-buket-tulip-single.png", Description: "Single tulip dengan wrapping dusty pink."},
		{Name: "BUKET FRESH SINGLE MINI", Slug: "buket-fresh-single-mini", CategorySlug: "mini", Price: 20000, ImagePath: "uploads/mini/10-buket-fresh-single.png", Description: "Single fresh flower segar dipetik hari ini."},

		// Hand Bouquet Wedding & Categories Lainnya
		{Name: "Hand Bouquet Wedding 01", Slug: "hand-bouquet-wedding-01", CategorySlug: "wedding", Price: 450000, ImagePath: "uploads/ready-stock/17-white-bloom.png", Description: "Hand bouquet pernikahan mewah dengan pita satin."},
		{Name: "Flower Box 01", Slug: "flower-box-01", CategorySlug: "flower-box", Price: 250000, ImagePath: "uploads/ready-stock/06-thumbelina-bloombox.png", Description: "Kotak bunga eksklusif untuk kado spesial."},
		{Name: "Vase Bouquet 01", Slug: "vase-bouquet-01", CategorySlug: "vase", Price: 300000, ImagePath: "uploads/ready-stock/02-trophy-flowers.png", Description: "Rangkaian bunga vas meja estetik untuk rumah & kantor."},
		{Name: "Round Bouquet Custom", Slug: "round-bouquet-custom", CategorySlug: "custom", Price: 200000, ImagePath: "uploads/ready-stock/01-colorful.png", Description: "Buket bulat kustom request bunga & warna wrapping."},
		{Name: "Unique Bouquet Custom", Slug: "unique-bouquet-custom", CategorySlug: "custom", Price: 275000, ImagePath: "uploads/ready-stock/10-happy-kila.png", Description: "Desain buket unik sesuai kreativitas & tema acara."},
		{Name: "Rotan Bouquet Custom", Slug: "rotan-bouquet-custom", CategorySlug: "custom", Price: 350000, ImagePath: "uploads/ready-stock/08-soft-peach.png", Description: "Rangkaian bunga keranjang rotan natural elegan."},
		{Name: "Karangan Bunga 01", Slug: "karangan-bunga-01", CategorySlug: "karangan", Price: 500000, ImagePath: "uploads/ready-stock/11-trophy-flowers.png", Description: "Karangan bunga ucapan selamat, sukses, atau duka cita."},
		{Name: "Money Bouquet 01", Slug: "money-bouquet-01", CategorySlug: "money", Price: 600000, ImagePath: "uploads/ready-stock/09-fairy.png", Description: "Buket uang asli / tarik uang untuk ulang tahun & wisuda."},
	}

	for _, pSeed := range productSeeds {
		var existingProduct models.Produk
		cat, catExists := categoryMap[pSeed.CategorySlug]

		if err := db.Where("slug = ?", pSeed.Slug).First(&existingProduct).Error; err != nil {
			// Buat produk baru
			newProduct := models.Produk{
				NamaProduk: pSeed.Name,
				Slug:       pSeed.Slug,
				Harga:      pSeed.Price,
				Deskripsi:  pSeed.Description,
			}
			if catExists {
				newProduct.Kategori = []models.Kategori{cat}
			}
			if err := db.Create(&newProduct).Error; err == nil {
				// Buat image record
				img := models.ProdukImage{
					ProdukID: newProduct.ID,
					FilePath: pSeed.ImagePath,
				}
				db.Create(&img)
			}
		} else {
			// Update data produk yang ada
			existingProduct.NamaProduk = pSeed.Name
			existingProduct.Harga = pSeed.Price
			existingProduct.Deskripsi = pSeed.Description
			if catExists {
				db.Model(&existingProduct).Association("Kategori").Replace([]models.Kategori{cat})
			}
			db.Save(&existingProduct)

			// Cek apakah image record sudah ada
			var existingImg models.ProdukImage
			if err := db.Where("produk_id = ?", existingProduct.ID).First(&existingImg).Error; err != nil {
				img := models.ProdukImage{
					ProdukID: existingProduct.ID,
					FilePath: pSeed.ImagePath,
				}
				db.Create(&img)
			} else {
				existingImg.FilePath = pSeed.ImagePath
				db.Save(&existingImg)
			}
		}
	}
	log.Printf("🌱 Seeder: Berhasil sinkronisasi %d Produk Katalog ke Database MySQL\n", len(productSeeds))

	// 4. Seed Pengaturan Awal
	var pengaturanCount int64
	db.Model(&models.Pengaturan{}).Count(&pengaturanCount)
	if pengaturanCount == 0 {
		initialSettings := models.Pengaturan{
			Whatsapp:   "089688035866",
			Instagram:  "https://instagram.com/crandyzflorist",
			Tiktok:     "https://tiktok.com/@crandyzflorist",
			Email:      "crandyzflorist@gmail.com",
			Alamat:     "Blok F No. 528, Perumahan Bumi Telukjambe, Kec. Telukjambe Timur, Karawang, Jawa Barat 41361",
			TemplateWa: "Halo CandyzFlorist, saya tertarik untuk memesan produk *{nama_produk}* dengan harga *Rp {harga}*. Apakah masih bisa dipesan?",
		}
		if err := db.Create(&initialSettings).Error; err == nil {
			log.Println("🌱 Seeder: Berhasil membuat data pengaturan toko awal")
		}
	}
}
