package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"io/ioutil"
	"lightRoom/pkg/models"
	"lightRoom/pkg/schemas"
	"lightRoom/pkg/utils"
	"log"
	"net/http"
	"strconv"
)

func CreateTag(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	var err error
	var tagPayload schemas.TagProfile

	err = json.Unmarshal(body, &tagPayload)
	if err != nil {
		utils.JSONResponse(writer, "body not valid", http.StatusUnprocessableEntity)
		return
	}
	err = validate.Struct(&tagPayload)
	if err != nil {
		validationError := err.(validator.ValidationErrors)

		utils.JSONResponse(writer, validationError.Error(), http.StatusUnprocessableEntity)
		return
	}

	_, err = models.FetchTag(tagPayload.Title)

	if err == nil {

		utils.JSONResponse(writer, "tag already exists", http.StatusUnprocessableEntity)
		return
	}

	tag := models.Tag{
		uuid.New(),
		tagPayload.Title,
		0,
	}

	err = models.CreateTag(tag)

	if err != nil {
		utils.JSONResponse(writer, "tag creation failed", http.StatusInternalServerError)
		return
	}

	utils.DSJsonResponse(writer, []byte(`{}`), http.StatusCreated)
	return
}

func GetTags(writer http.ResponseWriter, request *http.Request) {

	// Convert limit and offset to integers (default to 10 and 0)

	title := request.URL.Query().Get("title")

	limitParam := request.URL.Query().Get("limit")
	offsetParam := request.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetParam)

	if err != nil {
		offset = 10
	}

	tags, err := models.GetTags(title, limit, offset)
	if err != nil {
		log.Println(err)
	}

	tagJson, _ := json.Marshal(tags)
	utils.DSJsonResponse(writer, tagJson, http.StatusOK)
	return
}

func GetTag(writer http.ResponseWriter, request *http.Request) {
	idStr := chi.URLParam(request, "tag_id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.JSONResponse(writer, "tag id not valid", http.StatusUnprocessableEntity)
	}

	tag, err := models.GetTag(id)

	if err != nil {
		utils.JSONResponse(writer, "tag not found", http.StatusNotFound)
		return
	}

	tagJson, _ := json.Marshal(tag)

	utils.DSJsonResponse(writer, tagJson, http.StatusOK)
	return

}

// Portfolio API FUNCS
func CreatePortfolio(writer http.ResponseWriter, request *http.Request) {

	body, _ := ioutil.ReadAll(request.Body)
	userID := request.Context().Value("user_id").(string)

	parsedUUID, _ := uuid.Parse(userID)

	var err error
	var portfolioPayload schemas.PortfolioPayload
	err = json.Unmarshal(body, &portfolioPayload)

	if err != nil {
		utils.JSONResponse(writer, "body not valid", http.StatusUnprocessableEntity)
		return
	}

	err = validate.Struct(&portfolioPayload)

	if err != nil {
		validationError := err.(validator.ValidationErrors)
		utils.JSONResponse(writer, validationError.Error(), http.StatusUnprocessableEntity)
		return
	}
	if len(portfolioPayload.Tags) == 0 {
		utils.JSONResponse(writer, "tags not provided", http.StatusUnprocessableEntity)
		return
	}

	tags, err := models.GetTagViaIds(portfolioPayload.Tags)

	if err != nil {
		log.Println(err)
		utils.JSONResponse(writer, "tags not valid", http.StatusBadRequest)

		return
	}

	if len(tags) == 0 {
		utils.JSONResponse(writer, "tags not provided", http.StatusBadRequest)
		return
	}

	portfolio := models.Portfolio{
		ID:              uuid.New(),
		Title:           portfolioPayload.Title,
		Description:     models.GetStringPointer(portfolioPayload.Description),
		Price:           portfolioPayload.Price,
		Tags:            tags,
		PaywalledImages: portfolioPayload.PaywalledImages,
		Images:          portfolioPayload.Images,
		UserID:          parsedUUID,
	}

	err = models.CreatePortfolio(portfolio)

	if err != nil {
		log.Println(err)
		utils.JSONResponse(writer, "portfolio creation failed", http.StatusInternalServerError)
		return
	}
	portfolioJson, _ := json.Marshal(&portfolio)
	utils.DSJsonResponse(writer, portfolioJson, http.StatusCreated)
	return

}
