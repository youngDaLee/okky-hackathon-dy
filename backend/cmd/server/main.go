package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"okky-hackathon/fridge-master-backend/internal/fridge"
	"okky-hackathon/fridge-master-backend/internal/server"
	"okky-hackathon/fridge-master-backend/pkg/config"
	"okky-hackathon/fridge-master-backend/pkg/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := database.Connect(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatalf("mongodb: %v", err)
	}

	// Fridge domain wiring
	fridgeCol := database.GetCollection(mongoClient, cfg.MongoDB, "ingredients")
	fridgeRepo := fridge.NewIngredientRepository(fridgeCol)
	fridgeSvc := fridge.NewFridgeService(fridgeRepo)
	fridgeHandler := fridge.NewFridgeHandler(fridgeSvc)

	// Ensure indexes at startup (best-effort)
	if err := fridgeRepo.EnsureIndexes(context.Background()); err != nil {
		log.Printf("warn: index creation: %v", err)
	}

	r := server.NewRouter(server.RouterDeps{
		FridgeHandler: fridgeHandler,
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("server listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown: %v", err)
	}
}
