package cloudsigning

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

const (
	maxPostSize = 4 * 1024
)

func init() {
	functions.HTTP("RevokeAccess", RevokeAccess)
	functions.HTTP("RestrictAccess", RestrictAccess)
	functions.HTTP("OverrideEncryption", OverrideEncryption)
	functions.HTTP("RegisterAccess", RegisterAccess)
	functions.HTTP("NewS3Customer", NewS3Customer)
}

func handleCORS(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Accept-Language, Content-Language, Content-Length, Accept-Encoding")
		w.WriteHeader(http.StatusNoContent)
		return true
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	return false
}

func parseBodyJson(w http.ResponseWriter, r *http.Request, request any) bool {
	reader := http.MaxBytesReader(w, r.Body, maxPostSize)
	if err := json.NewDecoder(reader).Decode(request); err != nil {
		status := http.StatusUnprocessableEntity
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			status = http.StatusRequestEntityTooLarge
		}
		r.Response.StatusCode = status
		r.Response.Status = "Error while parsing request JSON"
		return true
	}
	return false
}
