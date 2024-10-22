package schemas

// Registration Payload
type UserPayload struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"gt=1,lte=15"`
}

// User Update Payload
type UserUpdatePayload struct {
	Name       *string `json:"name"`
	Password   *string `json:"password" validate:"gt=1,lte=15"`
	IsVerified *bool   `json:"is_verified"`
}

// Change Password Payload
type ChangePasswordPayload struct {
	OldPassword string `json:"old_password" validate:"gt=1,lte=15"`
	NewPassword string `json:"new_password" validate:"gt=1,lte=15"`
}

// Password Reset Payload
type PasswordResetPayload struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"gt=1,lte=15"`
}

// Login Payload
type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"gt=1,lte=15"`
}

// Token Payload
type TokenPayload struct {
	Token string `json:"token" validate:"required,token"`
}

// AccessToken Payload
type AccessTokenPayload struct {
	AccessToken string `json:"access_token" validate:"required"`
}

// LOGOUT PAYLOAD
type LogoutPayload struct {
	AccessTokenPayload
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type EmailPayload struct {
	Email string `json:"email" validate:"required,email"`
}

type AccessPayload struct {
	LogoutPayload
	AccountVerified bool `json:"account_verified"`
}
