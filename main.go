package main

import (
	"github/abbgo/yenil_yol/backend/config"
	"github/abbgo/yenil_yol/backend/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	// Database instance
	db, err := config.ConnDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := routes.Routes()

	// static file
	os.Mkdir("./uploads", os.ModePerm)
	r.Static("/uploads", "./uploads")

	// run routes
	if err := r.Run(":7723"); err != nil {
		log.Fatal(err)
	}

}
