package api

import (
	"bytes"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"html/template"
	"io/ioutil"
	"lightRoom/cache"
	"lightRoom/models"
	"lightRoom/schemas"
	"lightRoom/utils"
	"log"
	"net/http"
)

var validate *validator.Validate

// InitializeValidator initializes the validator instance.
func InitializeValidator() {
	validate = validator.New()
}

// Auth godoc
// @Tags Auth
// @Summary Create a New User
// @Accept json
// @Produce json
// @Param user body schemas.UserPayload true "Create User Payload"
// @Router /api/v1/auth/sign-up [post]
// @Success  200  {object}  models.User
// @Failure      400  {object} schemas.ErrorPayload
func CreateUser(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	var err error
	var userPayload schemas.UserPayload

	err = json.Unmarshal(body, &userPayload)
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnprocessableEntity)
		writer.Write([]byte(`{"detail": "user body not valid"}`))
		return
	}

	err = validate.Struct(userPayload)
	if err != nil {
		validationError := err.(validator.ValidationErrors)
		jsonResponse, _ := json.Marshal(map[string]string{"detail": validationError.Error()})
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(jsonResponse)
		return

	}
	//Check to see if email has a record already
	_, err = models.FetchViaMail(userPayload.Email)

	if err == nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(`{"detail": "user already exist"}`))
		return
	}

	userPayload.Password, _ = utils.HashPassword(userPayload.Password)

	user := models.User{
		ID:         uuid.New(),
		Name:       userPayload.Name,
		Email:      userPayload.Email,
		Password:   userPayload.Password,
		IsVerified: false,
	}

	err = models.CreateUser(user)

	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(
			[]byte(`{"detail": "user creation error"}`))
		return
	}

	//We need to send email with with a token for verification
	//	pluging in redis to store token here
	token := utils.TokenGenerator()
	cache.SetUserVerificationToken(user.ID, token)

	//Send Email Template
	// Parse the email template
	verificationTemplate, _ := template.ParseFiles("templates/verification_email.html")

	// Create a data structure to pass to the template
	data := struct {
		Name  string
		Token string
	}{
		Name:  user.Name,
		Token: token,
	}

	// Execute the template and store the result in a buffer
	var payload bytes.Buffer

	err = verificationTemplate.Execute(&payload, data)
	if err == nil {

		utils.SendMail([]string{user.Email}, payload)

	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	userJson, _ := json.Marshal(user)
	writer.Write(userJson)
}

// Auth godoc
// @Tags Auth
// @Summary Login
// @Accept json
// @Produce json
// @Param user body schemas.LoginPayload true "Login Payload"
// @Router /api/v1/auth/login [post]
// @Success  200  {object}  schemas.AccessPayload
// @Failure      400  {object} schemas.ErrorPayload
func Login(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	var err error
	var loginPayload schemas.LoginPayload

	err = json.Unmarshal(body, &loginPayload)
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnprocessableEntity)
		writer.Write([]byte(`{"detail": "login body not valid"}`))
		return
	}

	err = validate.Struct(loginPayload)
	if err != nil {
		validationError := err.(validator.ValidationErrors)
		jsonResponse, _ := json.Marshal(map[string]string{"detail": validationError.Error()})
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(jsonResponse)
		return

	}

	user, err := models.FetchViaMail(loginPayload.Email)

	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(`{"detail": "user email/password is incorrect"}`))
		return
	}
	if utils.ComparePasswords(user.Password, loginPayload.Password) == false {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(`{"detail": "user email/password is incorrect"}`))
		return
	}
	if user.IsVerified == false {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(`{"detail": "user account is not verified"}`))
		return

	}
	accessToken := utils.GenerateAccessToken(user.ID)
	refreshToken := utils.GenerateRefreshToken(user.ID)
	jsonResponse, _ := json.Marshal(map[string]string{"access_token": accessToken, "refresh_token": refreshToken, "account_verified": "verified"})
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonResponse)

}

// Auth godoc
// @Tags Auth
// @Summary Me
// @Produce json
// @Security BearerAuth
// @Router /api/v1/auth/me [get]
// @Success  200  {object}  models.User
// @Failure      400  {object} schemas.ErrorPayload
func Me(writer http.ResponseWriter, request *http.Request) {

	userID := request.Context().Value("user_id").(string)

	parsedUUID, err := uuid.Parse(userID)

	user, err := models.GetUser(parsedUUID)
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(`{"detail": "user not found"}`))
		return
	}
	userJson, _ := json.Marshal(user)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(userJson)
}

// Auth godoc
// @Tags Auth
// @Summary VerifyAccount
// @Accept json
// @Produce json
// @Param user body schemas.TokenPayload true "Token Payload"
// @Router /api/v1/auth/account-verification [post]
// @Success  200  {object} schemas.MessagePayload
// @Failure      400  {object} schemas.ErrorPayload
func Verify(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)

	// Parse the body into a struct
	var tokenPayload schemas.TokenPayload
	if err := json.Unmarshal(body, &tokenPayload); err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnprocessableEntity)
		writer.Write([]byte(`{"detail": "token body not valid"}`))
		return
	}

	err := validate.Struct(tokenPayload)
	if err != nil {
		validationError := err.(validator.ValidationErrors)
		jsonResponse, _ := json.Marshal(map[string]string{"detail": validationError.Error()})
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(jsonResponse)
		return

	}

	value, err := cache.GetUserVerificationToken(tokenPayload.Token)

	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(`{"detail": "token does not exist"}`))
		return
	}
	parsedUUID, _ := uuid.Parse(value)
	var updateUser models.User
	updateUser.IsVerified = true
	models.UpdateUser(parsedUUID, updateUser)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(`{"message": "user account verified", "status": "ok"}`))
}

