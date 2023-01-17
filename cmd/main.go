package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"shop/internal/delivery"
	"shop/internal/repository"
	"shop/internal/server"
	"shop/internal/service"
	"syscall"
	"time"
)

func main() {
	db, err := repository.Init()
	if err != nil {
		log.Print(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Can't close db err: %v\n", err)
		} else {
			log.Print("db closed")
		}
	}()
	if err := repository.CreateDatabase(db); err != nil {
		log.Print(err)
		return
	}
	repositories := repository.NewRepository(db)
	redis := repository.InitRedis()
	redisrepository := repository.NewRedisRepository(redis)
	services := service.NewServices(repositories, *redisrepository)
	handlers := delivery.NewHandler(services)
	server := new(server.Server)
	go func() {
		if err := server.Start(":8080", handlers.Handlers()); err != nil {
			log.Println(err)
			return
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	if err = server.Shutdown(ctx); err != nil {
		log.Print(err)
		return
	}
}
