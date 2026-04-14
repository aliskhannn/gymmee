// Package telegram provides the delivery layer for Telegram Bot interactions.
package telegram

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/aliskhannn/gym-log/internal/domain"
)

type UserService interface {
	GetOrCreateUser(ctx context.Context, telegramID int64, username *string) (*domain.User, error)
}

// Handler processes incoming updates from the Telegram bot.
type Handler struct {
	userService UserService
	miniAppURL  string
}

// NewHandler creates a new instance of Handler.
func NewHandler(userService UserService, miniAppURL string) *Handler {
	return &Handler{
		userService: userService,
		miniAppURL:  miniAppURL,
	}
}

// Start handles the /start command. It registers the user and sends a welcome
// message with an inline button to open the Mini App.
func (h *Handler) Start(ctx context.Context, b *bot.Bot, update *models.Update) {
	// Ensure the message and sender exist
	if update.Message == nil || update.Message.From == nil {
		return
	}

	tgUser := update.Message.From
	var username *string
	if tgUser.Username != "" {
		username = &tgUser.Username
	}

	// Register or get the user
	_, err := h.userService.GetOrCreateUser(ctx, tgUser.ID, username)
	if err != nil {
		log.Printf("error getting or creating user: %v\n", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Произошла ошибка при регистрации. Попробуйте позже.",
		})
		return
	}

	// Build the inline keyboard with the Mini App button
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text: "💪 Открыть GymLog",
					WebApp: &models.WebAppInfo{
						URL: h.miniAppURL,
					},
				},
			},
		},
	}

	// Send the welcome message
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Привет! Я твой персональный трекер тренировок.\n\nНажми кнопку ниже, чтобы открыть приложение и начать тренировку.",
		ReplyMarkup: kb,
	})
}
