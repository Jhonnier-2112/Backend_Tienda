package main

import (
	"fmt"
	"log"
	"tienda-backend/internal/adapters/repository"
	"tienda-backend/internal/config"
	"tienda-backend/internal/core/domain"

	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()
	db, err := repository.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// AutoMigrate new fields/tables
	db.AutoMigrate(&domain.Category{}, &domain.Product{}, &domain.ProductImage{})

	seed(db)
}

func seed(db *gorm.DB) {
	// 1. Clear existing data (optional, but good for clean validation)
	db.Exec("TRUNCATE categories, products, product_images RESTART IDENTITY CASCADE")

	// 2. Categories
	categories := []domain.Category{
		{Name: "Tecnología", Description: "Laptops, Smartphones y gadgets de última generación"},
		{Name: "Hogar Smart", Description: "Automatización y electrodomésticos inteligentes"},
		{Name: "Audio & Vídeo", Description: "Sonido premium y pantallas de alta definición"},
		{Name: "Accesorios", Description: "Complementos para tu ecosistema digital"},
	}

	for i := range categories {
		db.Create(&categories[i])
	}

	// 3. Products
	products := []domain.Product{
		{
			Name:        "MacBook Pro M3 Max",
			Description: "La laptop más potente de Apple para profesionales creativos.",
			SKU:         "APPLE-MBP-M3MX",
			Price:       15999000,
			CostPrice:   12000000,
			Stock:       10,
			MinStock:    2,
			CategoryID:  categories[0].ID,
			Specs:       `{"Procesador": "M3 Max 14-core", "Memoria": "36GB UC", "Almacenamiento": "1TB SSD", "Pantalla": "14.2\" Liquid Retina XDR"}`,
		},
		{
			Name:        "iPhone 15 Pro Titanium",
			Description: "Forjado en titanio, con el chip A17 Pro y botón de acción.",
			SKU:         "APPLE-I15P-256",
			Price:       5499000,
			CostPrice:   4200000,
			Stock:       25,
			MinStock:    5,
			CategoryID:  categories[0].ID,
			Specs:       `{"Cámara": "48MP Main", "Pantalla": "6.1\" ProMotion", "Chip": "A17 Pro", "Material": "Titanio aeroespacial"}`,
		},
		{
			Name:        "Sony WH-1000XM5",
			Description: "Líder en cancelación de ruido con un sonido excepcional.",
			SKU:         "SONY-XM5-BLK",
			Price:       1450000,
			CostPrice:   950000,
			Stock:       15,
			MinStock:    3,
			CategoryID:  categories[2].ID,
			Specs:       `{"Autonomía": "30 horas", "Carga": "3 min = 3 horas", "Micrófonos": "8 micrófonos", "Drivers": "30mm"}`,
		},
		{
			Name:        "Aspiradora Robot iRobot Roomba j7+",
			Description: "Limpia de forma inteligente y evita obstáculos de mascotas.",
			SKU:         "IROBOT-J7P-01",
			Price:       3890000,
			CostPrice:   2900000,
			Stock:       8,
			MinStock:    2,
			CategoryID:  categories[1].ID,
			Specs:       `{"Navegación": "PrecisionVision", "Vaciado": "Automático Clean Base", "Succión": "10x potencia", "Filtro": "Alta eficiencia"}`,
		},
	}

	for i := range products {
		db.Create(&products[i])

		// Add dummy gallery images
		images := []domain.ProductImage{
			{ProductID: products[i].ID, URL: "https://images.unsplash.com/photo-1517336714731-489689fd1ca8", IsPrimary: true},
			{ProductID: products[i].ID, URL: "https://images.unsplash.com/photo-1611186871348-b1ec696e5237", IsPrimary: false},
		}
		for j := range images {
			db.Create(&images[j])
		}
	}

	fmt.Println("Seed V2 completed successfully! 🚀")
}
