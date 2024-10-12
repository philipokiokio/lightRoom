package api

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"lightRoom/schemas"
	"lightRoom/utils"
	"net/http"
)

const maxUploadFile = 10 << 20

// Misc godoc
// @Tags Misc
// @Summary UploadFile
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param fileType query string true "Type of the file" Enums(PROFILE, PORTFOLIO)
// @Param files formData file true "Files to upload" multiple=true
// @Router /api/v1/misc/upload-file [post]
// @Success  200  {object} []string
// @Failure      400  {object} schemas.ErrorPayload
func UploadFile(writer http.ResponseWriter, request *http.Request) {
	var fileURL []string

	uploadFilePath := request.URL.Query().Get("fileType")
	if uploadFilePath == "" {
		utils.JSONResponse(writer, "fileType is required", http.StatusBadRequest)
		return
	}

	err := request.ParseMultipartForm(maxUploadFile)
	if err != nil {
		utils.JSONResponse(writer, "could not parse files", http.StatusBadRequest)
		return
	}
	files := request.MultipartForm.File["files"]

	if len(files) == 0 {
		utils.JSONResponse(writer, "no files provided", http.StatusBadRequest)
		return
	}
	for _, fileHeader := range files {

		file, err := fileHeader.Open()

		if err != nil {
			utils.JSONResponse(writer, "could not open file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		cloudFlareURL, err := utils.UploadPictures(fileHeader.Filename, uploadFilePath, file)

		if err != nil {
			utils.JSONResponse(writer, "could not upload file", http.StatusBadRequest)
			return
		}

		fileURL = append(fileURL, cloudFlareURL)

	}
	detail, _ := json.Marshal(fileURL)
	utils.DSJsonResponse(writer, detail, http.StatusOK)
}

// Misc godoc
// @Tags Misc
// @Summary DeleteFile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body schemas.DeletePayload true "DeleteFile Payload"
// @Router /api/v1/misc/delete-file [post]
// @Success 200 {object} map[string]interface{}
// @Failure      400  {object} schemas.ErrorPayload
func DeleteFile(writer http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	var DeletePayload schemas.DeletePayload

	err := json.Unmarshal(body, &DeletePayload)

	if err != nil {

		utils.JSONResponse(writer, "delete file body not valid", http.StatusUnprocessableEntity)
		return
	}

	err = validate.Struct(DeletePayload)
	if err != nil {
		validationError := err.(validator.ValidationErrors)
		utils.JSONResponse(writer, validationError.Error(), http.StatusBadRequest)
		return
	}

	err = utils.DeletePicture(DeletePayload.File)

	if err != nil {
		utils.JSONResponse(writer, "delete file failed", http.StatusBadRequest)
		return
	}
	utils.DSJsonResponse(writer, []byte(`{}`), http.StatusOK)
	return

}
