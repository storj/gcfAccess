// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package gcfaccess

import (
	"crypto/rand"
	"encoding/json"
	"net/http"

	"storj.io/uplink"
)

func OverrideEncryption(w http.ResponseWriter, r *http.Request) {
	if handleCORS(w, r) {
		return
	}
	var request struct {
		AccessGrant string             `json:"access_grant"`
		Passphrase  string             `json:"passphrase"`
		SaltBytes   []byte             `json:"salt"`
		Path        uplink.SharePrefix `json:"path"`
	}
	// input parsing and validation
	if parseBodyJson(w, r, &request) {
		return
	}
	access, err := uplink.ParseAccess(request.AccessGrant)
	if err != nil {
		http.Error(w, "Error while parsing access grant", http.StatusBadRequest)
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
	saltedUserKey, err := uplink.DeriveEncryptionKey(request.Passphrase, request.SaltBytes)
	if err != nil {
		http.Error(w, "Error deriving salted encryption key", http.StatusInternalServerError)
		return
	}
	err = access.OverrideEncryptionKey(request.Path.Bucket, request.Path.Prefix, saltedUserKey)
	if err != nil {
		http.Error(w, "Error overriding encryption key", http.StatusInternalServerError)
		return
	}
	var response struct {
		AccessGrant string `json:"access_grant"`
	}
	serializedAccess, err := access.Serialize()
	if err != nil {
		http.Error(w, "Error while serializing new access grant", http.StatusInternalServerError)
		return
	}
	response.AccessGrant = serializedAccess
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