// Auth godoc
// @Tags Auth
// @Summary Refresh
// @Param Refresh header string true "token"
// @Produce json
// @Router /api/v1/auth/refresh [post]
// @Success  200  {object} schemas.AccessTokenPayload
// @Failure  400  {object} schemas.ErrorPayload
func Refresh(writer http.ResponseWriter, request *http.Request) {
	refreshToken := request.Header.Get("Refresh")

	if refreshToken == "" {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(`{"detail": "refresh token not provided"}`))
		return
	}
	_, err := cache.GetToken(refreshToken)
	if err == nil {
		writer.Header().Set("WWW-Authenticate", "Bearer")
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte(`{"detail":"Invalidated token provided"}`))
		return
	}

	accessToken, err := utils.VerifyRefreshToken(refreshToken)

	if err != nil {
		jsonResponse, _ := json.Marshal(map[string]string{"detail": err.Error()})
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(jsonResponse))
		return
	}

	jsonResponse, _ := json.Marshal(map[string]string{"access_token": accessToken})
	log.Println(jsonResponse)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(jsonResponse)
}

// Logout godoc
// @Tags Auth
// @Summary Logout
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body schemas.LogoutPayload true "Logout Payload"
// @Router /api/v1/auth/logout [post]
// @Failure      400  {object} schemas.ErrorPayload
func LogOut(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)

	var logoutPayload schemas.LogoutPayload

	err := json.Unmarshal(body, &logoutPayload)

	if err != nil {

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnprocessableEntity)
		writer.Write([]byte(`{"detail": "logout body not valid"}`))
		return
	}

	err = validate.Struct(logoutPayload)
	if err != nil {
		validationError := err.(validator.ValidationErrors)
		jsonResponse, _ := json.Marshal(map[string]string{"detail": validationError.Error()})
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(jsonResponse)
		return

	}
	cache.SetToken(logoutPayload.AccessToken)
	cache.SetToken(logoutPayload.RefreshToken)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(`{}`))
}

// Auth godoc
// @Tags Auth
// @Summary PasswordReset
// @Accept json
// @Produce json
// @Param user body schemas.PasswordResetPayload true "PasswordReset Payload"
// @Router /api/v1/auth/reset-password [post]
// @Success  200  {object} schemas.MessagePayload
// @Failure      400  {object} schemas.ErrorPayload
func PasswordReset(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	var passwordResetPayload schemas.PasswordResetPayload
	err := json.Unmarshal(body, &passwordResetPayload)
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnprocessableEntity)
		writer.Write([]byte(`{"detail": "body not valid"}`))
		return
	}
	err = validate.Struct(passwordResetPayload)

	if err != nil {
		validationError := err.(validator.ValidationErrors)
		jsonResponse, _ := json.Marshal(map[string]string{"detail": validationError.Error()})
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnprocessableEntity)
		writer.Write(jsonResponse)
		return
	}

	userId, err := cache.GetPasswordToken(passwordResetPayload.Token)

	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(`{"detail": "reset token has expired/not found"}`))
		return
	}

	var userUpdate models.User
	userUpdate.Password, _ = utils.HashPassword(passwordResetPayload.Password)
	parsedUUID, _ := uuid.Parse(userId)

	err = models.UpdateUser(parsedUUID, userUpdate)

	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(`{"detail": "user not found"}`))
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(`{"message": "password reset success"}`))

}

// Auth godoc
// @Tags Auth
// @Summary ForgotPassword
// @Accept json
// @Produce json
// @Param user body schemas.EmailPayload true "ForgetPassword Payload"
// @Router /api/v1/auth/forgot-password [post]
// @Success  200  {object} schemas.MessagePayload
// @Failure      400  {object} schemas.ErrorPayload
func ForgotPassword(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	var emailPayload schemas.EmailPayload

	err := json.Unmarshal(body, &emailPayload)
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(`{"detail": " email not provided"}`))
		return
	}

	err = validate.Struct(emailPayload)
	if err != nil {
		validationError := err.(validator.ValidationErrors)
		jsonResponse, _ := json.Marshal(map[string]string{"detail": validationError.Error()})
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusUnprocessableEntity)
		writer.Write(jsonResponse)
		return
	}

	user, err := models.FetchViaMail(emailPayload.Email)

	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(`{"detail": "user not found"}`))
		return
	}

	token := utils.TokenGenerator()
	cache.SetPasswordToken(token, user.ID)
	verificationTemplate, _ := template.ParseFiles("templates/verification_email.html")

	// Create a data structure to pass to the template
	data := struct {
		Name  string
		Token string
	}{
		Name:  user.Name,
		Token: token,
	}

	// Execute the template and store the result in a buffer
	var payload bytes.Buffer

	err = verificationTemplate.Execute(&payload, data)
	if err == nil {
		utils.SendMail([]string{user.Email}, payload)

	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(`{"message": "password reset mail sent"}`))
}
