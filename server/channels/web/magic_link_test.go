// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package web

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mattermost/mattermost/server/public/model"
)

func TestMagicLinkRedirectUsesConfiguredSiteURL(t *testing.T) {
	th := Setup(t)

	const siteURL = "http://mattermost.example"
	th.App.UpdateConfig(func(cfg *model.Config) {
		*cfg.ServiceSettings.SiteURL = siteURL
		*cfg.GuestAccountsSettings.EnableGuestMagicLink = false
	})

	handler := th.Web.APIHandler(loginWithMagicLinkToken)

	for _, testCase := range []struct {
		name string
		host string
	}{
		{name: "configured host control", host: "mattermost.example"},
		{name: "untrusted request host", host: "evil.example"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/login/one_time_link?t=invalid&redirect_to=%2Fchannels%2Fx", nil)
			req.Host = testCase.host
			res := httptest.NewRecorder()

			handler.ServeHTTP(res, req)

			require.Equal(t, http.StatusFound, res.Code)
			require.Equal(t, siteURL+"/login?extra=login_error&message=Magic link is disabled", res.Header().Get("Location"), "browser redirects must not trust the request Host")
		})
	}
}
