package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DataBase_URL string
	ReSend_KEY   string
	Base_URL     string
	JWTSecret    string
}

func Load() Config {
	var cf Config

	err := godotenv.Load()
	if err != nil {
		log.Println("error loading env")
	}

	cf.Port = os.Getenv("PORT")
	if cf.Port == "" {
		cf.Port = "8080"
		log.Println("error no port")
	}

	cf.DataBase_URL = os.Getenv("DATABASE_URL")
	if cf.DataBase_URL == "" {
		log.Fatal("error no database_url")
	}

	cf.ReSend_KEY = os.Getenv("RESEND_API_KEY")
	if cf.ReSend_KEY == "" {
		log.Fatal("error no ReSend_KEY")
	}

	cf.Base_URL = os.Getenv("BASE_URL")
	if cf.Base_URL == "" {
		log.Fatal("error no Base_URL")
	}

	cf.JWTSecret = os.Getenv("JWT_SECRET")
	if cf.JWTSecret == "" {
		log.Fatal("error no JWTSecret")
	}

	return cf
}
