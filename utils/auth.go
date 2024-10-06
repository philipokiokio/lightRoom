package utils

import (
	"context"
	"errors"
	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"lightRoom/cache"
	"net/http"
	"strings"
	"time"
)

var TokenAuth *jwtauth.JWTAuth

func AuthInit() {
	var JwtSecret string = Settings.JwtSecret

	TokenAuth = jwtauth.New("HS256", []byte(JwtSecret), nil)
}

func GenerateAccessToken(userId uuid.UUID) string {
	var JwtSecret string = Settings.JwtSecret

	TokenAuth = jwtauth.New("HS256", []byte(JwtSecret), nil)
	_, tokenString, _ := TokenAuth.Encode(map[string]interface{}{"user_id": userId,
		"exp": time.Now().Add(time.Hour * 24).Unix()})
	return tokenString
}

func GenerateRefreshToken(userId uuid.UUID) string {
	var JwtSecret string = Settings.JwtSecret

	TokenAuth = jwtauth.New("HS256", []byte(JwtSecret), nil)

	_, tokenString, _ := TokenAuth.Encode(map[string]interface{}{"user_id": userId,
		"exp": time.Now().Add(time.Hour * 48).Unix()})

	return tokenString
}

func VerifyRefreshToken(refreshToken string) (string, error) {
	//Parsing the Token
	token, err := jwt.ParseString(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")

	}
	//	validate the token
	if err = jwt.Validate(token); err != nil {
		return "", errors.New("invalid or expired refresh token")
	}
	userID, _ := token.Get("user_id")
	userIDStr, ok := userID.(string)

	if !ok {
		return "", errors.New("user_id is not a valid string")
	}
	parsedUUID, _ := uuid.Parse(userIDStr)
	return GenerateAccessToken(parsedUUID), nil
}

// Authenticator is a default authentication middleware to enforce access from the
// Verifier middleware request context values. The Authenticator sends a 401 Unauthorized
// response for any unverified tokens and passes the good ones through.
func LightRoomTicator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		token, claims, err := jwtauth.FromContext(request.Context())

		if err != nil {
			writer.Header().Set("WWW-Authenticate", "Bearer")
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte(`{"detail":"Unauthorized"}`))
			return

		}

		if token == nil || jwt.Validate(token) != nil {
			writer.Header().Set("WWW-Authenticate", "Bearer")
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusUnauthorized)
			writer.Write([]byte(`{"detail":"Unauthorized, token invalid"}`))
			return

		}

		authHeader := request.Header.Get("Authorization")

		if strings.HasPrefix(authHeader, "Bearer ") {
			authRawToken := strings.TrimPrefix(authHeader, "Bearer ")
			cachedToken, _ := cache.GetToken(authRawToken)

			if cachedToken != "" {
				writer.Header().Set("WWW-Authenticate", "Bearer")
				writer.Header().Set("Content-Type", "application/json")
				writer.WriteHeader(http.StatusUnauthorized)
				writer.Write([]byte(`{"detail":"token blacklisted"}`))
				return

			}

		}

		userID, _ := claims["user_id"]
		contxt := context.WithValue(request.Context(), "user_id", userID)
		request = request.WithContext(contxt)
		// Token is authenticated, pass it through
		next.ServeHTTP(writer, request)
	})
}
