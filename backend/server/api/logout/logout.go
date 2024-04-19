/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package logout

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/server/api/shared"
	"github.com/gin-gonic/gin"
)

type LogoutRequest struct {
}

type LogoutResponse struct {
}
type ChallengeParameters struct {
}

// @Summary post logout
// @Description post logout
// @Tags framework/logout
// @Accept application/json
// @Param blueprint body LogoutRequest true "json"
// @Success 200  {object} LogoutResponse
// @Router /logout [post]
func Logout(ctx *gin.Context) {
	v := config.GetConfig()
	// Remove the HTTP-only cookies
	res := LogoutResponse{}
	cookieNames := []string{"AWSALBAuthNonce", "AWSELBAuthSessionCookie-0", "AWSELBAuthSessionCookie-1"}

	for _, cookieName := range cookieNames {
		cookie := &http.Cookie{
			Name:     cookieName,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			MaxAge:   -1,
		}
		http.SetCookie(ctx.Writer, cookie)
	}
	encodedLoginRedirect := url.QueryEscape(v.GetString("AWS_COGNITO_LOGIN_REDIRECT"))
	location := fmt.Sprintf("%s/logout?response_type=code&client_id=%s&redirect_uri=%s&scope=openid", v.GetString("AWS_COGNITO_DOMAIN"), v.GetString("AWS_COGNITO_CLIENT_ID"), encodedLoginRedirect)

	ctx.Redirect(http.StatusFound, location)
	shared.ApiOutputSuccess(ctx, res, http.StatusOK)
}
