package middleware

import (
	"context"
	"net/http"

	"github.com/deeploy-sh/deeploy/internal/server/auth"
	"github.com/deeploy-sh/deeploy/internal/server/cookie"
	"github.com/deeploy-sh/deeploy/internal/server/jwt"
	"github.com/deeploy-sh/deeploy/internal/server/service"
	"github.com/deeploy-sh/deeploy/internal/server/ui/pages"
)

type AuthMiddleWare struct {
	userService service.UserServiceInterface
}

func NewAuthMiddleware(userService service.UserServiceInterface) *AuthMiddleWare {
	return &AuthMiddleWare{userService: userService}
}

func getToken(r *http.Request) string {
	// CLI token
	authHeader := r.Header.Get("Authorization")
	const bearerPrefix = "Bearer "
	if len(authHeader) > len(bearerPrefix) && authHeader[:len(bearerPrefix)] == bearerPrefix {
		return authHeader[len(bearerPrefix):]
	}
	// Web token
	return cookie.GetTokenFromCookie(r)
}

func (m *AuthMiddleWare) Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := getToken(r)
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		t, claims, err := jwt.ValidateToken(token)
		if err != nil || !t.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := m.userService.GetUserByID(userID)
		if err != nil || user == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func RequireAuth(next http.HandlerFunc, redirectTo ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isCLI := r.URL.Query().Get("cli") == "true"

		token := getToken(r)
		if token == "" {
			path := "/"
			if len(redirectTo) > 0 {
				path = redirectTo[0]
			}

			// CLI Auth need CLI params
			if isCLI {
				if r.URL.RawQuery != "" {
					path += "?" + r.URL.RawQuery // Behalte cli=true&port=xyz
				}
			}

			http.Redirect(w, r, path, http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func RequireGuest(next http.HandlerFunc, redirectTo ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		isCLI := r.URL.Query().Get("cli") == "true"
		session := r.URL.Query().Get("session")

		token := getToken(r)
		if token != "" {
			// CLI flow - store token for polling session
			if isCLI && session != "" {
				auth.SetSessionToken(session, token)
				pages.CliAuthSuccess().Render(r.Context(), w)
				return
			}

			// Web flow
			path := "/dashboard"
			if len(redirectTo) > 0 {
				path = redirectTo[0]
			}
			http.Redirect(w, r, path, http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	}
}
