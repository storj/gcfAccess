// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package gcfaccess

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
		Path        uplink.SharePrefix `json:"path"`
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
		http.Error(w, "Error while parsing access grant", http.StatusBadRequest)
		return
	}
	if request.Permission == (uplink.Permission{}) {
		http.Error(w, "Error: no permissions specified", http.StatusBadRequest)
		return
	}
	if len(request.Path.Bucket) == 0 || len(request.Path.Prefix) == 0 {
		http.Error(w, "Error: no or invalid path specified", http.StatusBadRequest)
		return
	}
	if len(request.Passphrase) == 0 {
		http.Error(w, "Error: no passphrase specified", http.StatusBadRequest)
		return
	}
	if len(request.SaltBytes) == 0 {
		request.SaltBytes = make([]byte, 32)
		_, err = rand.Read(request.SaltBytes)
		if err != nil {
			http.Error(w, "Error generating salt bytes", http.StatusInternalServerError)
			return
		}
	}

	// restrict
	derivedAccess, err := originalAccess.Share(request.Permission, request.Path)
	if err != nil {
		http.Error(w, "Error while deriving restricted access grant", http.StatusInternalServerError)
		return
	}
	serializedAccess, err := derivedAccess.Serialize()
	if err != nil {
		http.Error(w, "Error while serializing new access grant", http.StatusInternalServerError)
		return
	}

	//override encryption
	saltedUserKey, err := uplink.DeriveEncryptionKey(request.Passphrase, request.SaltBytes)
	if err != nil {
		http.Error(w, "Error deriving salted encryption key", http.StatusInternalServerError)
		return
	}
	err = derivedAccess.OverrideEncryptionKey(request.Path.Bucket, request.Path.Prefix, saltedUserKey)
	if err != nil {
		http.Error(w, "Error overriding encryption key", http.StatusInternalServerError)
		return
	}

	//register for S3
	ctx := r.Context()
	var edgeConfig = edge.Config{AuthServiceAddress: "auth.storjshare.io:7777"}
	edgeCreds, err := edgeConfig.RegisterAccess(ctx, derivedAccess, &edge.RegisterAccessOptions{Public: request.Public})
	if err != nil {
		http.Error(w, "Error while registering access grant with Auth Service", http.StatusInternalServerError)
		return
	}

	//note that the access grant that is returned does not have an updated encryption passphrase
	//however, because it shares the same API key, revoking it will revoke the registered S3 grant
	var response struct {
		Edge            *edge.Credentials `json:"edge"`
		RevocableAccess string            `json:"revocable"`
	}
	response.Edge = edgeCreds
	response.RevocableAccess = serializedAccess

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
