// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package gcfaccess

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
		http.Error(w, "Error while parsing access grant", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	var edgeConfig = edge.Config{AuthServiceAddress: "auth.storjshare.io:7777"}
	response, err := edgeConfig.RegisterAccess(ctx, access, &edge.RegisterAccessOptions{Public: request.Public})
	if err != nil {
		http.Error(w, "Error while registering access grant with Auth Service", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
