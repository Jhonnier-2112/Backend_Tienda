package main

import (
	"fmt"
	"log"
	"tienda-backend/internal/config"
	"tienda-backend/internal/core/domain"
	"tienda-backend/internal/adapters/repository"
)

func main() {
	cfg := config.LoadConfig()
	db, err := repository.ConnectDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	var products []domain.Product
	db.Find(&products)

	fmt.Println("Current Products in Database:")
	for _, p := range products {
		fmt.Printf("ID: %d, Name: %s, SKU: %s\n", p.ID, p.Name, p.SKU)
	}
}
