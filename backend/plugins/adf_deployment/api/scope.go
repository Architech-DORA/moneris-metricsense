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

package api

import (
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/adf_deployment/models"
)

type ScopeRes struct {
	models.AdfDeployment
	TransformationRuleName string `json:"transformationRuleName,omitempty"`
}

type ScopeReq api.ScopeReq[models.AdfDeployment]

// PutScope create or update adf_deployment board
// @Summary create or update adf_deployment board
// @Description Create or update Jira board
// @Tags plugins/adf_deployment
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param scope body ScopeReq true "json"
// @Success 200  {object} []models.JiraBoard
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/adf_deployment/connections/{connectionId}/scopes [PUT]
func PutScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	print(input.Body, "put_input.Body")
	return scopeHelper.Put(input)
}

// UpdateScope patch to adf_deployment board
// @Summary patch to adf_deployment board
// @Description patch to adf_deployment board
// @Tags plugins/adf_deployment
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param scopeId path int false "board ID"
// @Param scope body models.JiraBoard true "json"
// @Success 200  {object} models.JiraBoard
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/adf_deployment/connections/{connectionId}/scopes/{scopeId} [PATCH]
func UpdateScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeHelper.Update(input, "id")
}

// GetScopeList get Jira boards
// @Summary get Jira boards
// @Description get Jira boards
// @Tags plugins/adf_deployment
// @Param connectionId path int false "connection ID"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Success 200  {object} []ScopeRes
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/adf_deployment/connections/{connectionId}/scopes/ [GET]
func GetScopeList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeHelper.GetScopeList(input)
}

// GetScope get one Jira board
// @Summary get one Jira board
// @Description get one Jira board
// @Tags plugins/adf_deployment
// @Param connectionId path int false "connection ID"
// @Param scopeId path int false "board ID"
// @Success 200  {object} ScopeRes
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/adf_deployment/connections/{connectionId}/scopes/{scopeId} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	print(input, "GetScope_input")
	return scopeHelper.GetScope(input, "id")
}

// func GetDeployment(op *tasks.AdfDeploymentOptions, apiClient aha.ApiClientAbstract) (*models.AdfDeploymentRevision, errors.Error) {
// 	boardRes := &models.AdfDeploymentRevision{}
// 	// TODO: get board id from op
// 	res, err := apiClient.Get(fmt.Sprintf("data/xx/%d", 1), nil, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()
// 	if res.StatusCode != http.StatusOK {
// 		return nil, errors.HttpStatus(res.StatusCode).New(fmt.Sprintf("unexpected status code when requesting repo detail from %s", res.Request.URL.String()))
// 	}
// 	body, err := errors.Convert01(io.ReadAll(res.Body))
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = errors.Convert(json.Unmarshal(body, boardRes))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return boardRes, nil
// }
