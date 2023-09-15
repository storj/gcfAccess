package cloudsigning

import (
	"encoding/json"
	"net/http"

	"storj.io/uplink"
	"storj.io/uplink/edge"
)

func RegisterAccess(w http.ResponseWriter, r *http.Request) {
	if handleCORS(w, r) {
		return
	}
	var request struct {
		AccessGrant string `json:"access_grant"`
		Public      bool   `json:"public"`
	}
	if parseBodyJson(w, r, &request) {
		return
	}
	access, err := uplink.ParseAccess(request.AccessGrant)
	if err != nil {
		r.Response.StatusCode = http.StatusBadRequest
		r.Response.Status = "Error while parsing access grant"
		return
	}

	ctx := r.Context()
	var edgeConfig = edge.Config{AuthServiceAddress: "auth.storjshare.io:7777"}
	response, err := edgeConfig.RegisterAccess(ctx, access, &edge.RegisterAccessOptions{Public: request.Public})
	if err != nil {
		r.Response.StatusCode = http.StatusInternalServerError
		r.Response.Status = "Error while registering access grant with Auth Service"
	}

	// var response struct {
	// 	AccessKeyID string `json:"AccessKeyID"`
	// 	SecretKey   string `json:"SecretKey"`
	// 	Endpoint    string `json:"Endpoint"`
	// }
	// response.AccessKeyID = s3.AccessKeyID
	// response.SecretKey = s3.SecretKey
	// response.Endpoint = s3.Endpoint

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
