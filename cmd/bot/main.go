// Package main is the entry point for the gym-log telegram bot backend.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/aliskhannn/gym-log/internal/config"
	"github.com/aliskhannn/gym-log/internal/delivery/telegram"
	"github.com/aliskhannn/gym-log/internal/infra/db"
	"github.com/aliskhannn/gym-log/internal/repository/sqlite"
	"github.com/aliskhannn/gym-log/internal/service"
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

	// 5. Initialize Services
	userService := service.NewUserService(userRepo)

	// 6. Initialize Delivery (Telegram Handler)
	tgHandler := telegram.NewHandler(userService, cfg.MiniAppURL)

	// 7. Initialize and configure Telegram Bot
	opts := []bot.Option{
		bot.WithDefaultHandler(func(ctx context.Context, b *bot.Bot, update *models.Update) {
			// Optional: handle unknown messages
		}),
	}

	tgBot, err := bot.New(cfg.BotToken, opts...)
	if err != nil {
		log.Fatalf("failed to init telegram bot: %v", err)
	}

	// Register commands
	tgBot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, tgHandler.Start)

	// 8. Start the bot
	log.Println("Bot is running...")
	tgBot.Start(ctx)

	log.Println("Graceful shutdown completed")
}
