// Package http provides REST API handlers and middlewares.
package http

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

// ContextKey is a custom type for context keys to avoid collisions.
type ContextKey string

const UserContextKey ContextKey = "telegram_user"

// TelegramUser represents the user data embedded in initData.
type TelegramUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// AuthMiddleware validates the Telegram WebApp initData passed in the Authorization header.
func AuthMiddleware(botToken string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ожидаем заголовок вида: Authorization: tma <initData>
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "tma ") {
			http.Error(w, "Unauthorized: missing or invalid authorization header", http.StatusUnauthorized)
			return
		}

		initData := strings.TrimPrefix(authHeader, "tma ")

		user, err := validateAndExtractUser(initData, botToken)
		if err != nil {
			http.Error(w, "Unauthorized: invalid signature", http.StatusUnauthorized)
			return
		}

		// Кладем пользователя в контекст запроса
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// validateAndExtractUser verifies the HMAC signature of Telegram initData.
func validateAndExtractUser(initData, token string) (*TelegramUser, error) {
	parsed, err := url.ParseQuery(initData)
	if err != nil {
		return nil, err
	}

	hash := parsed.Get("hash")
	if hash == "" {
		return nil, fmt.Errorf("hash is missing")
	}
	parsed.Del("hash")

	var keys []string
	for k := range parsed {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var dataCheckArr []string
	for _, k := range keys {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("%s=%s", k, parsed.Get(k)))
	}
	dataCheckString := strings.Join(dataCheckArr, "\n")

	secretKey := hmacSHA256([]byte("WebAppData"), []byte(token))
	calculatedHash := hex.EncodeToString(hmacSHA256(secretKey, []byte(dataCheckString)))

	if calculatedHash != hash {
		return nil, fmt.Errorf("invalid hash")
	}

	userStr := parsed.Get("user")
	var user TelegramUser
	if err := json.Unmarshal([]byte(userStr), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func hmacSHA256(key, data []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}
