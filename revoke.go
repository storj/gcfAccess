// Copyright (C) 2023 Storj Labs, Inc.
// See LICENSE for copying information.

package gcfaccess

import (
	"net/http"

	"storj.io/uplink"
)

func RevokeAccess(w http.ResponseWriter, r *http.Request) {
	if handleCORS(w, r) {
		return
	}
	var request struct {
		AuthGrant   string `json:"authorizing_access_grant"`
		RevokeGrant string `json:"access_grant_to_revoke"`
	}

	// input parsing and validation
	if parseBodyJson(w, r, &request) {
		return
	}
	authAccess, err := uplink.ParseAccess(request.AuthGrant)
	if err != nil {
		r.Response.StatusCode = http.StatusBadRequest
		r.Response.Status = "Error while parsing authorizing access grant"
		return
	}
	accessToRevoke, err := uplink.ParseAccess(request.RevokeGrant)
	if err != nil {
		r.Response.StatusCode = http.StatusBadRequest
		r.Response.Status = "Error while parsing access grant to be revoked"
		return
	}

	ctx := r.Context()
	config := &uplink.Config{UserAgent: "authservice"}
	project, err := config.OpenProject(ctx, authAccess)
	if err != nil {
		r.Response.StatusCode = http.StatusInternalServerError
		r.Response.Status = "Error while opening project"
		return
	}
	defer func() { _ = project.Close() }()

	if err := project.RevokeAccess(ctx, accessToRevoke); err != nil {
		r.Response.StatusCode = http.StatusInternalServerError
		r.Response.Status = "Error while revoking access grant"
		return
	}
	w.WriteHeader(http.StatusOK)
}
