package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DBUrl       string
	MongoURI    string
	MongoDBName string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system env")
	}

	dbUrl := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=" + os.Getenv("DB_SSLMODE")

	if os.Getenv("DB_CHANNEL_BINDING") != "" {
		dbUrl += " channel_binding=" + os.Getenv("DB_CHANNEL_BINDING")
	}

	if dbUrl == "" {
		log.Fatal("Error: Koneksi DB string kosong. Cek variabel DB di .env")
	}

	mongoURI := os.Getenv("MONGO_URI")
	mongoDBName := os.Getenv("MONGO_DB_NAME")

	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	if mongoDBName == "" {
		mongoDBName = "audit_logs"
	}

	return &Config{
		Port:        os.Getenv("PORT"),
		DBUrl:       dbUrl,
		MongoURI:    mongoURI,
		MongoDBName: mongoDBName,
	}
}
