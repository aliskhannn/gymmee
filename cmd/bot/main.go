// Package main is the entry point for the gymmee telegram bot backend.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/aliskhannn/gymmee/internal/config"
	deliveryhttp "github.com/aliskhannn/gymmee/internal/delivery/http"
	"github.com/aliskhannn/gymmee/internal/delivery/telegram"
	"github.com/aliskhannn/gymmee/internal/infra/db"
	"github.com/aliskhannn/gymmee/internal/repository/sqlite"
	"github.com/aliskhannn/gymmee/internal/service"
)

func main() {
	// 1. Load configuration
	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 2. Initialize context with graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// 3. Connect to database
	database, err := db.NewSqliteDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer database.Close()

	// 4. Initialize Repositories
	userRepo := sqlite.NewUserRepository(database)
	workoutRepo := sqlite.NewWorkoutRepository(database)
	exerciseRepo := sqlite.NewExerciseRepository(database)
	habitRepo := sqlite.NewHabitRepository(database)

	// 5. Initialize Services
	userService := service.NewUserService(userRepo)
	workoutService := service.NewWorkoutService(workoutRepo, exerciseRepo)
	habitService := service.NewHabitService(habitRepo)
	exerciseService := service.NewExerciseService(exerciseRepo)

	// 6. Initialize Delivery (Telegram & HTTP Handlers)
	tgHandler := telegram.NewHandler(userService, cfg.MiniAppURL)

	httpUserHandler := deliveryhttp.NewUserHandler(userService)
	httpWorkoutHandler := deliveryhttp.NewWorkoutHandler(workoutService, userService)
	httpHabitHandler := deliveryhttp.NewHabitHandler(habitService, userService)
	httpExerciseHandler := deliveryhttp.NewExerciseHandler(exerciseService, userService)

	router := deliveryhttp.SetupRouter(
		cfg.BotToken,
		httpUserHandler,
		httpWorkoutHandler,
		httpHabitHandler,
		httpExerciseHandler,
	)

	// 7. Configure HTTP Server
	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	// 8. Configure Telegram Bot
	opts := []bot.Option{
		bot.WithDefaultHandler(func(ctx context.Context, b *bot.Bot, update *models.Update) {}),
	}

	tgBot, err := bot.New(cfg.BotToken, opts...)
	if err != nil {
		log.Fatalf("failed to init telegram bot: %v", err)
	}

	// Register commands
	tgBot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, tgHandler.Start)

	// 9. Start HTTP Server in a goroutine
	go func() {
		log.Printf("HTTP server is running on %s...", cfg.ServerAddress)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// 10. Start Telegram Bot (blocks until ctx is canceled via SIGINT/SIGTERM)
	log.Println("Telegram bot is running...")
	tgBot.Start(ctx)

	// 11. Graceful Shutdown for HTTP Server
	log.Println("Shutting down gracefully...")

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}

	log.Println("Graceful shutdown completed")
}
