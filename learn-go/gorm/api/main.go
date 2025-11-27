package main

import (
	"context"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	ctx := context.Background()
	db.AutoMigrate(&Product{})

	// Generics API
	// Create
	err = gorm.G[Product](db).Create(ctx, &Product{Code: "D42", Price: 100})

	// Read
	product, err := gorm.G[Product](db).Where("id = ?", 1).First(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(product)
	// find product with integer primary key
	products, err := gorm.G[Product](db).Where("code = ?", "D42").Find(ctx) // find product with code D42
	if err != nil {
		panic(err)
	}
	fmt.Println(products)

}
