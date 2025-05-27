package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type strkey string

const key strkey = "userID"

func isWebSocketUpgrade(r *http.Request) bool {
	// хак, потому что хэдеры не передаются в вебсокет в playground
	return strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade") &&
		strings.ToLower(r.Header.Get("Upgrade")) == "websocket"
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isWebSocketUpgrade(r) {
			ctx := context.WithValue(r.Context(), key, 1)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		claims, err := ParseToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), key, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(key).(int)
	return userID, ok
}

func isSubscriptionOperation(ctx context.Context) bool {
	// тоже хак
	opCtx := graphql.GetOperationContext(ctx)
	if opCtx == nil {
		return false
	}
	return opCtx.Operation.Operation == "subscription"
}

func AuthMiddleware(ctx context.Context, obj any, next graphql.Resolver) (any, error) {
	if isSubscriptionOperation(ctx) {
		return next(ctx)
	}
	if _, ok := GetUserID(ctx); !ok {
		return nil, &gqlerror.Error{
			Message: "Access denied",
			Extensions: map[string]any{
				"code": "UNAUTHENTICATED",
			},
		}
	}
	return next(ctx)
}
