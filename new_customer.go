package cloudsigning

import (
	"crypto/rand"
	"encoding/json"
	"net/http"

	"storj.io/uplink"
	"storj.io/uplink/edge"
)

func NewS3Customer(w http.ResponseWriter, r *http.Request) {
	if handleCORS(w, r) {
		return
	}
	var request struct {
		AccessGrant string             `json:"access_grant"`
		Path        uplink.SharePrefix `json:"paths"`
		Permission  uplink.Permission  `json:"permission"`
		Passphrase  string             `json:"passphrase"`
		SaltBytes   []byte             `json:"salt"`
		Public      bool               `json:"public"`
	}
	// input parsing and validation
	if parseBodyJson(w, r, &request) {
		return
	}
	originalAccess, err := uplink.ParseAccess(request.AccessGrant)
	if err != nil {
		r.Response.StatusCode = http.StatusBadRequest
		r.Response.Status = "Error while parsing access grant"
		return
	}
	if request.Permission == (uplink.Permission{}) {
		r.Response.StatusCode = http.StatusUnprocessableEntity
		r.Response.Status = "Error: no permissions specified"
		return
	}
	if len(request.Path.Bucket) == 0 || len(request.Path.Prefix) == 0 {
		r.Response.StatusCode = http.StatusBadRequest
		r.Response.Status = "Error: no or invalid path specified"
		return
	}
	if len(request.Passphrase) == 0 {
		r.Response.StatusCode = http.StatusBadRequest
		r.Response.Status = "Error: no passphrase specified"
		return
	}
	if len(request.SaltBytes) == 0 {
		request.SaltBytes = make([]byte, 32)
		_, err = rand.Read(request.SaltBytes)
		if err != nil {
			r.Response.StatusCode = http.StatusInternalServerError
			r.Response.Status = "Error generating salt bytes"
			return
		}
	}

	// restrict
	derivedAccess, err := originalAccess.Share(request.Permission, request.Path)
	if err != nil {
		r.Response.StatusCode = http.StatusInternalServerError
		r.Response.Status = "Error while deriving restricted access grant"
	}

	//override encryption
	saltedUserKey, err := uplink.DeriveEncryptionKey(request.Passphrase, request.SaltBytes)
	if err != nil {
		r.Response.StatusCode = http.StatusInternalServerError
		r.Response.Status = "Error deriving salted encryption key"
		return
	}
	err = derivedAccess.OverrideEncryptionKey(request.Path.Bucket, request.Path.Prefix, saltedUserKey)
	if err != nil {
		r.Response.StatusCode = http.StatusInternalServerError
		r.Response.Status = "Error overriding encryption key"
		return
	}

	//register for S3
	ctx := r.Context()
	var edgeConfig = edge.Config{AuthServiceAddress: "auth.storjshare.io:7777"}
	response, err := edgeConfig.RegisterAccess(ctx, derivedAccess, &edge.RegisterAccessOptions{Public: request.Public})
	if err != nil {
		r.Response.StatusCode = http.StatusInternalServerError
		r.Response.Status = "Error while registering access grant with Auth Service"
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
