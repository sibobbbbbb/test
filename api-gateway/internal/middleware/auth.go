package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/sibobbbbbb/backend-engineer-challenge/api-gateway/internal/client"
)

// ContextKey adalah tipe untuk key pada context
type ContextKey string

// Constants untuk context
const (
	UserIDKey   ContextKey = "user_id"
	UsernameKey ContextKey = "username"
	RoleKey     ContextKey = "role"
)

// AuthMiddleware adalah middleware untuk autentikasi
type AuthMiddleware struct {
	userClient client.UserClient
}

// NewAuthMiddleware membuat instance baru AuthMiddleware
func NewAuthMiddleware(userClient client.UserClient) *AuthMiddleware {
	return &AuthMiddleware{
		userClient: userClient,
	}
}

// Authenticate memvalidasi token JWT dan menyimpan informasi user di context
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ambil token dari header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Format: "Bearer {token}"
		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}
		token := splitToken[1]

		// Validasi token
		valid, userID, username, role, err := m.userClient.ValidateToken(r.Context(), token)
		if err != nil || !valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Simpan informasi user di context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UsernameKey, username)
		ctx = context.WithValue(ctx, RoleKey, role)

		// Panggil handler selanjutnya dengan context yang sudah dimodifikasi
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RoleRequired memastikan user memiliki role yang diperlukan
func (m *AuthMiddleware) RoleRequired(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ambil role dari context
		userRole, ok := r.Context().Value(RoleKey).(string)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Periksa apakah role sesuai
		if userRole != role && userRole != "admin" { // Admin memiliki akses ke semua
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Panggil handler selanjutnya
		next.ServeHTTP(w, r)
	})
}

// AdminRequired memastikan user adalah admin
func (m *AuthMiddleware) AdminRequired(next http.Handler) http.Handler {
	return m.RoleRequired("admin", next)
}

// GetUserID mengambil user ID dari context
func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}

// GetUsername mengambil username dari context
func GetUsername(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(UsernameKey).(string)
	return username, ok
}

// GetRole mengambil role dari context
func GetRole(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(RoleKey).(string)
	return role, ok
}