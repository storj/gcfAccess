// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package gcfaccess

import (
	"encoding/json"
	"net/http"

	"storj.io/uplink"
)

func RestrictAccess(w http.ResponseWriter, r *http.Request) {
	if handleCORS(w, r) {
		return
	}
	var request struct {
		AccessGrant string               `json:"access_grant"`
		Paths       []uplink.SharePrefix `json:"paths"`
		Permission  uplink.Permission    `json:"permission"`
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
		http.Error(w, "Error: no permissions specified", http.StatusUnprocessableEntity)
		return
	}
	// if we aren't actually restricting anything, then we don't need to Share.
	if request.Permission == (uplink.Permission{
		AllowDelete:   true,
		AllowList:     true,
		AllowDownload: true,
		AllowUpload:   true,
	}) && len(request.Paths) == 0 {
		http.Error(w, "Error: no restrictions specified", http.StatusUnprocessableEntity)
		return
	}
	// note that Share() fail if permission is empty, but not paths
	derivedAccess, err := originalAccess.Share(request.Permission, request.Paths...)
	if err != nil {
		http.Error(w, "Error while deriving restricted access grant", http.StatusInternalServerError)
		return
	}
	var response struct {
		AccessGrant string `json:"access_grant"`
	}
	serializedAccess, err := derivedAccess.Serialize()
	if err != nil {
		http.Error(w, "Error while serializing restricted access grant", http.StatusInternalServerError)
		return
	}
	response.AccessGrant = serializedAccess
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
