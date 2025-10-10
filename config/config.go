package config

import (
	"log"
	"modulo/structs"
	"os"

	"github.com/joho/godotenv"
)

func CargaConfig() *structs.Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error al cargar el archivo .env")
	}
	return &structs.Config{
		BaseEmail: os.Getenv("baseEmail"),
		PasswordBaseEmail: os.Getenv("passwordBaseEmail"),
	}
}
