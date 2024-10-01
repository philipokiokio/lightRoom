package utils

import (
	"context"
	"errors"
	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"net/http"
	"os"
	"time"
)

var TokenAuth *jwtauth.JWTAuth
var JwtSecret string = os.Getenv("JWT_SECRET")

func GenerateAccessToken(userId uuid.UUID) string {
	TokenAuth = jwtauth.New("HS256", []byte(JwtSecret), nil)

	_, tokenString, _ := TokenAuth.Encode(map[string]interface{}{"user_id": userId,
		"exp": time.Now().Add(time.Hour * 24).Unix()})

	return tokenString
}

func GenerateRefreshToken(userId uuid.UUID) string {
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

	parsedUUID := userID.(uuid.UUID)

	return GenerateAccessToken(parsedUUID), nil
}

// Authenticator is a default authentication middleware to enforce access from the
// Verifier middleware request context values. The Authenticator sends a 401 Unauthorized
// response for any unverified tokens and passes the good ones through.
func LightRoomTicator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		//authHeader := request.Header.Get("Authorization")

		// Extract the token from the Authorization header (assuming "Bearer <token>")
		//tokenStr := authHeader[len("Bearer "):]
		//_, err := cache.GetToken(tokenStr)
		//if err == nil {
		//	writer.Header().Set("WWW-Authenticate", "Bearer")
		//	writer.Header().Set("Content-Type", "application/json")
		//	writer.WriteHeader(http.StatusUnauthorized)
		//	writer.Write([]byte(`{"detail":"Invalidated token provided"}`))
		//	return
		//}
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

		userID, _ := claims["user_id"]
		contxt := context.WithValue(request.Context(), "user_id", userID)
		request.WithContext(contxt)

		// Token is authenticated, pass it through
		next.ServeHTTP(writer, request)
	})
}
