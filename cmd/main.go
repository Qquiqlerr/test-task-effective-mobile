package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	_ "testEffectiveMobile/docs"
	"testEffectiveMobile/internal/controller"
	"testEffectiveMobile/internal/repository"
	"testEffectiveMobile/internal/service"
	"testEffectiveMobile/internal/service/mocks"
	"testEffectiveMobile/internal/utils/config"
	"testEffectiveMobile/internal/utils/logger"
	"testEffectiveMobile/internal/utils/storage"
	"time"
)

// @title Test task Effective Mobile API
// @version 1.0
// @description This is a test task for Effective Mobile
// @host localhost:8080
// @BasePath /

func main() {
	//Инициализация логгера, конфига, хранилища
	cfg := config.MustLoad()
	log := logger.NewLogger(cfg.Logger.Level)
	storage.MustLoadPostgres(cfg.Database)

	//Инициализация слоев приложения
	songRepo := repository.NewRepository(log, storage.MustLoadPostgres(cfg.Database))
	apiClient := mocks.NewAPIClientMock()
	songService := service.NewService(songRepo, log, apiClient)
	songController := controller.NewController(songService, log)

	//Загрузка роутов
	router := LoadRoutes(songController)

	//Запуск сервера
	var server = http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}
	log.Info("Server started on port " + cfg.Server.Port)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Error("Server error: " + err.Error())
		}
	}()

	//Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Info("Server shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown")
	}
	log.Info("Server shutdown")
}

func LoadRoutes(controller *controller.SongController) *chi.Mux {
	router := chi.NewRouter()
	router.Get("/swagger/*", httpSwagger.Handler())
	router.Get("/songs", controller.GetSongs)
	router.Post("/songs", controller.CreateSong)
	router.Get("/songs/{id}/verses", controller.GetVersesByID)
	router.Delete("/songs/{id}", controller.DeleteSong)
	router.Put("/songs/{id}", controller.UpdateSong)
	return router
}
