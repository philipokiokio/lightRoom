package utils

import (
	"encoding/json"
	"net/http"
)

func JSONResponse(writer http.ResponseWriter, detail string, statusCode int) {

	writer.Header().Set("Content-Type", "application/json")
	jsonResponse, _ := json.Marshal(map[string]string{"detail": detail})
	writer.WriteHeader(statusCode)
	writer.Write(jsonResponse)

}

func DSJsonResponse(writer http.ResponseWriter, detail []byte, statusCode int) {

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	writer.Write(detail)

}
