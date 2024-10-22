package cache

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

var contxt = context.Background()

func verificationTokenGenerator(token string) string {

	return fmt.Sprintf("light-room-user-verification-%v", token)
}

func SetUserVerificationToken(userId uuid.UUID, token string) {
	key := verificationTokenGenerator(token)
	_ = LRedis.Set(
		contxt, key, userId.String(), 24*time.Hour,
	).Err()

}

func GetUserVerificationToken(token string) (string, error) {
	key := verificationTokenGenerator(token)
	return LRedis.Get(contxt, key).Result()
}

func DeleteUserVerificationToken(token string) error {
	key := verificationTokenGenerator(token)
	return LRedis.Del(contxt, key).Err()
}

func TokenKeyGeneration(token string) string {
	return fmt.Sprintf("light-room-token-%v", token)
}

func SetToken(token string) {
	key := TokenKeyGeneration(token)

	_ = LRedis.Set(contxt, key, token, 48*time.Hour).Err()
}

func GetToken(token string) (string, error) {
	key := TokenKeyGeneration(token)
	return LRedis.Get(contxt, key).Result()
}

func PasswordResetKey(token string) string {
	return fmt.Sprintf("light-room-password-reset-%v", token)
}

func SetPasswordToken(token string, userId uuid.UUID) {
	key := PasswordResetKey(token)
	_ = LRedis.Set(contxt, key, userId.String(), 15*time.Minute).Err()

}

func GetPasswordToken(token string) (string, error) {
	key := PasswordResetKey(token)
	return LRedis.Get(contxt, key).Result()
}
