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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/datafactory/mgmt/datafactory"
	"github.com/apache/incubator-devlake/core/errors"

	"github.com/Azure/go-autorest/autorest/date"
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

type PipelineRunReturnObject struct {
	ID string
	Status string
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
		adfApiClient, err := adfDeploymentHelper.NewAdfApiClient(credentialsInput)
		if err != nil {
			body.Success = false
			body.Message = "Unable to establish connection to ADF"
			return &plugin.ApiResourceOutput{Body: body, Status: 400}, nil
		}

		factoryName := credentialsInput["factoryName"].(string)
		resourceGroupName := credentialsInput["resourceGroupName"].(string)		

		df, dfError := adfApiClient.Client.Get(context.Background(), resourceGroupName, factoryName, "")
		if dfError != nil {
			body.Success = false
			body.Message = "Unable to establish connection to ADF"
			return &plugin.ApiResourceOutput{Body: body, Status: 400}, nil
		}
		fmt.Printf("Data Factory Name: %v\n", *df.Name)
		fmt.Printf("Data Factory ID: %v\n", *df.ID)

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

// @Summary Get adf pipelines
// @Description Get adf_deployment pipelines
// @Tags plugins/adf_deployment
// @Success 200  {object} models.AdfConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/adf_deployment/connections/{connectionId}/pipelines [GET]
func GetPipelines(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	fmt.Println("input.Params: -->", input.Params)

	connection := &models.AdfConnection{}
	err := connectionHelper.First(connection, input.Params)

	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get ADF API client instance")
	}

	var credentials map[string]interface{}
	unmarshalErr := json.Unmarshal([]byte(connection.Credentials), &credentials)
	if unmarshalErr != nil {
		return nil, errors.BadInput.New("credentials is not a valid json")
	}

	body := AdfDeploymentTestConnResponse{}
	body.Success = true
	body.Message = "success"
	body.Connection = nil
	
	adfPipelinesClient, err  := adfDeploymentHelper.NewAdfPipelinesClient(credentials)
	if err != nil {
		body.Success = false
		body.Message = "Unable to establish connection to ADF"
		return &plugin.ApiResourceOutput{Body: body, Status: 400}, nil
	}

	factoryName := credentials["factoryName"].(string)
	resourceGroupName := credentials["resourceGroupName"].(string)

	pipelines, pipelinesError := adfPipelinesClient.Client.ListByFactory(context.Background(), resourceGroupName, factoryName)
	if pipelinesError != nil {
		body.Success = false
		body.Message = "Unable to establish connection to ADF"
		return &plugin.ApiResourceOutput{Body: body, Status: 400}, nil
	}

	var pipelineList = make([]string, len(pipelines.Values()))

	for i, pipeline := range pipelines.Values() {
		pipelineList[i] = *pipeline.Name
	
	}

	return &plugin.ApiResourceOutput{Body: pipelineList, Status: http.StatusOK}, nil

}


// @Summary Get adf pipeline runs
// @Description Get adf_deployment pipeline runs
// @Tags plugins/adf_deployment
// @Success 200  {object} models.AdfConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/adf_deployment/connections/{connectionId}/pipelines/{pipeline}/runs [GET]
func GetPipelineRuns(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	fmt.Println("input.Params: -->", input.Params)

	
	daysAgo := input.Query.Get("dayAgo")

	if daysAgo == ""{
		daysAgo = "30"
	}

	fmt.Println("daysAgo: -->", daysAgo)



	connection := &models.AdfConnection{}
	err := connectionHelper.First(connection, input.Params)

	

	if err != nil {
		return nil, errors.Default.Wrap(err, "unable to get ADF API client instance")
	}

	

	var credentials map[string]interface{}
	unmarshalErr := json.Unmarshal([]byte(connection.Credentials), &credentials)

	if unmarshalErr != nil {
		return nil, errors.BadInput.New("credentials is not a valid json")
	}

	

	body := AdfDeploymentTestConnResponse{}
	body.Success = true
	body.Message = "success"
	body.Connection = nil


	adfPipelineRunsClient, err := adfDeploymentHelper.NewAdfPipelineRunsClient(credentials)
	if err != nil {
		fmt.Println("Error PPERROR:", err)
		body.Success = false
		body.Message = "Unable to establish connection to ADF"
		return &plugin.ApiResourceOutput{Body: body, Status: 400}, nil
	}

	factoryName := credentials["factoryName"].(string)
	resourceGroupName := credentials["resourceGroupName"].(string)
	pipelineName := input.Params["pipeline"]


	var filters []datafactory.RunQueryFilter


	filters = append(filters, datafactory.RunQueryFilter{
		Operand: "PipelineName",
		Operator: "Equals",
		Values: &[]string{pipelineName},
	})

	daysAgoInt, _ := strconv.Atoi(daysAgo)
	lastUpdatedAfter := date.Time{Time: time.Now().AddDate(0, 0, -daysAgoInt)} // 30 days ago

	filterParams := datafactory.RunFilterParameters{
		LastUpdatedAfter:  &lastUpdatedAfter,
		LastUpdatedBefore: &date.Time{Time: time.Now()},
		Filters:           &filters,
	}


	pipelineRuns, pipelineRunsError := adfPipelineRunsClient.Client.QueryByFactory(context.Background(), resourceGroupName, factoryName, filterParams)

	if pipelineRunsError != nil {
		fmt.Println("pipelineRunsError:", pipelineRunsError)
		body.Success = false
		body.Message = "Unable to retrieve pipeline runs from ADF"
		return &plugin.ApiResourceOutput{Body: body, Status: 400}, nil
	}

	

	var pipelineRunList = make([]PipelineRunReturnObject, len(*pipelineRuns.Value))

	for i, pipelineRun := range *pipelineRuns.Value {
		pipelineRunList[i] = PipelineRunReturnObject{
			ID: *pipelineRun.RunID,
			Status: *pipelineRun.Status,
		}
	}

	return &plugin.ApiResourceOutput{Body: pipelineRunList, Status: http.StatusOK}, nil

}
