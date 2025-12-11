package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/HIUNCY/simple-multi-tenant-notes-api/internal/config"
	"github.com/HIUNCY/simple-multi-tenant-notes-api/internal/database"
	"github.com/HIUNCY/simple-multi-tenant-notes-api/internal/handler"
	"github.com/HIUNCY/simple-multi-tenant-notes-api/internal/middleware"
	"github.com/HIUNCY/simple-multi-tenant-notes-api/internal/repository"
	"github.com/HIUNCY/simple-multi-tenant-notes-api/internal/service"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db, err := database.NewPostgresDB(cfg.DBUrl)
	if err != nil {
		log.Fatalf("Gagal connect database: %v", err)
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Fatalf("Gagal close database: %v", err)
		}
	}(db)

	mongoClient, err := database.NewMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatalf("Gagal connect MongoDB: %v", err)
	}

	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			log.Printf("Error disconnecting Mongo: %v", err)
		}
	}()

	log.Println("Menjalankan migrasi database...")
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}
	log.Println("Tabel 'notes' berhasil dibuat/dipastikan ada!")

	// NOTE FEATURE
	noteRepo := repository.NewPostgresNoteRepository(db)
	auditRepo := repository.NewMongoAuditRepository(mongoClient, cfg.MongoDBName)
	noteService := service.NewNoteService(noteRepo, auditRepo)
	noteHandler := handler.NewNoteHandler(noteService)

	// AUTH FEATURE
	authHandler := handler.NewAuthHandler(cfg)

	enforcer, err := casbin.NewEnforcer("model.conf", "policy.csv")
	if err != nil {
		log.Fatalf("Gagal init Casbin: %v", err)
	}
	log.Println("Casbin Enforcer siap!")

	r := gin.Default()

	r.POST("/login", authHandler.Login)

	api := r.Group("/api")
	{
		api.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		api.Use(middleware.CasbinMiddleware(enforcer))
		{
			api.POST("/notes", noteHandler.Create)
			api.GET("/notes", noteHandler.GetAll)
			api.GET("/notes/:id", noteHandler.GetByID)
		}
	}

	log.Printf("Server berjalan di port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
