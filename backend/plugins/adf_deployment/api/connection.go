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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"

	"github.com/apache/incubator-devlake/core/plugin"
	adfDeploymentHelper "github.com/apache/incubator-devlake/plugins/adf_deployment/helper"
	"github.com/apache/incubator-devlake/plugins/adf_deployment/models"
	"github.com/apache/incubator-devlake/server/api/shared"
)

type AdfDeploymentTestConnResponse struct {
	shared.ApiBody
	Connection *models.AdfConn
}

type ReturnObject struct {
	models.AdfConnection `mapstructure:",squash"`
	Credentials           map[string]interface{} `mapstructure:"credentials" json:"credentials"`
}

// @Summary test adf_deployment connection
// @Description Test adf_deployment Connection. endpoint: "https://dev.adf_deployment.com/{organization}/
// @Tags plugins/adf_deployment
// @Param body body models.KubeConn true "json body"
// @Success 200  {object} KubeDeploymentTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/adf_deployment/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	
	body := AdfDeploymentTestConnResponse{}
	

	body.Success = true
	body.Message = "success"
	body.Connection = nil

	credentialsInput, ok := input.Body["credentials"].(map[string]interface{})

	if ok {
		_, err := adfDeploymentHelper.NewAdfApiClient(credentialsInput)
		if err != nil {
			body.Success = false
			body.Message = "Unable to establish connection to ADF"
			return &plugin.ApiResourceOutput{Body: body, Status: 400}, nil
		}
	}

	// output
	return &plugin.ApiResourceOutput{Body: body, Status: 200}, nil
}

// TODO Please modify the folowing code to adapt to your plugin
// @Summary create adf_deployment connection
// @Description Create adf_deployment connection
// @Tags plugins/adf_deployment
// @Param body body models.AdfConnection true "json body"
// @Success 200  {object} models.AdfConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/adf_deployment/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// update from request and save to database
	connection := &models.AdfConnection{}


	credentials := input.Body["credentials"].(map[string]interface{})
	token, _ := json.Marshal(credentials)

	if input.Body == nil {
		return nil, errors.BadInput.New("missing connectionId")
	}
	fmt.Println("token: -->", string(token))
	input.Body["credentials"] = string(token)
	fmt.Println("input.Body[credentials]: -->", input.Body["credentials"])
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}

	returnObject := ReturnObject{
		AdfConnection: *connection,
		Credentials:    credentials,
	}

	return &plugin.ApiResourceOutput{Body: returnObject, Status: http.StatusOK}, nil
}

// TODO Please modify the folowing code to adapt to your plugin
// @Summary patch adf_deployment connection
// @Description Patch adf_deployment connection
// @Tags plugins/adf_deployment
// @Param body body models.AdfConnection true "json body"
// @Success 200  {object} models.AdfConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/adf_deployment/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.AdfConnection{}

	credentialsInput := input.Body["credentials"].(map[string]interface{})
	token, _ := json.Marshal(credentialsInput)

	if input.Body == nil {
		return nil, errors.BadInput.New("missing connectionId")
	}

	input.Body["credentials"] = string(token)

	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}

	var credentials map[string]interface{}
	unmarshalErr := json.Unmarshal([]byte(connection.Credentials), &credentials)

	if unmarshalErr != nil {
		return nil, errors.BadInput.New("credentials is not a valid json")
	}

	returnObject := ReturnObject{
		AdfConnection: *connection,
		Credentials:    credentials,
	}

	return &plugin.ApiResourceOutput{Body: returnObject}, nil
}

// @Summary delete a adf_deployment connection
// @Description Delete a adf_deployment connection
// @Tags plugins/adf_deployment
// @Success 200  {object} models.AdfConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/adf_deployment/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.AdfConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &plugin.ApiResourceOutput{Body: connection}, err
}

// @Summary get all adf_deployment connections
// @Description Get all adf_deployment connections
// @Tags plugins/adf_deployment
// @Success 200  {object} []models.AdfConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/adf_deployment/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connections []models.AdfConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}

	returnObjects := make([]ReturnObject, len(connections))
	for i, connection := range connections {
		var credentials map[string]interface{}
		unmarshalErr := json.Unmarshal([]byte(connection.Credentials), &credentials)

		if unmarshalErr != nil {
			returnObjects[i].Credentials = map[string]interface{}{}
		} else {
			returnObjects[i].Credentials = credentials
		}
		returnObjects[i].AdfConnection = connection
	}

	return &plugin.ApiResourceOutput{Body: returnObjects, Status: http.StatusOK}, nil
}

// TODO Please modify the folowing code to adapt to your plugin
// @Summary get adf_deployment connection detail
// @Description Get adf_deployment connection detail
// @Tags plugins/adf_deployment
// @Success 200  {object} models.AdfConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/adf_deployment/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.AdfConnection{}
	err := connectionHelper.First(connection, input.Params)

	var credentials map[string]interface{}
	unmarshalErr := json.Unmarshal([]byte(connection.Credentials), &credentials)

	if unmarshalErr != nil {
		return nil, errors.BadInput.New("credentials is not a valid json")
	}

	returnObject := ReturnObject{
		AdfConnection: *connection,
		Credentials:    credentials,
	}
	return &plugin.ApiResourceOutput{Body: returnObject}, err
}

// @Summary Get kubernetes namespaces
// @Description Get kubernetes namespaces
// @Tags plugins/adf_deployment
// @Success 200  {object} models.AdfConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/adf_deployment/connections/{connectionId}/namespaces [GET]
func GetNameSpaces(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	fmt.Println("input.Params: -->", input.Params)
	
	body := AdfDeploymentTestConnResponse{}

	return &plugin.ApiResourceOutput{Body: body, Status: http.StatusOK}, nil
}

// @Summary Get kubernetes namespaces
// @Description Get kubernetes namespaces
// @Tags plugins/adf_deployment
// @Success 200  {object} models.AdfConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/adf_deployment/connections/{connectionId}/{namespace}/deployments [GET]
func GetDeployments(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	fmt.Println("input.Params: -->", input.Params)

	body := AdfDeploymentTestConnResponse{}

	return &plugin.ApiResourceOutput{Body: body, Status: http.StatusOK}, nil
}
