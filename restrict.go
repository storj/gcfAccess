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
		r.Response.StatusCode = http.StatusBadRequest
		r.Response.Status = "Error while parsing access grant"
		return
	}
	if request.Permission == (uplink.Permission{}) {
		r.Response.StatusCode = http.StatusUnprocessableEntity
		r.Response.Status = "Error: no permissions specified"
		return
	}
	// if we aren't actually restricting anything, then we don't need to Share.
	if request.Permission == (uplink.Permission{
		AllowDelete:   true,
		AllowList:     true,
		AllowDownload: true,
		AllowUpload:   true,
	}) && len(request.Paths) == 0 {
		r.Response.StatusCode = http.StatusUnprocessableEntity
		r.Response.Status = "Error: no restrictions specified"
		return
	}
	// note that Share() fail if permission is empty, but not paths
	derivedAccess, err := originalAccess.Share(request.Permission, request.Paths...)
	if err != nil {
		r.Response.StatusCode = http.StatusInternalServerError
		r.Response.Status = "Error while deriving restricted access grant"
	}
	var response struct {
		AccessGrant string `json:"access_grant"`
	}
	serializedAccess, err := derivedAccess.Serialize()
	if err != nil {
		r.Response.StatusCode = http.StatusInternalServerError
		r.Response.Status = "Error while serializing restricted access grant"
	}
	response.AccessGrant = serializedAccess
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
